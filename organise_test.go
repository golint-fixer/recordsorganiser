package main

import (
	"testing"

	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

type testBridge struct{}

func (discogsBridge testBridge) getReleases(folders []int32) []*pbd.Release {
	var result []*pbd.Release

	result = append(result, &pbd.Release{
		Labels:         []*pbd.Label{&pbd.Label{Name: "FirstLabel"}},
		FormatQuantity: 2,
	})
	result = append(result, &pbd.Release{
		Labels:         []*pbd.Label{&pbd.Label{Name: "SecondLabel"}},
		FormatQuantity: 1,
	})
	result = append(result, &pbd.Release{
		Labels:         []*pbd.Label{&pbd.Label{Name: "ThirdLabel"}},
		FormatQuantity: 1,
	})

	return result
}

func TestAddLocation(t *testing.T) {
	testServer := &Server{saveLocation: ".testout", bridge: testBridge{}}

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
}
