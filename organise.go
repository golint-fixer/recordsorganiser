package main

import (
	"flag"
	"log"
	"sort"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/discogssyncer/server"
	pbdi "github.com/brotherlogic/discovery/proto"
	pbd "github.com/brotherlogic/godiscogs"
)

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
	dClient := pb.NewDiscogsServiceClient(dConn)

	list := &pb.FolderList{}
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
