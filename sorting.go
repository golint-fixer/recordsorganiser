package main

import (
	pbd "github.com/brotherlogic/recordsorganiser/proto"
)

// ByDateAdded allows sorting of releases by the date they were added
type ByDateAdded []*pbd.CombinedRelease

func (a ByDateAdded) Len() int           { return len(a) }
func (a ByDateAdded) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDateAdded) Less(i, j int) bool { return a[i].Metadata.DateAdded < a[j].Metadata.DateAdded }
