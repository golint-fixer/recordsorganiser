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

func locateRelease(ctx context.Context, c pb.OrganiserServiceClient, id int32) {
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
	recs, err := client.GetRecords(ctx, &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbgd.Release{Id: id}}})
	if err != nil {
		log.Fatalf("Unable to get record %v -> %v", id, err)
	}
	for _, rec := range recs.GetRecords() {
		location, err := c.Locate(ctx, &pb.LocateRequest{InstanceId: rec.GetRelease().InstanceId})
		if err != nil {
			fmt.Printf("Unable to locate instance (%v) of %v because %v\n", rec.GetRelease().InstanceId, rec.GetRelease().Title, err)
		} else {
			fmt.Printf("%v (%v) is in %v\n", rec.GetRelease().Title, rec.GetRelease().InstanceId, location.GetFoundLocation().GetName())

			for i, r := range location.GetFoundLocation().GetReleasesLocation() {
				if r.GetInstanceId() == rec.GetRelease().InstanceId {
					fmt.Printf("Slot %v\n", r.GetSlot())
					fmt.Printf("%v. %v\n", i-1, getReleaseString(location.GetFoundLocation().GetReleasesLocation()[i-1].InstanceId))
					fmt.Printf("%v. %v\n", i, getReleaseString(location.GetFoundLocation().GetReleasesLocation()[i].InstanceId))
					fmt.Printf("%v. %v\n", i+1, getReleaseString(location.GetFoundLocation().GetReleasesLocation()[i+1].InstanceId))
				}
			}
		}
	}
}

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
		fmt.Printf("%v (%v) -> %v [%v] with %v\n", loc.GetName(), len(loc.GetReleasesLocation()), loc.GetFolderIds(), loc.GetQuota(), loc.Sort.String())
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
	case "locate":
		locateFlags := flag.NewFlagSet("Locate", flag.ExitOnError)
		var id = locateFlags.Int("id", -1, "The id of the release")
		if err := locateFlags.Parse(os.Args[2:]); err == nil {
			locateRelease(ctx, client, int32(*id))
		}
	case "sell":
		sellFlags := flag.NewFlagSet("sell", flag.ExitOnError)
		var name = sellFlags.String("name", "", "The name of the location to get")

		if err := sellFlags.Parse(os.Args[2:]); err == nil {
			loc, err := client.GetOrganisation(ctx, &pb.GetOrganisationRequest{Locations: []*pb.Location{&pb.Location{Name: *name}}})
			if err != nil {
				log.Fatalf("ERRR: %v", err)
			}

			records := make([]*pbrc.Record, 0)
			minScore := int32(6)
			for _, l := range loc.GetLocations() {
				for _, id := range l.GetFolderIds() {
					host, port, err := utils.Resolve("recordcollection")
					if err != nil {
						log.Fatalf("Unable to reach collection: %v", err)
					}
					conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
					defer conn.Close()

					if err != nil {
						log.Fatalf("Unable to dial: %v", err)
					}

					rclient := pbrc.NewRecordCollectionServiceClient(conn)
					recs, err := rclient.GetRecords(context.Background(), &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbgd.Release{}, Metadata: &pbrc.ReleaseMetadata{GoalFolder: id}}})
					if err != nil {
						log.Fatalf("Error : %v", err)
					}

					for _, r := range recs.GetRecords() {
						if r.GetRelease().Rating > 0 {
							records = append(records, r)
							if r.GetRelease().Rating < minScore {
								minScore = r.GetRelease().Rating
							}
						}
					}
				}
			}

			log.Printf("FOUND %v [%v]", len(records), minScore)
			for _, r := range records {
				if r.GetRelease().Rating == minScore {
					fmt.Printf("SELL: %v\n", r.GetRelease().Title)
				}
			}
		}
	case "update":
		updateLocationFlags := flag.NewFlagSet("UpdateLocation", flag.ExitOnError)
		var name = updateLocationFlags.String("name", "", "The name of the new location")
		var folder = updateLocationFlags.Int("folder", 0, "The folder to add to the location")
		var quota = updateLocationFlags.Int("quota", 0, "The new quota to add to the location")
		var sort = updateLocationFlags.String("sort", "", "The new sorting mechanism")
		var alert = updateLocationFlags.Bool("alert", true, "Whether we should alert on this location")
		if err := updateLocationFlags.Parse(os.Args[2:]); err == nil {
			if *folder > 0 {
				client.UpdateLocation(ctx, &pb.UpdateLocationRequest{Location: *name, Update: &pb.Location{FolderIds: []int32{int32(*folder)}}})
			}
			if *quota > 0 {
				client.UpdateLocation(ctx, &pb.UpdateLocationRequest{Location: *name, Update: &pb.Location{Quota: &pb.Quota{NumOfSlots: int32(*quota)}}})
			}
			if len(*sort) > 0 {
				if *sort == "time" {
					client.UpdateLocation(ctx, &pb.UpdateLocationRequest{Location: *name, Update: &pb.Location{Sort: pb.Location_BY_DATE_ADDED}})
				}
			}
			if !*alert {
				client.UpdateLocation(ctx, &pb.UpdateLocationRequest{Location: *name, Update: &pb.Location{NoAlert: true}})
			}
		}
	}
}
