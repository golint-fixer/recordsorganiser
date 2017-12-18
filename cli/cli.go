package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"github.com/brotherlogic/keystore/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbdi "github.com/brotherlogic/discovery/proto"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

func findServer(name string) (string, int) {
	conn, err := grpc.Dial(utils.Discover, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot reach discover server: %v (trying to discover %v)", err, name)
	}
	defer conn.Close()

	registry := pbdi.NewDiscoveryServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	re := &pbdi.RegistryEntry{Name: name}
	r, err := registry.Discover(ctx, re)

	e, ok := status.FromError(err)
	if ok && e.Code() == codes.Unavailable {
		log.Printf("RETRY")
		r, err = registry.Discover(ctx, re)
	}

	if err != nil {
		return "", -1
	}
	return r.Ip, int(r.Port)
}

func read() *pb.Organisation {
	client := *keystoreclient.GetClient(findServer)
	thing, _, _ := client.Read("/github.com/brotherlogic/recordsorganiser/curr", &pb.Organisation{})
	return thing.(*pb.Organisation)
}

func main() {
	org := read()

	for i, o := range org.GetLocations() {
		fmt.Printf("%v. - %v\n", i, o.GetName())
	}
}
