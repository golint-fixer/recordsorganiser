package main

import (
	"sort"

	"github.com/brotherlogic/goserver"

	pbs "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

// Server the configuration for the syncer
type Server struct {
	*goserver.GoServer
	bridge discogsBridge
	org    *pb.Organisation
}

type discogsBridge interface {
	getReleases(folders []int32) ([]*pbrc.Record, error)
	getMetadata(release *pbd.Release) (*pbs.ReleaseMetadata, error)
	moveToFolder(releaseMove *pbs.ReleaseMove)
	GetIP(string) (string, int)
}

func (s *Server) prepareForReorg() {

}

func (s *Server) organise(c *pb.Organisation) error {
	for _, l := range s.org.Locations {
		err := s.organiseLocation(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) organiseLocation(c *pb.Location) error {
	fr, err := s.bridge.getReleases(c.GetFolderIds())
	if err != nil {
		return err
	}

	switch c.GetSort() {
	case pb.Location_BY_DATE_ADDED:
		sort.Sort(ByLabelCat(fr))
	}

	return nil
}
