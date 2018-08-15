package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brotherlogic/goserver/utils"
	"github.com/golang/protobuf/proto"

	pbgd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/recordsorganiser/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

// UpdateLocation updates a given location
func (s *Server) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.UpdateLocationResponse, error) {
	for _, loc := range s.org.GetLocations() {
		if loc.GetName() == req.GetLocation() {
			proto.Merge(loc, req.Update)
		}
	}

	s.saveOrg(ctx)
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
	s.saveOrg(ctx)

	_, err := s.organise(s.org)

	if err != nil {
		return &pb.AddLocationResponse{}, err
	}

	return &pb.AddLocationResponse{Now: s.org}, nil
}

// GetOrganisation gets a given organisation
func (s *Server) GetOrganisation(ctx context.Context, req *pb.GetOrganisationRequest) (*pb.GetOrganisationResponse, error) {
	ctx = s.LogTrace(ctx, "GetOrganisation", time.Now(), pbt.Milestone_START_FUNCTION)
	locations := make([]*pb.Location, 0)
	num := int32(0)

	for _, rloc := range req.GetLocations() {
		for _, loc := range s.org.GetLocations() {
			if utils.FuzzyMatch(rloc, loc) {
				if req.ForceReorg {
					n, err := s.organiseLocation(ctx, loc)
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
		s.saveOrg(ctx)
	}

	s.LogTrace(ctx, "GetOrganisation", time.Now(), pbt.Milestone_END_FUNCTION)
	return &pb.GetOrganisationResponse{Locations: locations, NumberProcessed: num}, nil
}

// GetQuota fills out the quota response
func (s *Server) GetQuota(ctx context.Context, req *pb.QuotaRequest) (*pb.QuotaResponse, error) {
	s.LogTrace(ctx, "GetQuota", time.Now(), pbt.Milestone_START_FUNCTION)

	instanceIds := []int32{}

	folderIds := []int32{}
	for _, loc := range s.org.GetLocations() {
		if loc.Name == req.Name {
			folderIds = append(folderIds, loc.FolderIds...)
		}
	}

	if len(folderIds) == 0 {
		folderIds = append(folderIds, req.FolderId)
	}

	//Compute the count of valid records in the listening pile
	count := 0
	for _, loc := range s.org.GetLocations() {
		log.Printf("Trying %v", loc.Name)
		if loc.Name == "Listening Pile" {
			for _, place := range loc.ReleasesLocation {
				meta, err := s.bridge.getMetadata(&pbgd.Release{InstanceId: place.InstanceId})
				if err == nil {
					for _, fid := range folderIds {
						if meta.GoalFolder == fid {
							if meta.Category != pbrc.ReleaseMetadata_UNLISTENED &&
								meta.Category != pbrc.ReleaseMetadata_STAGED &&
								meta.Category != pbrc.ReleaseMetadata_STAGED_TO_SELL &&
								meta.Category != pbrc.ReleaseMetadata_SOLD &&
								meta.Category != pbrc.ReleaseMetadata_PREPARE_TO_SELL &&
								meta.Category != pbrc.ReleaseMetadata_PRE_FRESHMAN {
								instanceIds = append(instanceIds, place.InstanceId)
								count++
							}
						}
					}
				}
			}
		}
	}

	for _, loc := range s.org.GetLocations() {
		for _, id := range loc.GetFolderIds() {
			if id == req.GetFolderId() || (req.Name == loc.Name) {
				if loc.GetQuota().GetNumOfSlots() > 0 && len(loc.GetReleasesLocation())+count >= int(loc.GetQuota().GetNumOfSlots()) {
					s.LogTrace(ctx, "GetQuota", time.Now(), pbt.Milestone_END_FUNCTION)
					for _, in := range loc.ReleasesLocation {
						meta, err := s.bridge.getMetadata(&pbgd.Release{InstanceId: in.InstanceId})
						if err == nil {
							if meta.Category != pbrc.ReleaseMetadata_STAGED_TO_SELL &&
								meta.Category != pbrc.ReleaseMetadata_SOLD {
								instanceIds = append(instanceIds, in.InstanceId)
							}
						}
					}

					if len(instanceIds) > int(loc.GetQuota().GetNumOfSlots()) {
						if !loc.GetNoAlert() {
							s.gh.alert(loc)
						}
					}

					return &pb.QuotaResponse{SpillFolder: loc.SpillFolder, OverQuota: len(instanceIds) > int(loc.GetQuota().GetNumOfSlots()), LocationName: loc.GetName(), InstanceId: instanceIds}, nil
				}

				s.LogTrace(ctx, "GetQuota", time.Now(), pbt.Milestone_END_FUNCTION)
				return &pb.QuotaResponse{OverQuota: false, LocationName: loc.GetName()}, nil
			}
		}
	}

	s.LogTrace(ctx, "GetQuota", time.Now(), pbt.Milestone_END_FUNCTION)
	return &pb.QuotaResponse{}, status.Error(codes.InvalidArgument, fmt.Sprintf("Unable to locate folder in request (%v)", req.GetFolderId()))
}

// AddExtractor adds an extractor
func (s *Server) AddExtractor(ctx context.Context, req *pb.AddExtractorRequest) (*pb.AddExtractorResponse, error) {
	s.org.Extractors = append(s.org.Extractors, req.Extractor)
	s.saveOrg(ctx)
	return &pb.AddExtractorResponse{}, nil
}
