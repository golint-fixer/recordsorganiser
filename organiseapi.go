package main

import (
	"errors"
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
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

const (
	//CurrKey the current state of the collection
	CurrKey = "/github.com/brotherlogic/recordsorganiser/curr"

	//PrevKey the old state of the collection
	PrevKey = "/github.com/brotherlogic/recordsorganiser/prev"

	backoffTime = time.Second * 5
	retries     = 5
)

// Bridge that accesses discogs syncer server
type prodBridge struct {
	Resolver func(string) (string, int)
}

func (discogsBridge prodBridge) GetIP(name string) (string, int) {
	return discogsBridge.Resolver(name)
}

func (discogsBridge prodBridge) getMetadata(rel *pbd.Release) (*pbs.ReleaseMetadata, error) {
	for i := 0; i < retries; i++ {
		ip, port := discogsBridge.GetIP("discogssyncer")
		conn, err := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
		if err == nil {
			defer conn.Close()
			client := pbs.NewDiscogsServiceClient(conn)
			meta, err := client.GetMetadata(context.Background(), rel)
			if err == nil {
				return meta, nil
			}
		}
		time.Sleep(backoffTime)
	}

	return nil, errors.New("Unable to get release metadata")
}

func (discogsBridge prodBridge) moveToFolder(move *pbs.ReleaseMove) {
	ip, port := discogsBridge.GetIP("discogssyncer")
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pbs.NewDiscogsServiceClient(conn)
	client.MoveToFolder(context.Background(), move)
}

func (discogsBridge prodBridge) getReleases(folders []int32) ([]*pbd.Release, error) {
	err := errors.New("First Pass Fail")
	for i := 0; i < retries; i++ {
		var result []*pbd.Release

		list := &pbs.FolderList{}
		for _, id := range folders {
			list.Folders = append(list.Folders, &pbd.Folder{Id: id})
		}

		ip, port := discogsBridge.GetIP("discogssyncer")
		conn, err2 := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
		err = err2

		if err == nil {
			defer conn.Close()
			client := pbs.NewDiscogsServiceClient(conn)

			rel, err3 := client.GetReleasesInFolder(context.Background(), list)
			err = err3
			if err == nil {
				result = rel.GetReleases()
				return result, nil
			}
		}
		time.Sleep(backoffTime)
	}

	return nil, fmt.Errorf("Unable to read releases: %v", err)
}

func (discogsBridge prodBridge) getRelease(ID int32) (*pbd.Release, error) {
	for i := 0; i < retries; i++ {
		ip, port := discogsBridge.GetIP("discogssyncer")
		conn, err := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
		if err == nil {

			defer conn.Close()
			client := pbs.NewDiscogsServiceClient(conn)

			rel, err := client.GetSingleRelease(context.Background(), &pbd.Release{Id: ID})
			if err == nil {
				return rel, err
			}
		}
		time.Sleep(backoffTime / time.Duration(retries))
	}
	return nil, errors.New("Unable to reach discogs")
}

func compare(collectionStart *pb.Organisation, collectionEnd *pb.Organisation) []*pb.LocationMove {
	var moves []*pb.LocationMove

	for _, folder := range collectionEnd.Locations {
		var matcher = &pb.Location{}
		for _, otherFolder := range collectionStart.Locations {
			if otherFolder.Name == folder.Name {
				matcher = otherFolder
			}
		}

		for i := 1; i <= int(folder.Units); i++ {
			diff := getMoves(matcher.ReleasesLocation, folder.ReleasesLocation, i, folder.Name)
			moves = append(moves, diff...)
		}
	}

	// Correct the slot moves
	for _, move := range moves {
		if !move.SlotMove {
			if move.New != nil {
				for _, location := range collectionStart.Locations {
					for _, placement := range location.ReleasesLocation {
						if placement.ReleaseId == move.New.ReleaseId {
							if location.Name == move.New.Folder {
								move.SlotMove = true
							}
						}
					}
				}
			} else {
				for _, location := range collectionStart.Locations {
					for _, placement := range location.ReleasesLocation {
						if placement.ReleaseId == move.Old.ReleaseId {
							if location.Name == move.Old.Folder {
								move.SlotMove = true
							}
						}
					}
				}
			}
		}
	}

	return moves
}

// DoRegister does RPC registration
func (s Server) DoRegister(server *grpc.Server) {
	pb.RegisterOrganiserServiceServer(server, &s)
}

// Mote promotes/demotes this server
func (s Server) Mote(master bool) error {
	return nil
}

func (s Server) save() {

	// Always update the timestamp on a save
	s.currOrg.Timestamp = time.Now().Unix()
	s.KSclient.Save(CurrKey, s.currOrg)
	s.KSclient.Save(PrevKey, s.pastOrg)
}

func (s Server) load(key string) (*pb.Organisation, error) {
	collection := &pb.Organisation{}
	data, err := s.KSclient.Read(key, collection)
	if err != nil {
		return nil, err
	}

	return data.(*pb.Organisation), nil
}

func (s *Server) loadLatest() error {
	curr, err := s.load(CurrKey)
	if err != nil {
		return err
	}
	past, err := s.load(PrevKey)
	if err != nil {
		return err
	}

	s.currOrg = curr
	s.pastOrg = past

	log.Printf("LOADED %v", curr)

	return nil
}

// InitServer builds an initial server
func InitServer() Server {
	server := Server{&goserver.GoServer{}, prodBridge{}, &pb.Organisation{}, &pb.Organisation{}}
	server.PrepServer()
	server.bridge = &prodBridge{Resolver: server.GetIP}
	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)
	err := server.loadLatest()

	if err != nil {
		panic(err)
	}

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
	log.Printf("Collection size: %v -> %v", len(server.currOrg.Locations), server.currOrg)
	server.Serve()
}
