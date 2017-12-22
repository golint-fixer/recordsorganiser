package main

import (
	"context"
	"testing"

	pb "github.com/brotherlogic/recordsorganiser/proto"
)

func TestLocate(t *testing.T) {
	testServer := getTestServer(".testLocate")
	location := &pb.Location{
		Name:      "TestName",
		Slots:     1,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_DATE_ADDED,
		ReleasesLocation: []*pb.ReleasePlacement{
			&pb.ReleasePlacement{InstanceId: 1234, Index: 1, Slot: 1},
		},
	}
	testServer.org.Locations = append(testServer.org.Locations, location)

	f, err := testServer.Locate(context.Background(), &pb.LocateRequest{InstanceId: 1234})

	if err != nil {
		t.Fatalf("Error locating record: %v", err)
	}

	if f.FoundLocation.GetName() != "TestName" {
		t.Errorf("Error on spotted location: %v", f)
	}
}

func TestLocateFail(t *testing.T) {
	testServer := getTestServer(".testLocate")
	location := &pb.Location{
		Name:      "TestName",
		Slots:     1,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_DATE_ADDED,
		ReleasesLocation: []*pb.ReleasePlacement{
			&pb.ReleasePlacement{InstanceId: 1234, Index: 1, Slot: 1},
		},
	}
	testServer.org.Locations = append(testServer.org.Locations, location)

	f, err := testServer.Locate(context.Background(), &pb.LocateRequest{InstanceId: 12345})

	if err == nil {
		t.Fatalf("Failed locate has not failed: %v", f)
	}
}

func TestGetLocation(t *testing.T) {
	testServer := getTestServer(".testAddLocation")
	location := &pb.Location{
		Name:      "TestName",
		Slots:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_DATE_ADDED,
	}

	_, err := testServer.AddLocation(context.Background(), &pb.AddLocationRequest{Add: location})
	if err != nil {
		t.Fatalf("Unable to add location: %v", err)
	}

	resp, err := testServer.GetOrganisation(context.Background(), &pb.GetOrganisationRequest{ForceReorg: true, Locations: []*pb.Location{&pb.Location{Name: "TestName"}}})
	if err != nil {
		t.Fatalf("Unable to get organisation %v", err)
	}

	if len(resp.GetLocations()) != 1 {
		t.Errorf("Bad location response: %v", resp)
	}
}
func TestGetLocationOrgFail(t *testing.T) {
	testServer := getTestServer(".testAddLocation")
	location := &pb.Location{
		Name:      "TestName",
		Slots:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_DATE_ADDED,
	}

	_, err := testServer.AddLocation(context.Background(), &pb.AddLocationRequest{Add: location})
	if err != nil {
		t.Fatalf("Unable to add location: %v", err)
	}

	testServer.bridge = testBridgeFail{}

	_, err = testServer.GetOrganisation(context.Background(), &pb.GetOrganisationRequest{ForceReorg: true, Locations: []*pb.Location{&pb.Location{Name: "TestName"}}})
	if err == nil {
		t.Fatalf("Failing bridge did not fail reorg %v", err)
	}

}
