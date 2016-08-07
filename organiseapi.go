package main

import (
	"flag"

	"google.golang.org/grpc"

	"github.com/brotherlogic/goserver"

	pb "github.com/brotherlogic/recordsorganiser/proto"
)

// Server the configuration for the syncer
type Server struct {
	*goserver.GoServer
	saveLocation string
}

// DoRegister does RPC registration
func (s Server) DoRegister(server *grpc.Server) {
	pb.RegisterOrganiserServiceServer(server, &s)
}

// InitServer builds an initial server
func InitServer(folder *string) Server {
	server := Server{&goserver.GoServer{}, *folder}
	server.Register = server
	return server
}

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs/", "Location to store the records")
	flag.Parse()
	server := InitServer(folder)

	server.PrepServer()
	server.RegisterServer("recordorganizer", false)
	server.Serve()
}
