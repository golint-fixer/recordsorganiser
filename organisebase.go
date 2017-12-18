package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbs "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	pbgs "github.com/brotherlogic/goserver/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

const (
	//Key the current state of the collection
	Key = "/github.com/brotherlogic/recordsorganiser/organisation"
)

// Bridge that accesses discogs syncer server
type prodBridge struct {
	Resolver func(string) (string, int)
}

func (discogsBridge prodBridge) GetIP(name string) (string, int) {
	return discogsBridge.Resolver(name)
}

func (discogsBridge prodBridge) getMetadata(rel *pbd.Release) (*pbs.ReleaseMetadata, error) {
	ip, port := discogsBridge.GetIP("discogssyncer")
	conn, err2 := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err2 == nil {
		defer conn.Close()
		client := pbs.NewDiscogsServiceClient(conn)
		meta, err2 := client.GetMetadata(context.Background(), rel)
		if err2 == nil {
			return meta, nil
		}
	}

	return nil, fmt.Errorf("Unable to get release metadata: %v", err2)
}

func (discogsBridge prodBridge) moveToFolder(move *pbs.ReleaseMove) {
	ip, port := discogsBridge.GetIP("discogssyncer")
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pbs.NewDiscogsServiceClient(conn)
	client.MoveToFolder(context.Background(), move)
}

func (discogsBridge prodBridge) getReleases(folders []int32) ([]*pbrc.Record, error) {
	var result []*pbrc.Record

	for _, id := range folders {
		ip, port := discogsBridge.GetIP("recordcollection")
		conn, err2 := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
		if err2 != nil {
			return result, err2
		}

		if err2 == nil {
			defer conn.Close()
			client := pbrc.NewRecordCollectionServiceClient(conn)

			rel, err3 := client.GetRecords(context.Background(), &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbd.Release{FolderId: id}}})
			if err3 != nil {
				return result, err3
			}
			result = append(result, rel.GetRecords()...)
		}
	}

	return result, nil
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterOrganiserServiceServer(server, s)
}

// Mote promotes/demotes this server
func (s *Server) Mote(master bool) error {
	return s.loadLatest()
}

// GetState gets the state of the server
func (s Server) GetState() []*pbgs.State {
	return []*pbgs.State{}
}

func (s Server) save() {

	// Always update the timestamp on a save
	s.org.Timestamp = time.Now().Unix()
	s.KSclient.Save(Key, s.org)
}

func (s Server) load(key string) (*pb.Organisation, error) {
	collection := &pb.Organisation{}
	data, _, err := s.KSclient.Read(key, collection)
	if err != nil {
		return nil, err
	}

	return data.(*pb.Organisation), nil
}

func (s *Server) loadLatest() error {
	curr, err := s.load(Key)
	if err != nil {
		return err
	}

	s.org = curr

	log.Printf("LOADED %v", curr)

	return nil
}

// InitServer builds an initial server
func InitServer() *Server {
	server := &Server{&goserver.GoServer{}, prodBridge{}, &pb.Organisation{}}
	server.PrepServer()
	server.bridge = &prodBridge{Resolver: server.GetIP}
	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)

	server.Register = server

	return server
}

// ReportHealth alerts if we're not healthy
func (s Server) ReportHealth() bool {
	return true
}

func main() {
	var quiet = flag.Bool("quiet", true, "Show log output")
	flag.Parse()

	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("Logging is on!")

	server := InitServer()

	server.GoServer.Killme = true
	server.RegisterServer("recordsorganiser", false)
	server.Serve()
}
