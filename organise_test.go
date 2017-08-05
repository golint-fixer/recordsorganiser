package main

import (
	"log"
	"testing"
	"time"

	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"

	pbs "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	"github.com/brotherlogic/goserver"
	pb "github.com/brotherlogic/recordsorganiser/proto"
	"github.com/golang/protobuf/proto"
)

type testBridge struct{}

type testBridgeMove struct {
	move bool
}

func (discogsBridge testBridge) getMetadata(rel *pbd.Release) *pbs.ReleaseMetadata {
	metadata := &pbs.ReleaseMetadata{}
	switch rel.Id {
	case 1:
		metadata.DateAdded = time.Now().Unix()
	case 2:
		metadata.DateAdded = time.Now().Unix() - 100
	case 3:
		metadata.DateAdded = time.Now().Unix() + 100
	}
	return metadata
}

func (discogsBridge testBridgeMove) getMetadata(rel *pbd.Release) *pbs.ReleaseMetadata {
	metadata := &pbs.ReleaseMetadata{}
	switch rel.Id {
	case 1:
		metadata.DateAdded = time.Now().Unix()
	case 2:
		metadata.DateAdded = time.Now().AddDate(0, -4, 0).Unix()
	case 3:
		metadata.DateAdded = time.Now().AddDate(0, -4, 0).Unix()
	}
	return metadata
}

func (discogsBridge testBridge) getReleases(folders []int32) []*pbd.Release {
	var result []*pbd.Release

	result = append(result, &pbd.Release{
		Id:             1,
		Labels:         []*pbd.Label{&pbd.Label{Name: "FirstLabel"}},
		Formats:        []*pbd.Format{&pbd.Format{Name: "12"}},
		FormatQuantity: 2,
	})
	result = append(result, &pbd.Release{
		Id:             2,
		Labels:         []*pbd.Label{&pbd.Label{Name: "SecondLabel"}},
		Formats:        []*pbd.Format{&pbd.Format{Name: "12"}},
		FormatQuantity: 1,
	})
	result = append(result, &pbd.Release{
		Id:             3,
		Labels:         []*pbd.Label{&pbd.Label{Name: "ThirdLabel"}},
		Formats:        []*pbd.Format{&pbd.Format{Name: "CD"}},
		FormatQuantity: 1,
	})

	return result
}

func (discogsBridge testBridgeMove) getReleases(folders []int32) []*pbd.Release {
	var result []*pbd.Release

	result = append(result, &pbd.Release{
		Id:             1,
		Labels:         []*pbd.Label{&pbd.Label{Name: "FirstLabel"}},
		Formats:        []*pbd.Format{&pbd.Format{Name: "12"}},
		FormatQuantity: 2,
	})
	result = append(result, &pbd.Release{
		Id:             2,
		Labels:         []*pbd.Label{&pbd.Label{Name: "SecondLabel"}},
		Formats:        []*pbd.Format{&pbd.Format{Name: "12"}},
		FormatQuantity: 1,
	})
	result = append(result, &pbd.Release{
		Id:             3,
		Labels:         []*pbd.Label{&pbd.Label{Name: "FourthLabel"}},
		Formats:        []*pbd.Format{&pbd.Format{Name: "12"}},
		FormatQuantity: 1,
	})
	result = append(result, &pbd.Release{
		Id:             4,
		Labels:         []*pbd.Label{&pbd.Label{Name: "ThirdLabel"}},
		Formats:        []*pbd.Format{&pbd.Format{Name: "CD"}},
		FormatQuantity: 1,
		Rating:         5,
	})

	return result
}

func (discogsBridge testBridge) getRelease(ID int32) *pbd.Release {
	if ID < 3 {
		return &pbd.Release{Id: ID, Formats: []*pbd.Format{&pbd.Format{Name: "12"}}, Labels: []*pbd.Label{&pbd.Label{Name: "SomethingElse"}}}
	}
	return &pbd.Release{Id: ID, Formats: []*pbd.Format{&pbd.Format{Name: "CD"}}, Labels: []*pbd.Label{&pbd.Label{Name: "Numero"}}}
}

func (discogsBridge testBridgeMove) getRelease(ID int32) *pbd.Release {
	if ID < 4 {
		return &pbd.Release{Id: ID, Formats: []*pbd.Format{&pbd.Format{Name: "12"}}, Labels: []*pbd.Label{&pbd.Label{Name: "SomethingElse"}}}
	}
	return &pbd.Release{Id: ID, Formats: []*pbd.Format{&pbd.Format{Name: "CD"}}, Labels: []*pbd.Label{&pbd.Label{Name: "Numero"}}}
}

func (discogsBridge testBridge) moveToFolder(move *pbs.ReleaseMove) {
	//Do nothing
}

