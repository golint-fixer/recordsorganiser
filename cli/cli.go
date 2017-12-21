package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pbgd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordsorganiser/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func getReleaseString(instanceID int32) string {
	host, port, err := utils.Resolve("recordcollection")
	if err != nil {
		log.Fatalf("Unable to reach collection: %v", err)
	}
	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}

	client := pbrc.NewRecordCollectionServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	rel, err := client.GetRecords(ctx, &pbrc.GetRecordsRequest{Force: true, Filter: &pbrc.Record{Release: &pbgd.Release{InstanceId: instanceID}}})
	if err != nil {
		log.Fatalf("unable to get record: %v", err)
	}
	return rel.GetRecords()[0].GetRelease().Title
}

func get(ctx context.Context, client pb.OrganiserServiceClient, name string, force bool, slot int32) {
	locs, err := client.GetOrganisation(ctx, &pb.GetOrganisationRequest{ForceReorg: force, Locations: []*pb.Location{&pb.Location{Name: name}}})
	if err != nil {
		log.Fatalf("Error reading locations: %v", err)
	}

	for _, loc := range locs.GetLocations() {
		fmt.Printf("%v (%v)\n", loc.GetName(), len(loc.GetReleasesLocation()))
		for j, rloc := range loc.GetReleasesLocation() {
			if rloc.GetSlot() == slot {
				fmt.Printf("%v. %v\n", j, getReleaseString(rloc.GetInstanceId()))
			}
		}
	}

	fmt.Printf("Summary: %v\n", locs.GetNumberProcessed())

	if len(locs.GetLocations()) == 0 {
		fmt.Printf("No Locations Found!\n")
	}
}

func list(ctx context.Context, client pb.OrganiserServiceClient) {
	locs, err := client.GetOrganisation(ctx, &pb.GetOrganisationRequest{Locations: []*pb.Location{&pb.Location{}}})
	if err != nil {
		log.Fatalf("Error reading locations: %v", err)
	}

	for i, loc := range locs.GetLocations() {
		fmt.Printf("%v. %v\n", i, loc.GetName())
	}

	if len(locs.GetLocations()) == 0 {
		fmt.Printf("No Locations Found!\n")
	}
}

func add(ctx context.Context, client pb.OrganiserServiceClient, name string, folders []int32, slots int32) {
	loc, err := client.AddLocation(ctx, &pb.AddLocationRequest{Add: &pb.Location{Name: name, FolderIds: folders, Slots: slots}})

	if err != nil {
		log.Fatalf("Error adding location: %v", err)
	}

	fmt.Printf("Added location: %v\n", len(loc.GetNow().GetLocations()))
}

func main() {
	host, port, err := utils.Resolve("recordsorganiser")
	if err != nil {
		log.Fatalf("Unable to reach organiser: %v", err)
	}
	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}

	client := pb.NewOrganiserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	switch os.Args[1] {
	case "list":
		list(ctx, client)
	case "get":
		getLocationFlags := flag.NewFlagSet("GetLocation", flag.ExitOnError)
		var name = getLocationFlags.String("name", "", "The name of the location")
		var force = getLocationFlags.Bool("force", false, "Force a reorg")
		var slot = getLocationFlags.Int("slot", 1, "Slot to view")

		if err := getLocationFlags.Parse(os.Args[2:]); err == nil {
			get(ctx, client, *name, *force, int32(*slot))
		}
	case "add":
		addLocationFlags := flag.NewFlagSet("AddLocation", flag.ExitOnError)
		var name = addLocationFlags.String("name", "", "The name of the new location")
		var slots = addLocationFlags.Int("slots", 0, "The number of slots in the location")
		var folderIds = addLocationFlags.String("folders", "", "The list of folder IDs")

		if err := addLocationFlags.Parse(os.Args[2:]); err == nil {
			nums := make([]int32, 0)
			for _, folderID := range strings.Split(*folderIds, ",") {
				v, err := strconv.Atoi(folderID)
				if err != nil {
					log.Fatalf("Cannot parse folderid: %v", err)
				}
				nums = append(nums, int32(v))
			}
			add(ctx, client, *name, nums, int32(*slots))
		}
	}
}
