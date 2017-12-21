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

func (s *Server) organise(c *pb.Organisation) (int32, error) {
	num := int32(0)
	for _, l := range s.org.Locations {
		n, err := s.organiseLocation(l)
		if err != nil {
			return -1, err
		}
		num += n
	}
	return num, nil
}

func (s *Server) organiseLocation(c *pb.Location) (int32, error) {
	fr, err := s.bridge.getReleases(c.GetFolderIds())
	if err != nil {
		return -1, err
	}

	switch c.GetSort() {
	case pb.Location_BY_DATE_ADDED:
		sort.Sort(ByLabelCat(fr))
	}

	return int32(len(fr)), nil
}