func (discogsBridge testBridgeMove) moveToFolder(move *pbs.ReleaseMove) {
	//Do nothing
}

func getTestServer(dir string) *Server {
	testServer := &Server{GoServer: &goserver.GoServer{}, bridge: testBridge{}, currOrg: &pb.Organisation{}}
	testServer.Register = testServer
	testServer.GoServer.KSclient = *keystoreclient.GetTestClient(dir)
	return testServer
}

func getTestServerWithMove(dir string) *Server {
	testServer := &Server{GoServer: &goserver.GoServer{}, bridge: testBridgeMove{}, currOrg: &pb.Organisation{}}
	testServer.Register = testServer
	testServer.GoServer.KSclient = *keystoreclient.GetTestClient(dir)
	return testServer
}

func TestDeleteNonLocation(t *testing.T) {
	testServer := getTestServer(".testgetreleaselocation")

	_, err := testServer.DeleteLocation(context.Background(), &pb.Location{Name: "Made Up"})

	if err == nil {
		t.Errorf("Delete location has not returned an error")
	}
}

func TestDeleteLocation(t *testing.T) {
	testServer := getTestServer(".testgetreleaselocation")
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}
	testServer.AddLocation(context.Background(), location)

	locations, err := testServer.GetOrganisation(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Unable to pull locations %v", err)
	}
	if len(locations.GetLocations()) != 1 {
		t.Errorf("Location has not been added")
	}

	testServer.DeleteLocation(context.Background(), &pb.Location{Name: "TestName"})

	locations, err = testServer.GetOrganisation(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Unable to pull locations %v", err)
	}
	if len(locations.GetLocations()) != 0 {
		t.Errorf("Location has not been deleted")
	}
}

func TestGetReleaseLocation(t *testing.T) {
	testServer := getTestServer(".testgetreleaselocation")
	log.Printf("HERE = %v", testServer)
	testServer.SkipLog = true
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
	testServer := getTestServer(".testlistlocations")
	testServer.SkipLog = true
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

func TestCompareMoves(t *testing.T) {
	start := []*pb.ReleasePlacement{
		&pb.ReleasePlacement{ReleaseId: 1, Index: 1, Slot: 1},
		&pb.ReleasePlacement{ReleaseId: 2, Index: 2, Slot: 1},
		&pb.ReleasePlacement{ReleaseId: 3, Index: 3, Slot: 1},
	}
	end := []*pb.ReleasePlacement{
		&pb.ReleasePlacement{ReleaseId: 1, Index: 2, Slot: 1},
		&pb.ReleasePlacement{ReleaseId: 2, Index: 1, Slot: 1},
		&pb.ReleasePlacement{ReleaseId: 4, Index: 3, Slot: 1},
	}
	expectedMoves := []*pb.LocationMove{
		&pb.LocationMove{SlotMove: false, Old: &pb.ReleasePlacement{ReleaseId: 3, Index: 3, BeforeReleaseId: 2, Slot: 1, Folder: "MadeUp"}},
		&pb.LocationMove{SlotMove: false, New: &pb.ReleasePlacement{ReleaseId: 4, Index: 3, BeforeReleaseId: 2, Slot: 1, Folder: "MadeUp"}},
		&pb.LocationMove{SlotMove: true, Old: &pb.ReleasePlacement{ReleaseId: 1, Index: 1, AfterReleaseId: 2, Slot: 1, Folder: "MadeUp"}, New: &pb.ReleasePlacement{ReleaseId: 1, Index: 2, BeforeReleaseId: 2, AfterReleaseId: 4, Slot: 1, Folder: "MadeUp"}},
	}

	moves := getMoves(start, end, 1, "MadeUp")
	if len(moves) != 3 {
		t.Errorf("Not enough moves: %v", moves)
	}
	for i := range moves {
		if !proto.Equal(moves[i], expectedMoves[i]) {
			t.Errorf("Bad move : %v (expected %v)", moves[i], expectedMoves[i])
		}
	}
}

func TestIdentifySlotMovesCorrectly(t *testing.T) {
	start := &pb.Organisation{
		Timestamp: 1,
		Locations: []*pb.Location{
			&pb.Location{
				Name:  "Test",
				Units: 2,
				ReleasesLocation: []*pb.ReleasePlacement{
					&pb.ReleasePlacement{ReleaseId: 1, Index: 1, Slot: 1},
				},
			},
		},
	}

	end := &pb.Organisation{
		Timestamp: 1,
		Locations: []*pb.Location{
			&pb.Location{
				Name:  "Test",
				Units: 2,
				ReleasesLocation: []*pb.ReleasePlacement{
					&pb.ReleasePlacement{ReleaseId: 1, Index: 1, Slot: 2},
				},
			},
		},
	}

	expectedMoves := []*pb.LocationMove{
		&pb.LocationMove{SlotMove: true, Old: &pb.ReleasePlacement{ReleaseId: 1, Index: 1, Slot: 1, Folder: "Test"}},
		&pb.LocationMove{SlotMove: true, New: &pb.ReleasePlacement{ReleaseId: 1, Index: 1, Slot: 2, Folder: "Test"}},
	}

	moves := compare(start, end)
	if len(moves) != len(expectedMoves) {
		t.Errorf("Not enough moves (%v vs %v): %v", len(moves), len(expectedMoves), moves)
	} else {
		for i := range moves {
			if !proto.Equal(moves[i], expectedMoves[i]) {
				t.Errorf("Bad move : %v (expected %v)", moves[i], expectedMoves[i])
			}
		}
	}
}

func TestCleanLocation(t *testing.T) {
	testServer := getTestServer(".testCleanLocation")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:           "TestName",
		Units:          2,
		FolderIds:      []int32{10},
		Sort:           pb.Location_BY_LABEL_CATNO,
		ExpectedFormat: "12",
	}
	testServer.AddLocation(context.Background(), location)

	cleans, err := testServer.CleanLocation(context.Background(), location)
	if err != nil || len(cleans.Entries) != 1 || cleans.Entries[0].Id != 3 {
		t.Errorf("Cleaning error: %v, %v", cleans, err)
	}
}

