package main

import (
	"golang.org/x/net/context"

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
