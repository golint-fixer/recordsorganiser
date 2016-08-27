package main

import (
	pb "github.com/brotherlogic/godiscogs"
	pbd "github.com/brotherlogic/recordsorganiser/proto"
)

// ByDateAdded allows sorting of releases by the date they were added
type ByDateAdded []*pbd.CombinedRelease

func (a ByDateAdded) Len() int           { return len(a) }
func (a ByDateAdded) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDateAdded) Less(i, j int) bool { return a[i].Metadata.DateAdded < a[j].Metadata.DateAdded }

// ByEarliestReleaseDate allows sorting by the earliest release date
type ByEarliestReleaseDate []*pb.Release

func (a ByEarliestReleaseDate) Len() int      { return len(a) }
func (a ByEarliestReleaseDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByEarliestReleaseDate) Less(i, j int) bool {
	return a[i].EarliestReleaseDate < a[j].EarliestReleaseDate
}
