package main

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/brotherlogic/diffmove"
	"github.com/brotherlogic/goserver"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	pbs "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

// Server the configuration for the syncer
type Server struct {
	*goserver.GoServer
	bridge  discogsBridge
	currOrg *pb.Organisation
	pastOrg *pb.Organisation
}

type discogsBridge interface {
	getReleases(folders []int32) []*pbd.Release
	getRelease(ID int32) (*pbd.Release, error)
	getMetadata(release *pbd.Release) *pbs.ReleaseMetadata
	moveToFolder(releaseMove *pbs.ReleaseMove)
	GetIP(string) (string, int)
}

func getMoves(start []*pb.ReleasePlacement, end []*pb.ReleasePlacement, slot int, folder string) []*pb.LocationMove {
	var moves []*pb.LocationMove

	inStartSlot := 0
	for _, startRec := range start {
		if int(startRec.Slot) == slot && int(startRec.Index) > inStartSlot {
			inStartSlot = int(startRec.Index)
		}
	}
	inEndSlot := 0
	for _, endRec := range end {
		if int(endRec.Slot) == slot && int(endRec.Index) > inEndSlot {
			inEndSlot = int(endRec.Index)
		}
	}

	//Build out the arrays for diffmove
	startNumbers := make([]int, inStartSlot+1)
	endNumbers := make([]int, inEndSlot+1)
	for _, startRec := range start {
		if int(startRec.Slot) == slot {
			startNumbers[startRec.Index-1] = int(startRec.ReleaseId)
		}
	}
	for _, endRec := range end {
		if int(endRec.Slot) == slot {
			endNumbers[endRec.Index-1] = int(endRec.ReleaseId)
		}
	}

	diffMoves := diffmove.Diff(startNumbers, endNumbers)
	for _, move := range diffMoves {
		switch move.Move {
		case "Add":
			moves = append(moves, &pb.LocationMove{SlotMove: false,
				New: &pb.ReleasePlacement{ReleaseId: int32(move.Value), Index: int32(move.Start + 1),
					BeforeReleaseId: int32(move.StartPrior), AfterReleaseId: int32(move.StartFollow), Slot: int32(slot), Folder: folder},
			})
		case "Delete":
			moves = append(moves, &pb.LocationMove{SlotMove: false,
				Old: &pb.ReleasePlacement{ReleaseId: int32(move.Value), Index: int32(move.Start + 1),
					BeforeReleaseId: int32(move.StartPrior), AfterReleaseId: int32(move.StartFollow), Slot: int32(slot), Folder: folder},
			})
		case "Move":
			moves = append(moves, &pb.LocationMove{SlotMove: true,
				Old: &pb.ReleasePlacement{ReleaseId: int32(move.Value), Index: int32(move.Start + 1),
					BeforeReleaseId: int32(move.StartPrior), AfterReleaseId: int32(move.StartFollow), Slot: int32(slot), Folder: folder},
				New: &pb.ReleasePlacement{ReleaseId: int32(move.Value), Index: int32(move.End + 1),
					BeforeReleaseId: int32(move.EndPrior), AfterReleaseId: int32(move.EndFollow), Slot: int32(slot), Folder: folder},
			})
		}
	}

	return moves
}

// Diff computes the diff between two slot organisations
func (s *Server) Diff(ctx context.Context, in *pb.DiffRequest) (*pb.OrganisationMoves, error) {

	var locStart *pb.Location
	var locEnd *pb.Location
	for _, location := range s.currOrg.Locations {
		if location.Name == in.LocationName {
			locStart = location
		}
	}
	for _, location := range s.pastOrg.Locations {
		if location.Name == in.LocationName {
			locEnd = location
		}
	}

	if locEnd == nil || locStart == nil {
		return nil, errors.New("Unable to find location " + in.LocationName)
	}
	moves := getMoves(locStart.ReleasesLocation, locEnd.ReleasesLocation, int(in.Slot), in.LocationName)
	res := &pb.OrganisationMoves{
		StartTimestamp: s.currOrg.Timestamp,
		EndTimestamp:   s.pastOrg.Timestamp,
		Moves:          moves,
	}
	return res, nil
}

