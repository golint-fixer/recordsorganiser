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

func TestGetOverQuota(t *testing.T) {
	testServer := getTestServer(".testQuota")
	location := &pb.Location{
		Name:      "TestName",
		Slots:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_DATE_ADDED,
		Quota:     &pb.Quota{NumOfSlots: 2},
	}

	l, err := testServer.AddLocation(context.Background(), &pb.AddLocationRequest{Add: location})
	if err != nil {
		t.Fatalf("Error in adding location: %v", err)
	}

	if len(l.GetNow().GetLocations()[0].GetReleasesLocation()) == 0 {
		t.Fatalf("No releases at the new location")
	}

	quota, err := testServer.GetQuota(context.Background(), &pb.QuotaRequest{FolderId: 10})

	if err != nil {
		t.Fatalf("Error getting quota: %v", err)
	}

	if !quota.GetOverQuota() {
		t.Errorf("Reported under quota?: %v", quota)
	}
}

func TestGetUnderQuota(t *testing.T) {
	testServer := getTestServer(".testQuota")
	location := &pb.Location{
		Name:      "TestName",
		Slots:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_DATE_ADDED,
		Quota:     &pb.Quota{NumOfSlots: 4},
	}

	l, err := testServer.AddLocation(context.Background(), &pb.AddLocationRequest{Add: location})
	if err != nil {
		t.Fatalf("Error in adding location: %v", err)
	}

	if len(l.GetNow().GetLocations()[0].GetReleasesLocation()) == 0 {
		t.Fatalf("No releases at the new location")
	}

	quota, err := testServer.GetQuota(context.Background(), &pb.QuotaRequest{FolderId: 10})

	if err != nil {
		t.Fatalf("Error getting quota: %v", err)
	}

	if quota.GetOverQuota() {
		t.Errorf("Reported over quota?: %v", quota)
	}
}

func TestGetQuotaFail(t *testing.T) {
	testServer := getTestServer(".testQuota")
	location := &pb.Location{
		Name:      "TestName",
		Slots:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_DATE_ADDED,
		Quota:     &pb.Quota{NumOfSlots: 4},
	}

	l, err := testServer.AddLocation(context.Background(), &pb.AddLocationRequest{Add: location})
	if err != nil {
		t.Fatalf("Error in adding location: %v", err)
	}

	if len(l.GetNow().GetLocations()[0].GetReleasesLocation()) == 0 {
		t.Fatalf("No releases at the new location")
	}

	_, err = testServer.GetQuota(context.Background(), &pb.QuotaRequest{FolderId: 20})

	if err == nil {
		t.Errorf("No errror on bad quota: %v", err)
	}
}
