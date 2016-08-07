package main

import (
	"flag"
	"log"
	"sort"
	"strconv"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/brotherlogic/goserver"

	pbs "github.com/brotherlogic/discogssyncer/server"
	pbdi "github.com/brotherlogic/discovery/proto"
	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

// Bridge that accesses discogs syncer server
type prodBridge struct{}

func (discogsBridge prodBridge) getReleases(folders []int32) []*pbd.Release {
	var result []*pbd.Release

	list := &pbs.FolderList{}
	for _, id := range folders {
		list.Folders = append(list.Folders, &pbd.Folder{Id: id})
	}

	ip, port := getIP("discogsyncer", "10.0.1.17", 50055)
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()
	client := pbs.NewDiscogsServiceClient(conn)

	rel, _ := client.GetReleasesInFolder(context.Background(), list)
	result = rel.GetReleases()

	return result
}

// DoRegister does RPC registration
func (s Server) DoRegister(server *grpc.Server) {
	pb.RegisterOrganiserServiceServer(server, &s)
}

// InitServer builds an initial server
func InitServer(folder *string) Server {
	server := Server{&goserver.GoServer{}, *folder, prodBridge{}}
	server.Register = server
	return server
}

func getIP(servername string, ip string, port int) (string, int) {
	conn, _ := grpc.Dial(ip+":"+strconv.Itoa(port), grpc.WithInsecure())
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	entry := pbdi.RegistryEntry{Name: servername}
	r, _ := registry.Discover(context.Background(), &entry)
	return r.Ip, int(r.Port)
}

func test() {
	var host = flag.String("host", "10.0.1.17", "Hostname of server.")
	var port = flag.String("port", "50055", "Port number of server")
	flag.Parse()
	portVal, _ := strconv.Atoi(*port)
	dServer, dPort := getIP("discogssyncer", *host, portVal)

	dConn, _ := grpc.Dial(dServer+":"+strconv.Itoa(dPort), grpc.WithInsecure())
	defer dConn.Close()
	dClient := pbs.NewDiscogsServiceClient(dConn)

	list := &pbs.FolderList{}
	list.Folders = append(list.Folders, &pbd.Folder{Name: "12s"})

	releases, err := dClient.GetReleasesInFolder(context.Background(), list)

	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(pbd.ByLabelCat(releases.Releases))
	splits := pbd.Split(releases.Releases, 8)
	count := 1
	for _, split := range splits[0] {
		log.Printf("%v - %v", count, split.Title)
		count++
	}
}

func main() {
	var folder = flag.String("folder", "/home/simon/.discogs/", "Location to store the records")
	flag.Parse()
	server := InitServer(folder)

	server.PrepServer()
	server.RegisterServer("recordsorganiser", false)
	server.Serve()
}