// GetOrganisation Gets the current organisation
func (s *Server) GetOrganisation(ctx context.Context, in *pb.Empty) (*pb.Organisation, error) {
	return s.currOrg, nil
}

// Locate gets the location of a given release
func (s *Server) Locate(ctx context.Context, in *pbd.Release) (*pb.ReleaseLocation, error) {
	t := time.Now()
	relLoc := &pb.ReleaseLocation{}
	foundIndex := -1
	for _, loc := range s.currOrg.Locations {
		for _, rel := range loc.ReleasesLocation {
			if rel.ReleaseId == in.Id {
				foundIndex = int(rel.Index)
				relLoc.Location = loc
				relLoc.Slot = rel.Slot
			}
		}
		if foundIndex >= 0 {
			for _, rel := range loc.ReleasesLocation {
				if rel.Slot == relLoc.Slot {
					if int(rel.Index) == foundIndex-1 {
						relLoc.Before, _ = s.bridge.getRelease(rel.ReleaseId)
					}
					if int(rel.Index) == foundIndex+1 {
						relLoc.After, _ = s.bridge.getRelease(rel.ReleaseId)
					}
				}
			}
			s.LogFunction("Locate", t)
			return relLoc, nil
		}
	}
	s.LogFunction("Locate-fail", t)
	return nil, errors.New("Unable to locate record with id " + strconv.Itoa(int(in.Id)))
}

func (s Server) runOrgSteps() {
	s.moveOldRecordsToPile()
}

func (s Server) moveOldRecordsToPile() {
	records := s.bridge.getReleases([]int32{673768})

	for _, record := range records {
		meta := s.bridge.getMetadata(record)
		if meta != nil {
			if meta.DateAdded < (time.Now().AddDate(0, -3, 0).Unix()) {
				if record.Rating > 0 {
					s.bridge.moveToFolder(&pbs.ReleaseMove{Release: record, NewFolderId: 242017})
				} else {
					s.bridge.moveToFolder(&pbs.ReleaseMove{Release: record, NewFolderId: 812802})
				}
			}
		}
	}
}

// DeleteLocation removes a location
func (s *Server) DeleteLocation(ctx context.Context, in *pb.Location) (*pb.Empty, error) {
	for i, folder := range s.currOrg.Locations {
		if folder.GetName() == in.GetName() {
			s.pastOrg = proto.Clone(s.currOrg).(*pb.Organisation)
			s.currOrg.Locations = append(s.currOrg.Locations[:(i)], s.currOrg.Locations[(i)+1:]...)
			s.save()
			return &pb.Empty{}, nil
		}
	}

	return nil, errors.New("Unable to locate " + in.GetName())
}

// Organise Organises out the whole collection
func (s *Server) Organise(ctx context.Context, in *pb.Empty) (*pb.OrganisationMoves, error) {
	s.runOrgSteps()
	newList := &pb.Organisation{}

	for _, folder := range s.currOrg.Locations {
		newList.Locations = append(newList.Locations, s.arrangeLocation(folder))
	}

	diffs := compare(s.currOrg, newList)
	s.pastOrg = s.currOrg
	s.currOrg = newList
	s.save()

	return &pb.OrganisationMoves{StartTimestamp: s.pastOrg.Timestamp, EndTimestamp: s.currOrg.Timestamp, Moves: diffs}, nil
}

// GetLocation Gets an existing location
func (s *Server) GetLocation(ctx context.Context, location *pb.Location) (*pb.Location, error) {
	t := time.Now()
	for _, storedLocation := range s.currOrg.GetLocations() {
		if storedLocation.Name == location.Name {
			s.LogFunction("GetLocation-curr", t)
			return storedLocation, nil
		}
	}

	s.LogFunction("GetLocation-fail", t)
	return &pb.Location{}, errors.New("Cannot find location called " + location.Name)
}

