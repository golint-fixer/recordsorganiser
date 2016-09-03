package main

import (
	"strings"

	pb "github.com/brotherlogic/godiscogs"
	pbd "github.com/brotherlogic/recordsorganiser/proto"
)

// ByDateAdded allows sorting of releases by the date they were added
type ByDateAdded []*pbd.CombinedRelease

func (a ByDateAdded) Len() int      { return len(a) }
func (a ByDateAdded) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDateAdded) Less(i, j int) bool {
	if a[i].Metadata.DateAdded != a[j].Metadata.DateAdded {
		return a[i].Metadata.DateAdded < a[j].Metadata.DateAdded
	}
	return strings.Compare(a[i].Release.Title, a[j].Release.Title) < 0
}

// ByEarliestReleaseDate allows sorting by the earliest release date
type ByEarliestReleaseDate []*pb.Release

func (a ByEarliestReleaseDate) Len() int      { return len(a) }
func (a ByEarliestReleaseDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByEarliestReleaseDate) Less(i, j int) bool {
	if a[i].EarliestReleaseDate != a[j].EarliestReleaseDate {
		return a[i].EarliestReleaseDate < a[j].EarliestReleaseDate
	}
	return strings.Compare(a[i].Title, a[j].Title) < 0
}
