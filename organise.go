package main

import (
	"errors"
	"sort"

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
	org          *pb.Organisation
}

type discogsBridge interface {
	getReleases(folders []int32) []*pbd.Release
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