func (s *Server) arrangeLocation(location *pb.Location) *pb.Location {
	releases := s.bridge.getReleases(location.FolderIds)
	retLocation := &pb.Location{Name: location.Name, Units: location.Units, FolderIds: location.FolderIds, Sort: location.Sort, Quota: location.Quota, ExpectedFormat: location.ExpectedFormat, UnexpectedLabel: location.UnexpectedLabel}

	switch location.Sort {
	case pb.Location_BY_LABEL_CATNO:
		sort.Sort(pbd.ByLabelCat(releases))
	case pb.Location_BY_DATE_ADDED:
		var combined []*pb.CombinedRelease
		for _, release := range releases {
			meta := s.bridge.getMetadata(release)
			if meta != nil && meta.DateAdded > 0 {
				comb := &pb.CombinedRelease{Release: release, Metadata: meta}
				combined = append(combined, comb)
			}
		}
		sort.Sort(ByDateAdded(combined))
		newReleases := make([]*pbd.Release, len(releases))
		for i, comb := range combined {
			newReleases[i] = comb.Release
		}
		releases = newReleases
	case pb.Location_BY_RELEASE_DATE:
		sort.Sort(ByEarliestReleaseDate(releases))
	}
	splits := pbd.Split(releases, float64(location.Units))

	var locations []*pb.ReleasePlacement
	for i, split := range splits {
		for j, rel := range split {
			place := &pb.ReleasePlacement{
				ReleaseId: rel.Id,
				Index:     int32(j) + 1,
				Slot:      int32(i) + 1,
			}
			locations = append(locations, place)
		}
	}

	retLocation.ReleasesLocation = locations
	return retLocation
}

// AddLocation Adds a new location to the organiser
func (s *Server) AddLocation(ctx context.Context, location *pb.Location) (*pb.Location, error) {
	newLocation := s.arrangeLocation(location)
	s.currOrg.Locations = append(s.currOrg.Locations, newLocation)
	s.save()
	return newLocation, nil
}

//UpdateLocation updates the location with new properties
func (s *Server) UpdateLocation(ctx context.Context, in *pb.Location) (*pb.Location, error) {
	for i, loc := range s.currOrg.Locations {
		if loc.Name == in.Name {
			s.pastOrg = proto.Clone(s.currOrg).(*pb.Organisation)
			proto.Merge(loc, in)
			newLocation := s.arrangeLocation(loc)
			s.currOrg.Locations[i] = newLocation
			s.save()
			return newLocation, nil
		}
	}
	return nil, errors.New("Cannot find location in org")
}

// GetQuotaViolations gets the quota violations for the whole collection
func (s *Server) GetQuotaViolations(ctx context.Context, in *pb.Empty) (*pb.LocationList, error) {
	t := time.Now()
	violations := &pb.LocationList{}

	for _, location := range s.currOrg.Locations {
		q := location.Quota
		if q.NumOfUnits > 0 && len(location.GetReleasesLocation()) > int(q.NumOfUnits) {
			violations.Locations = append(violations.Locations, location)
		}
	}

	s.LogFunction("GetQuotaViolations", t)
	return violations, nil
}

//CleanLocation reports infractions on a given location
func (s *Server) CleanLocation(ctx context.Context, in *pb.Location) (*pb.CleanList, error) {
	t := time.Now()
	var loc *pb.Location
	for _, l := range s.currOrg.GetLocations() {
		if l.Name == in.Name {
			loc = l
		}
	}

	if loc == nil {
		s.LogFunction("CleanLocation-fail", t)
		return nil, errors.New("Unable to find location " + in.Name)
	}

	list := &pb.CleanList{}
	for _, entry := range loc.ReleasesLocation {
		record, err := s.bridge.getRelease(entry.ReleaseId)
		if record == nil {
			s.Log(fmt.Sprintf("Checking %v from %v given %v", record, entry, err))
		}
		match := false
		for _, format := range record.GetFormats() {
			for _, desc := range format.Descriptions {
				m, _ := regexp.MatchString(loc.ExpectedFormat, desc)
				if m {
					match = true
				}
			}
		}

		badlabel := false
		if len(loc.UnexpectedLabel) > 0 {
			for _, label := range record.GetLabels() {
				m, _ := regexp.MatchString(loc.UnexpectedLabel, label.Name)
				if m {
					badlabel = true
				}
			}
		}

		if !match || badlabel {
			list.Entries = append(list.Entries, record)
		}
	}

	s.LogFunction("CleanLocation", t)
	return list, nil
}
