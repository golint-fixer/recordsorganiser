package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/recordsorganiser/proto"
	"github.com/golang/protobuf/proto"
)

// UpdateLocation updates a given location
func (s *Server) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.UpdateLocationResponse, error) {
	for _, loc := range s.org.GetLocations() {
		if loc.GetName() == req.GetLocation() {
			proto.Merge(loc, req.Update)
		}
	}

	s.saveOrg()
	return &pb.UpdateLocationResponse{}, nil
}

//Locate finds a record in the collection
func (s *Server) Locate(ctx context.Context, req *pb.LocateRequest) (*pb.LocateResponse, error) {
	for _, loc := range s.org.GetLocations() {
		for _, r := range loc.GetReleasesLocation() {
			if r.GetInstanceId() == req.GetInstanceId() {
				return &pb.LocateResponse{FoundLocation: loc}, nil
			}
		}
	}

	return &pb.LocateResponse{}, fmt.Errorf("Unable to locate %v in collection", req.GetInstanceId())
}

//AddLocation adds a location
func (s *Server) AddLocation(ctx context.Context, req *pb.AddLocationRequest) (*pb.AddLocationResponse, error) {
	s.prepareForReorg()

	s.org.Locations = append(s.org.Locations, req.GetAdd())

	_, err := s.organise(s.org)

	if err != nil {
		return &pb.AddLocationResponse{}, err
	}

	s.saveOrg()

	return &pb.AddLocationResponse{Now: s.org}, nil
}

// GetOrganisation gets a given organisation
func (s *Server) GetOrganisation(ctx context.Context, req *pb.GetOrganisationRequest) (*pb.GetOrganisationResponse, error) {
	locations := make([]*pb.Location, 0)
	num := int32(0)

	for _, rloc := range req.GetLocations() {
		for _, loc := range s.org.GetLocations() {
			if utils.FuzzyMatch(rloc, loc) {
				if req.ForceReorg {
					n, err := s.organiseLocation(loc)
					num = n
					if err != nil {
						return &pb.GetOrganisationResponse{}, err
					}
				}
				locations = append(locations, loc)
			}
		}
	}

	if req.GetForceReorg() {
		s.saveOrg()
	}
	return &pb.GetOrganisationResponse{Locations: locations, NumberProcessed: num}, nil
}

// GetQuota fills out the quota response
func (s *Server) GetQuota(ctx context.Context, req *pb.QuotaRequest) (*pb.QuotaResponse, error) {
	t := time.Now()
	for _, loc := range s.org.GetLocations() {
		for _, id := range loc.GetFolderIds() {
			if id == req.GetFolderId() {
				s.organiseLocation(loc)

				if loc.GetQuota().GetNumOfSlots() > 0 && len(loc.GetReleasesLocation()) >= int(loc.GetQuota().GetNumOfSlots()) {
					s.LogFunction("GetQuota-true", t)
					if !loc.GetNoAlert() {
						s.gh.alert(loc)
					}
					return &pb.QuotaResponse{SpillFolder: loc.SpillFolder, OverQuota: true}, nil
				}

				s.LogFunction("GetQuota-false", t)
				return &pb.QuotaResponse{OverQuota: false}, nil
			}
		}
	}

	s.LogFunction("GetQuota-notfound", t)
	return &pb.QuotaResponse{}, fmt.Errorf("Unable to locate folder in request (%v)", req.GetFolderId())
}
