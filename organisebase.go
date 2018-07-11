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
	pbgh "github.com/brotherlogic/githubcard/proto"
	pbd "github.com/brotherlogic/godiscogs"
	pbgs "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

type prodGh struct{}

func (gh *prodGh) alert(r *pb.Location) error {
	host, port, err := utils.Resolve("githubcard")

	if err != nil {
		return err
	}

	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := pbgh.NewGithubClient(conn)
	_, err = client.AddIssue(ctx, &pbgh.Issue{Title: "Quota Issue", Body: fmt.Sprintf("%v is out of quota", r.GetName()), Service: "recordsorganiser"})
	return err
}

// Bridge that accesses discogs syncer server
type prodBridge struct {
	Resolver func(string) (string, int)
}

const (
	//KEY is where we store the org
	KEY = "github.com/brotherlogic/recordsorganiser/org"
)

func (s *Server) readOrg() error {
	org := &pb.Organisation{}
	data, _, err := s.KSclient.Read(KEY, org)

	if err != nil {
		return err
	}
	s.org = data.(*pb.Organisation)
	return nil
}

func (s *Server) saveOrg() {
	s.KSclient.Save(KEY, s.org)
}

func (discogsBridge prodBridge) GetIP(name string) (string, int) {
	return discogsBridge.Resolver(name)
}

func (discogsBridge prodBridge) getMetadata(rel *pbd.Release) (*pbrc.ReleaseMetadata, error) {
	ip, port := discogsBridge.GetIP("recordcollection")
	conn, err2 := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	if err2 == nil {
		defer conn.Close()
		client := pbrc.NewRecordCollectionServiceClient(conn)
		meta, err3 := client.GetRecords(context.Background(), &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: rel}})
		if err3 == nil && len(meta.Records) == 0 {
			return meta.Records[0].Metadata, nil
		} else {
			return nil, fmt.Errorf("Problem getting meta %v and %v", err3, len(meta.Records))
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

			rel, err3 := client.GetRecords(context.Background(), &pbrc.GetRecordsRequest{Force: true, Strip: true, Filter: &pbrc.Record{Release: &pbd.Release{FolderId: id}}})
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
	if master {
		err := s.readOrg()
		return err
	}

	return nil
}

// GetState gets the state of the server
func (s Server) GetState() []*pbgs.State {
	return []*pbgs.State{}
}

// InitServer builds an initial server
func InitServer() *Server {
	server := &Server{&goserver.GoServer{}, prodBridge{}, &pb.Organisation{}, &prodGh{}}
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
