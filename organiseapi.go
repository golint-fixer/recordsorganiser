package main

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

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
	for _, loc := range s.org.GetLocations() {
		for _, id := range loc.GetFolderIds() {
			if id == req.GetFolderId() {
				if loc.GetQuota().GetNumOfSlots() > 0 && len(loc.GetReleasesLocation()) >= int(loc.GetQuota().GetNumOfSlots()) {
					return &pb.QuotaResponse{OverQuota: true}, nil
				}

				return &pb.QuotaResponse{OverQuota: false}, nil
			}
		}
	}

	return &pb.QuotaResponse{}, fmt.Errorf("Unable to locate folder in request")
}