func TestCleanLocationFail(t *testing.T) {
	testServer := getTestServer(".testCleanLocation")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:           "TestName",
		Units:          2,
		FolderIds:      []int32{10},
		Sort:           pb.Location_BY_LABEL_CATNO,
		ExpectedFormat: "12",
	}
	testServer.AddLocation(context.Background(), location)

	_, err := testServer.CleanLocation(context.Background(), &pb.Location{Name: "MadeUpName"})
	if err == nil {
		t.Errorf("No cleaning error: %v", err)
	}
}

func TestCleanLocationOfLabels(t *testing.T) {
	testServer := getTestServer(".testCleanLocationOfLabels")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:            "TestName",
		Units:           2,
		FolderIds:       []int32{10},
		Sort:            pb.Location_BY_LABEL_CATNO,
		UnexpectedLabel: "Numero",
	}
	testServer.AddLocation(context.Background(), location)

	cleans, err := testServer.CleanLocation(context.Background(), location)
	if err != nil || len(cleans.Entries) != 1 || cleans.Entries[0].Id != 3 {
		t.Errorf("Cleaning error: %v, %v", cleans, err)
	}
}

func TestStraightClean(t *testing.T) {
	testServer := getTestServer(".testStraightClean")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}
	testServer.AddLocation(context.Background(), location)

	cleans, err := testServer.CleanLocation(context.Background(), location)
	if err != nil || len(cleans.Entries) != 0 {
		t.Errorf("Cleaning error: %v, %v", cleans, err)
	}
}

func TestQuotaFail(t *testing.T) {
	testServer := getTestServer(".testQuotaFail")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
		Quota:     2,
	}

	testServer.AddLocation(context.Background(), location)
	violations, err := testServer.GetQuotaViolations(context.Background(), &pb.Empty{})
	if err != nil {
		t.Errorf("Error in getting quota: %v", err)
	}
	if len(violations.Locations) != 1 || violations.Locations[0].Name != "TestName" {
		t.Errorf("Violations are not returned correctly: %v", violations)
	}
}

func TestAddLocationByReleaseDate(t *testing.T) {
	testServer := getTestServer(".testAddLocation")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_RELEASE_DATE,
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

func TestAddLocation(t *testing.T) {
	testServer := getTestServer(".testAddLocation")
	testServer.SkipLog = true
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

func TestGetLocation(t *testing.T) {
	testServer := getTestServer(".testGetLocation")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	_, err := testServer.AddLocation(context.Background(), location)
	if err != nil {
		t.Errorf("Error on adding location: %v", err)
	}

	loc, err := testServer.GetLocation(context.Background(), location)
	if err != nil {
		t.Errorf("Error getting location: %v", err)
	}

	if len(loc.ReleasesLocation) != 3 {
		t.Errorf("All the releases have not been organised")
	}

	if loc.ReleasesLocation[1].Slot != 2 {
		t.Errorf("Second release is in the wrong slot: %v", loc.ReleasesLocation[1])
	}

	if loc.ReleasesLocation[1].Index != 1 {
		t.Errorf("Second release has the wrong index %v", loc)
	}
}

func TestGetLocationFail(t *testing.T) {
	testServer := getTestServer(".testGetLocation")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	_, err := testServer.AddLocation(context.Background(), location)
	if err != nil {
		t.Errorf("Error on adding location: %v", err)
	}

	_, err = testServer.GetLocation(context.Background(), &pb.Location{Name: "MadeUp"})
	if err == nil {
		t.Errorf("Location pull has not failed: %v", err)
	}
}

func TestOverallOrg(t *testing.T) {
	testServer := getTestServer(".testOverallOrg")
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}
	testServer.AddLocation(context.Background(), location)

	moves, err := testServer.Organise(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Error in performing org")
	}

	if len(moves.GetMoves()) != 0 {
		t.Errorf("Wrong nubmer of moves (%v)", len(moves.GetMoves()))
	}
}

func TestRateRecordInPile(t *testing.T) {
	testServer := getTestServerWithMove(".testOverallOrg")

	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}
	testServer.AddLocation(context.Background(), location)

	moves, err := testServer.Organise(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Error in performing org")
	}

	if len(moves.GetMoves()) != 0 {
		t.Errorf("Wrong nubmer of moves (%v)", len(moves.GetMoves()))
	}
}

