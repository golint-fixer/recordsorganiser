package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
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

// ByReleaseDate sorts by the given release date
type ByReleaseDate []*pbrc.Record

func (a ByReleaseDate) Len() int      { return len(a) }
func (a ByReleaseDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByReleaseDate) Less(i, j int) bool {
	if a[i].GetRelease().Released != a[j].GetRelease().Released {
		return a[i].GetRelease().Released > a[j].GetRelease().Released
	}
	return strings.Compare(a[i].GetRelease().Title, a[j].GetRelease().Title) < 0
}

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

	if len(recs.GetRecords()) == 0 {
		fmt.Printf("No records with that id\n")
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
	return rel.GetRecords()[0].GetRelease().Title + " [" + strconv.Itoa(int(instanceID)) + "]"
}

func get(ctx context.Context, client pb.OrganiserServiceClient, name string, force bool, slot int32) {
	locs, err := client.GetOrganisation(ctx, &pb.GetOrganisationRequest{ForceReorg: force, Locations: []*pb.Location{&pb.Location{Name: name}}})
	if err != nil {
		log.Fatalf("Error reading locations: %v", err)
	}

	for _, loc := range locs.GetLocations() {
		fmt.Printf("%v (%v) -> %v [%v] with %v (%v) %v\n", loc.GetName(), len(loc.GetReleasesLocation()), loc.GetFolderIds(), loc.GetQuota(), loc.Sort.String(), loc.GetNoAlert(), loc.GetSpillFolder())
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
	case "quota":
		quotaFlags := flag.NewFlagSet("quota", flag.ExitOnError)
		var name = quotaFlags.String("name", "", "The name of the location to evaluate")

		if err := quotaFlags.Parse(os.Args[2:]); err == nil {
			loc, err := client.GetOrganisation(ctx, &pb.GetOrganisationRequest{Locations: []*pb.Location{&pb.Location{Name: *name}}})

			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			fmt.Printf("%v / %v is taken for %v\n", len(loc.GetLocations()[0].ReleasesLocation), loc.GetLocations()[0].Quota.GetNumOfSlots(), *name)

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

			count := 0
			for _, l := range loc.GetLocations() {
				for _, id := range l.GetFolderIds() {
					recs, err := rclient.GetRecords(context.Background(), &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbgd.Release{}, Metadata: &pbrc.ReleaseMetadata{GoalFolder: id}}})
					if err != nil {
						log.Fatalf("Error: %v", err)
					}
					for _, r := range recs.GetRecords() {
						if r.GetMetadata().Category != pbrc.ReleaseMetadata_STAGED && r.GetMetadata().Category != pbrc.ReleaseMetadata_STAGED_TO_SELL && r.GetMetadata().Category != pbrc.ReleaseMetadata_SOLD {
							count++
						}
					}
				}
			}
			fmt.Printf("COUNT = %v\n", count)
			fmt.Printf("We have to sell %v records to meet this\n", count-int(loc.GetLocations()[0].Quota.GetNumOfSlots())+10)
		}
	case "sell":
		sellFlags := flag.NewFlagSet("sell", flag.ExitOnError)
		var name = sellFlags.String("name", "", "The name of the location to get")
		var assess = sellFlags.Bool("assess", false, "Auto assess the for sale records")
		var forcesell = sellFlags.Bool("force", false, "Auto assess the for sale records")

		if err := sellFlags.Parse(os.Args[2:]); err == nil {
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

			loc, err := client.GetOrganisation(ctx, &pb.GetOrganisationRequest{Locations: []*pb.Location{&pb.Location{Name: *name}}})
			if err != nil {
				log.Fatalf("ERRR: %v", err)
			}

			records := make([]*pbrc.Record, 0)
			minScore := int32(6)
			foundOthers := false
			for _, l := range loc.GetLocations() {
				for _, id := range l.GetFolderIds() {
					recs, err := rclient.GetRecords(context.Background(), &pbrc.GetRecordsRequest{Filter: &pbrc.Record{Release: &pbgd.Release{}, Metadata: &pbrc.ReleaseMetadata{GoalFolder: id}}})
					if err != nil {
						log.Fatalf("Error : %v", err)
					}

					for _, r := range recs.GetRecords() {
						if r.GetRelease().Rating > 0 && r.GetMetadata().Category != pbrc.ReleaseMetadata_STAGED && r.GetMetadata().Category != pbrc.ReleaseMetadata_STAGED_TO_SELL && r.GetMetadata().Category != pbrc.ReleaseMetadata_SOLD {
							records = append(records, r)
							if r.GetRelease().Rating < minScore {
								minScore = r.GetRelease().Rating
							}
						}
					}
				}
			}

			for _, r := range records {
				if r.GetMetadata().Others {
					foundOthers = true
					break
				}

			}

			fmt.Printf("FOUND %v [%v] - need to sell %v\n", len(records), minScore, len(records)-int(loc.GetLocations()[0].GetQuota().GetNumOfSlots()))
			if foundOthers {
				fmt.Printf("These have others - auto sell\n")
			}

			total := len(records) - int(loc.GetLocations()[0].GetQuota().GetNumOfSlots()) + 10
			count := 0

			//Sort by release date
			sort.Sort(ByReleaseDate(records))
			for _, r := range records {
				if count > total {
					break
				}
				if r.GetRelease().Rating == minScore && (!foundOthers || r.GetMetadata().Others) {
					count++
					fmt.Printf("SELL: [%v] %v\n", r.GetRelease().InstanceId, r.GetRelease().Title)
					if *assess {
						up := &pbrc.UpdateRecordRequest{Update: &pbrc.Record{Release: &pbgd.Release{InstanceId: r.GetRelease().InstanceId}, Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_ASSESS}}}
						_, err = rclient.UpdateRecord(context.Background(), up)
						if err != nil {
							log.Fatalf("Error updating record: %v", err)
						}
					}
					if *forcesell {
						up := &pbrc.UpdateRecordRequest{Update: &pbrc.Record{Release: &pbgd.Release{InstanceId: r.GetRelease().InstanceId}, Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_PREPARE_TO_SELL}}}
						_, err = rclient.UpdateRecord(context.Background(), up)
						if err != nil {
							log.Fatalf("Error updating record: %v", err)
						}
					}

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
		var spill = updateLocationFlags.Int("spill", 0, "The spill folder for this location")
		var slots = updateLocationFlags.Int("slots", 0, "The new number of slots for this location")
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
			if *spill != 0 {
				client.UpdateLocation(ctx, &pb.UpdateLocationRequest{Location: *name, Update: &pb.Location{SpillFolder: int32(*spill)}})
			}
			if *slots != 0 {
				client.UpdateLocation(ctx, &pb.UpdateLocationRequest{Location: *name, Update: &pb.Location{Slots: int32(*slots)}})
			}
		}
	case "extractor":
		extractFlags := flag.NewFlagSet("Extract", flag.ExitOnError)
		var label = extractFlags.Int("id", -1, "The ID of the label")
		var reg = extractFlags.String("extract", "", "The extractor")
		if err := extractFlags.Parse(os.Args[2:]); err == nil {
			client.AddExtractor(ctx, &pb.AddExtractorRequest{Extractor: &pb.LabelExtractor{LabelId: int32(*label), Extractor: *reg}})
		}
	}

}
