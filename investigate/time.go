package main

import (
	"context"
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
	r, err := registry.Discover(ctx, &pbdi.DiscoverRequest{Request: re})

	e, ok := status.FromError(err)
	if ok && e.Code() == codes.Unavailable {
		log.Printf("RETRY")
		r, err = registry.Discover(ctx, &pbdi.DiscoverRequest{Request: re})
	}

	if err != nil {
		return "", -1
	}
	return r.GetService().Ip, int(r.GetService().Port)
}

func run() (time.Duration, int64, error) {
	s := time.Now()
	client := *keystoreclient.GetClient(findServer)
	_, details, err := client.Read("/github.com/brotherlogic/recordsorganiser/curr", &pb.Organisation{})
	return time.Now().Sub(s), details.GetReadTime(), err
}

func main() {
	for i := 0; i < 100; i++ {
		t, ot, err := run()
		if err != nil {
			log.Fatalf("Fatal error on run: %v", err)
		}
		if t.Seconds() > 0.5 {
			log.Printf("Excessive length of run: %v but %v", t, ot)
		}
	}
}
