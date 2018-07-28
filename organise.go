package main

import (
	"sort"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"

	pbs "github.com/brotherlogic/discogssyncer/server"
	pbd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordsorganiser/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

// Server the configuration for the syncer
type Server struct {
	*goserver.GoServer
	bridge        discogsBridge
	org           *pb.Organisation
	gh            gh
	lastOrgTime   time.Duration
	lastOrgFolder string
}

type gh interface {
	alert(r *pb.Location) error
}

type discogsBridge interface {
	getReleases(folders []int32) ([]*pbrc.Record, error)
	getMetadata(release *pbd.Release) (*pbrc.ReleaseMetadata, error)
	moveToFolder(releaseMove *pbs.ReleaseMove)
	GetIP(string) (string, int)
}

func (s *Server) prepareForReorg() {

}

func (s *Server) organise(c *pb.Organisation) (int32, error) {
	num := int32(0)
	for _, l := range s.org.Locations {
		n, err := s.organiseLocation(context.Background(), l)
		if err != nil {
			return -1, err
		}
		num += n
	}
	return num, nil
}

func convert(exs []*pb.LabelExtractor) map[int32]string {
	m := make(map[int32]string)
	for _, ex := range exs {
		m[ex.LabelId] = ex.Extractor
	}
	return m
}

func (s *Server) organiseLocation(ctx context.Context, c *pb.Location) (int32, error) {
	t := time.Now()
	ctx = s.LogTrace(ctx, "organiseLocation", time.Now(), pbt.Milestone_START_FUNCTION)
	s.lastOrgFolder = c.Name
	fr, err := s.bridge.getReleases(c.GetFolderIds())
	if err != nil {
		return -1, err
	}

	switch c.GetSort() {
	case pb.Location_BY_DATE_ADDED:
		sort.Sort(ByDateAdded(fr))
	case pb.Location_BY_LABEL_CATNO:
		sort.Sort(ByLabelCat{fr, convert(s.org.GetExtractors()), s.Log})
	}

	records := Split(fr, float64(c.GetSlots()))
	c.ReleasesLocation = []*pb.ReleasePlacement{}
	for slot, recs := range records {
		for i, rinloc := range recs {
			c.ReleasesLocation = append(c.ReleasesLocation, &pb.ReleasePlacement{Slot: int32(slot + 1), Index: int32(i), InstanceId: rinloc.GetRelease().InstanceId})
		}
	}

	s.lastOrgTime = time.Now().Sub(t)
	s.LogTrace(ctx, "organiseLocation", time.Now(), pbt.Milestone_END_FUNCTION)
	return int32(len(fr)), nil
}
