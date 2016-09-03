package main

import (
	"sort"
	"testing"

	pbs "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

func TestSortByDateAdded(t *testing.T) {
	releases := []*pb.CombinedRelease{
		&pb.CombinedRelease{Release: &pbd.Release{Id: 2}, Metadata: &pbs.ReleaseMetadata{DateAdded: 125}},
		&pb.CombinedRelease{Release: &pbd.Release{Id: 3}, Metadata: &pbs.ReleaseMetadata{DateAdded: 124}},
		&pb.CombinedRelease{Release: &pbd.Release{Id: 4}, Metadata: &pbs.ReleaseMetadata{DateAdded: 123}},
	}

	sort.Sort(ByDateAdded(releases))

	if releases[0].Release.Id != 4 {
		t.Errorf("Releases are not correctly ordered: %v", releases)
	}
}

func TestSortByDateAddedWithFallback(t *testing.T) {
	releases := []*pb.CombinedRelease{
		&pb.CombinedRelease{Release: &pbd.Release{Title: "Second", Id: 2}, Metadata: &pbs.ReleaseMetadata{DateAdded: 124}},
		&pb.CombinedRelease{Release: &pbd.Release{Title: "Third", Id: 3}, Metadata: &pbs.ReleaseMetadata{DateAdded: 124}},
		&pb.CombinedRelease{Release: &pbd.Release{Title: "First", Id: 4}, Metadata: &pbs.ReleaseMetadata{DateAdded: 124}},
	}

	sort.Sort(ByDateAdded(releases))

	if releases[0].Release.Id != 4 {
		t.Errorf("Releases are not correctly ordered: %v", releases)
	}
}

func TestSortByMasterReleaseDate(t *testing.T) {
	releases := []*pbd.Release{
		&pbd.Release{Id: 2, EarliestReleaseDate: 15},
		&pbd.Release{Id: 3, EarliestReleaseDate: 10},
		&pbd.Release{Id: 4, EarliestReleaseDate: 20},
	}

	sort.Sort(ByEarliestReleaseDate(releases))

	if releases[0].Id != 3 {
		t.Errorf("Releases are not correctly ordered: %v", releases)
	}
}
