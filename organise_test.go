package main

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordsorganiser/proto"
	"github.com/golang/protobuf/proto"
)

type testBridge struct{}

func (discogsBridge testBridge) getReleases(folders []int32) []*pbd.Release {
	var result []*pbd.Release

	result = append(result, &pbd.Release{
		Id:             1,
		Labels:         []*pbd.Label{&pbd.Label{Name: "FirstLabel"}},
		FormatQuantity: 2,
	})
	result = append(result, &pbd.Release{
		Id:             2,
		Labels:         []*pbd.Label{&pbd.Label{Name: "SecondLabel"}},
		FormatQuantity: 1,
	})
	result = append(result, &pbd.Release{
		Id:             3,
		Labels:         []*pbd.Label{&pbd.Label{Name: "ThirdLabel"}},
		FormatQuantity: 1,
	})

	return result
}

func (discogsBridge testBridge) getRelease(ID int32) *pbd.Release {
	return &pbd.Release{Id: ID}
}

func TestGetReleaseLocation(t *testing.T) {
	testServer := &Server{saveLocation: ".testgetlocation", bridge: testBridge{}, org: &pb.Organisation{}}
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}
	testServer.AddLocation(context.Background(), location)

	relLocation, err := testServer.Locate(context.Background(), &pbd.Release{Id: 2})

	if err != nil {
		t.Errorf("GetReleaseLocation has failed: %v", err)
	}

	if relLocation == nil {
		t.Errorf("Err: %v", relLocation)
	}

	if relLocation.Slot != 2 {
		t.Errorf("Slot has come back wrong: %v", relLocation)
	}

	if relLocation.Before != nil {
		t.Errorf("Release location has come back wrong: %v", relLocation)
	}
	if relLocation.After == nil || relLocation.After.Id != 3 {
		t.Errorf("Release location has come back wrong: %v", relLocation)
	}

	relLocation, _ = testServer.Locate(context.Background(), &pbd.Release{Id: 1})
	if relLocation.Before != nil || relLocation.After != nil {
		t.Errorf("Release location has come back wrong: %v", relLocation)
	}
}

func TestListLocations(t *testing.T) {
	testServer := &Server{saveLocation: ".testgetorgs", bridge: testBridge{}, org: &pb.Organisation{}}
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	testServer.AddLocation(context.Background(), location)

	org, err := testServer.GetOrganisation(context.Background(), &pb.Empty{})

	if err != nil {
		t.Errorf("Error retrieving current organisation")
	}

	if len(org.Locations) != 1 {
		t.Errorf("Too Many Locations: %v", org)
	}
	if org.Locations[0].Name != "TestName" {
		t.Errorf("Location name is incorrect: %v", org.Locations[0])
	}
}

func TestGetOrganisations(t *testing.T) {
	testServer := &Server{saveLocation: ".testgetorgs", bridge: testBridge{}, org: &pb.Organisation{}}
	testServer.save()

	//Sleep for 1.5 seconds to bump the timestamp
	time.Sleep(time.Millisecond * 1500)
	testServer.save()

	organisations, err := testServer.GetOrganisations(context.Background(), &pb.Empty{})
	if err != nil {
		t.Errorf("Get Organisations returned an errors: %v", err)
	}
	if len(organisations.Organisations) != 2 {
		t.Errorf("Organisations has returned wrong: %v", organisations)
	}

	if organisations.Organisations[0].Timestamp <= 0 {
		t.Errorf("Timestamp has come back as wrong: %v", organisations.Organisations[0])
	}

	if organisations.Organisations[0].Timestamp == organisations.Organisations[1].Timestamp {
		t.Errorf("Timestamps are equal between two saves: %v", organisations.Organisations)
	}
}

func TestCompareMoves(t *testing.T) {
	start := []*pb.ReleasePlacement{
		&pb.ReleasePlacement{ReleaseId: 1, Index: 1},
		&pb.ReleasePlacement{ReleaseId: 2, Index: 2},
		&pb.ReleasePlacement{ReleaseId: 3, Index: 3},
	}
	end := []*pb.ReleasePlacement{
		&pb.ReleasePlacement{ReleaseId: 1, Index: 2},
		&pb.ReleasePlacement{ReleaseId: 2, Index: 1},
		&pb.ReleasePlacement{ReleaseId: 4, Index: 3},
	}
	expectedMoves := []*pb.LocationMove{
		&pb.LocationMove{Old: &pb.ReleasePlacement{ReleaseId: 3, Index: 3}},
		&pb.LocationMove{New: &pb.ReleasePlacement{ReleaseId: 4, Index: 3}},
		&pb.LocationMove{Old: &pb.ReleasePlacement{ReleaseId: 1, Index: 1}, New: &pb.ReleasePlacement{ReleaseId: 1, Index: 2}},
	}

	moves := getMoves(start, end)
	if len(moves) != 3 {
		t.Errorf("Not enough moves: %v", moves)
	}
	for i := range moves {
		if !proto.Equal(moves[i], expectedMoves[i]) {
			t.Errorf("Bad move : %v (expected %v)", moves[i], expectedMoves[i])
		}
	}
}

func TestGetLocation(t *testing.T) {
	testServer := &Server{saveLocation: ".testoutget", bridge: testBridge{}, org: &pb.Organisation{}}

	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	testServer.AddLocation(context.Background(), location)

	testServer2 := &Server{saveLocation: ".testoutget", bridge: testBridge{}, org: &pb.Organisation{}}
	testServer2.org = loadLatest(".testoutget")
	retrLocation := &pb.Location{Name: "TestName"}
	retr, err := testServer2.GetLocation(context.Background(), retrLocation)
	if err != nil {
		t.Errorf("Error on getting location: %v", err)
	}

	if retr.Name != "TestName" {
		t.Errorf("Location name is wrong: %v", retrLocation)
	}

	if len(retr.FolderIds) == 0 || retr.FolderIds[0] != 10 {
		t.Errorf("Folder Id has come back wrong: %v", retrLocation)
	}
}

func TestGetLocationFail(t *testing.T) {
	testServer := &Server{saveLocation: ".testoutget", bridge: testBridge{}, org: &pb.Organisation{}}

	location2 := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	testServer.AddLocation(context.Background(), location2)

	testServer2 := &Server{saveLocation: ".testoutget", bridge: testBridge{}, org: &pb.Organisation{}}
	testServer2.org = loadLatest(".testoutget")
	retrLocation := &pb.Location{Name: "MadeUpName"}
	_, err := testServer2.GetLocation(context.Background(), retrLocation)
	if err == nil {
		t.Errorf("Error on getting location: %v", err)
	}
}

func TestAddLocation(t *testing.T) {
	testServer := &Server{saveLocation: ".testaddalocation", bridge: testBridge{}, org: &pb.Organisation{}}

	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	newLocation, err := testServer.AddLocation(context.Background(), location)
	if err != nil {
		t.Errorf("Error on adding location: %v", err)
	}

	if len(newLocation.ReleasesLocation) != 3 {
		t.Errorf("All the releases have not been organised")
	}

	if newLocation.ReleasesLocation[1].Slot != 2 {
		t.Errorf("Second release is in the wrong slot: %v", newLocation.ReleasesLocation[1])
	}

	if newLocation.ReleasesLocation[1].Index != 1 {
		t.Errorf("Second release has the wrong index %v", newLocation)
	}
}
