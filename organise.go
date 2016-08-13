package main

import (
	"errors"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/brotherlogic/diffmove"
	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"

	pbd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

// Server the configuration for the syncer
type Server struct {
	*goserver.GoServer
	saveLocation string
	bridge       discogsBridge
	org          pb.Organisation
}

type discogsBridge interface {
	getReleases(folders []int32) []*pbd.Release
}

func getMoves(start []*pb.ReleasePlacement, end []*pb.ReleasePlacement) []*pb.LocationMove {
	var moves []*pb.LocationMove

	//Build out the arrays for diffmove
	startNumbers := make([]int, len(start))
	endNumbers := make([]int, len(end))
	for _, startRec := range start {
		startNumbers[startRec.Index-1] = int(startRec.ReleaseId)
	}
	for _, endRec := range end {
		endNumbers[endRec.Index-1] = int(endRec.ReleaseId)
	}
	log.Printf("DIFFING %v -> %v", startNumbers, endNumbers)
	diffMoves := diffmove.Diff(startNumbers, endNumbers)
	log.Printf("RES: %v", diffMoves)
	for _, move := range diffMoves {
		switch move.Move {
		case "Add":
			moves = append(moves, &pb.LocationMove{
				New: &pb.ReleasePlacement{ReleaseId: int32(move.Value), Index: int32(move.Start + 1)},
			})
		case "Delete":
			moves = append(moves, &pb.LocationMove{
				Old: &pb.ReleasePlacement{ReleaseId: int32(move.Value), Index: int32(move.Start + 1)},
			})
			log.Printf("WHAT = %v from %v", moves, move)
		case "Move":
			moves = append(moves, &pb.LocationMove{
				Old: &pb.ReleasePlacement{ReleaseId: int32(move.Value), Index: int32(move.Start + 1)},
				New: &pb.ReleasePlacement{ReleaseId: int32(move.Value), Index: int32(move.End + 1)},
			})
		}
	}

	return moves
}

// GetOrganisation Gets the current organisation
func (s *Server) GetOrganisation(ctx context.Context, in *pb.Empty) (*pb.Organisation, error) {
	return &s.org, nil
}

// GetOrganisations Gets all the available organisations
func (s *Server) GetOrganisations(ctx context.Context, in *pb.Empty) (*pb.OrganisationList, error) {
	orgList := &pb.OrganisationList{}
	files, _ := ioutil.ReadDir(s.saveLocation)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".data") {
			org, _ := load(s.saveLocation, file.Name()[0:len(file.Name())-5])
			orgList.Organisations = append(orgList.Organisations, org)
		}
	}

	return orgList, nil
}

// GetLocation Gets an existing location
func (s *Server) GetLocation(ctx context.Context, location *pb.Location) (*pb.Location, error) {
	for _, storedLocation := range s.org.GetLocations() {
		if storedLocation.Name == location.Name {
			return storedLocation, nil
		}
	}

	return &pb.Location{}, errors.New("Cannot find location called " + location.Name)
}

// AddLocation Adds a new location to the organiser
func (s *Server) AddLocation(ctx context.Context, location *pb.Location) (*pb.Location, error) {
	var locations []*pb.ReleasePlacement
	releases := s.bridge.getReleases(location.FolderIds)

	sort.Sort(pbd.ByLabelCat(releases))
	splits := pbd.Split(releases, float64(location.Units))

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

	location.ReleasesLocation = locations

	s.org.Locations = append(s.org.Locations, location)
	s.save()

	return location, nil
}