func TestUpdateLocationFail(t *testing.T) {
	testServer := getTestServer(".testdiff")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	testServer.AddLocation(context.Background(), location)
	locationUpdate := &pb.Location{
		Sort: pb.Location_BY_DATE_ADDED,
		Name: "MadeuPName",
	}
	//Wait 2 seconds to let the timestamps change
	time.Sleep(time.Second * 2)
	_, err := testServer.UpdateLocation(context.Background(), locationUpdate)
	if err == nil {
		t.Errorf("Update locaiton has not failed")
	}
}
func TestDiff(t *testing.T) {
	testServer := getTestServer(".testdiff")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	testServer.AddLocation(context.Background(), location)
	locationUpdate := &pb.Location{
		Sort: pb.Location_BY_DATE_ADDED,
		Name: "TestName",
	}
	//Wait 2 seconds to let the timestamps change
	time.Sleep(time.Second * 2)
	testServer.UpdateLocation(context.Background(), locationUpdate)

	diffRequest := &pb.DiffRequest{
		LocationName: "TestName",
		Slot:         1,
	}
	moves, err := testServer.Diff(context.Background(), diffRequest)
	if err != nil {
		t.Errorf("Error running diff %v", err)
	}
	if len(moves.Moves) != 2 {
		t.Errorf("Moves are wrong on diff: %v", moves)
	}

	//Validate the move
	log.Printf("MOVE = %v", moves)
	if moves.Moves[0].Old.Slot != 1 {
		t.Errorf("Slot rep is wrong")
	}
	if moves.Moves[0].Old.Folder != "TestName" {
		t.Errorf("Folder rep is wrong")
	}
}

func TestDiffFail(t *testing.T) {
	testServer := getTestServer(".testdiff")
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     2,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}

	testServer.AddLocation(context.Background(), location)
	locationUpdate := &pb.Location{
		Sort: pb.Location_BY_DATE_ADDED,
		Name: "TestName",
	}
	//Wait 2 seconds to let the timestamps change
	time.Sleep(time.Second * 2)
	testServer.UpdateLocation(context.Background(), locationUpdate)

	diffRequest := &pb.DiffRequest{
		LocationName: "MadeUp",
		Slot:         1,
	}
	_, err := testServer.Diff(context.Background(), diffRequest)
	if err == nil {
		t.Errorf("No error when running diff %v", err)
	}
}

func TestGetReleaseLocationFull(t *testing.T) {
	testServer := getTestServerWithMove(".testgetreleaselocation")
	log.Printf("HERE = %v", testServer)
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     1,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}
	testServer.AddLocation(context.Background(), location)

	relLocation, err := testServer.Locate(context.Background(), &pbd.Release{Id: 3})

	if err != nil {
		t.Errorf("GetReleaseLocation has failed: %v", err)
	}

	if relLocation == nil {
		t.Errorf("Err: %v", relLocation)
	}

	if relLocation.Slot != 1 {
		t.Errorf("Slot has come back wrong: %v", relLocation)
	}

	if relLocation.Before == nil || relLocation.Before.Id != 1 {
		t.Errorf("Release location has come back wrong: %v", relLocation)
	}
	if relLocation.After == nil || relLocation.After.Id != 2 {
		t.Errorf("Release location has come back wrong: %v", relLocation)
	}
}

func TestGetReleaseLocationFail(t *testing.T) {
	testServer := getTestServerWithMove(".testgetreleaselocation")
	log.Printf("HERE = %v", testServer)
	testServer.SkipLog = true
	location := &pb.Location{
		Name:      "TestName",
		Units:     1,
		FolderIds: []int32{10},
		Sort:      pb.Location_BY_LABEL_CATNO,
	}
	testServer.AddLocation(context.Background(), location)

	_, err := testServer.Locate(context.Background(), &pbd.Release{Id: 200})

	if err == nil {
		t.Errorf("GetReleaseLocation has not failed: %v", err)
	}
}
