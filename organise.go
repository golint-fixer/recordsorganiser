package main

import (
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
}

type discogsBridge interface {
	getReleases(folders []int32) []*pbd.Release
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
	return location, nil
}
