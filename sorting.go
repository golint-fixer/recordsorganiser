package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	pb "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

// ByDateAdded allows sorting of releases by the date they were added
type ByDateAdded []*pbrc.Record

func (a ByDateAdded) Len() int      { return len(a) }
func (a ByDateAdded) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDateAdded) Less(i, j int) bool {
	if a[i].Metadata.DateAdded != a[j].Metadata.DateAdded {
		return a[i].Metadata.DateAdded < a[j].Metadata.DateAdded
	}
	return strings.Compare(a[i].Release.Title, a[j].Release.Title) < 0
}

// ByLabelCat allows sorting of releases by the date they were added
type ByLabelCat []*pbrc.Record

func (a ByLabelCat) Len() int      { return len(a) }
func (a ByLabelCat) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLabelCat) Less(i, j int) bool {
	return sortByLabelCat(*a[i].GetRelease(), *a[j].GetRelease()) < 0
}

func split(str string) []string {
	return regexp.MustCompile("[0-9]+|[a-z]+|[A-Z]+").FindAllString(str, -1)
}

func sortByLabelCat(rel1, rel2 pb.Release) int {
	label1 := pb.GetMainLabel(rel1.Labels)
	label2 := pb.GetMainLabel(rel2.Labels)

	labelSort := strings.Compare(label1.Name, label2.Name)
	if labelSort != 0 {
		return labelSort
	}

	cat1Elems := split(label1.Catno)
	cat2Elems := split(label2.Catno)

	toCheck := len(cat1Elems)

	for i := 0; i < toCheck; i++ {
		if unicode.IsNumber(rune(cat1Elems[i][0])) && unicode.IsNumber(rune(cat2Elems[i][0])) {
			num1, _ := strconv.Atoi(cat1Elems[i])
			num2, _ := strconv.Atoi(cat2Elems[i])
			if num1 > num2 {
				return 1
			} else if num2 > num1 {
				return -1
			}
		} else {
			catComp := strings.Compare(cat1Elems[i], cat2Elems[i])
			if catComp != 0 {
				return catComp
			}
		}
	}

	//Fallout to sorting by title
	titleComp := strings.Compare(rel1.Title, rel2.Title)
	return titleComp
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

// Split splits a releases list into buckets
func Split(releases []*pbrc.Record, n float64) [][]*pbrc.Record {
	var solution [][]*pbrc.Record

	var count int32
	count = 0
	for _, rel := range releases {
		count += rel.GetRelease().FormatQuantity
	}

	boundaryAccumulator := float64(count) / n
	boundaryValue := boundaryAccumulator
	currentValue := 0.0
	var currentReleases []*pbrc.Record
	for _, rel := range releases {
		if currentValue+float64(rel.GetRelease().FormatQuantity) > boundaryValue {
			solution = append(solution, currentReleases)
			currentReleases = make([]*pbrc.Record, 0)
			boundaryValue += boundaryAccumulator
		}

		currentReleases = append(currentReleases, rel)
		currentValue += float64(rel.GetRelease().FormatQuantity)
	}
	solution = append(solution, currentReleases)

	return solution
}
