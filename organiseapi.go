package main

import (
	"golang.org/x/net/context"

	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/recordsorganiser/proto"
)

//AddLocation adds a location
func (s *Server) AddLocation(ctx context.Context, req *pb.AddLocationRequest) (*pb.AddLocationResponse, error) {
	s.prepareForReorg()

	s.org.Locations = append(s.org.Locations, req.GetAdd())

	err := s.organise(s.org)

	if err != nil {
		return &pb.AddLocationResponse{}, err
	}

	return &pb.AddLocationResponse{Now: s.org}, nil
}

// GetOrganisation gets a given organisation
func (s *Server) GetOrganisation(ctx context.Context, req *pb.GetOrganisationRequest) (*pb.GetOrganisationResponse, error) {
	locations := make([]*pb.Location, 0)

	for _, rloc := range req.GetLocations() {
		for _, loc := range s.org.GetLocations() {
			if utils.FuzzyMatch(rloc, loc) {
				locations = append(locations, loc)
			}
		}
	}
	return &pb.GetOrganisationResponse{Locations: locations}, nil
}
