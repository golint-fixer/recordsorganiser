// Code generated by protoc-gen-go. DO NOT EDIT.
// source: organise.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	organise.proto

It has these top-level messages:
	Empty
	CombinedRelease
	ReleasePlacement
	Location
	LocationList
	Organisation
	OrganisationList
	LocationMove
	OrganisationMoves
	ReleaseLocation
	DiffRequest
	CleanList
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import godiscogs "github.com/brotherlogic/godiscogs"
import discogsserver "github.com/brotherlogic/discogssyncer/server"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// The means by which the folder is sorted
type Location_Sorting int32

const (
	Location_BY_LABEL_CATNO  Location_Sorting = 0
	Location_BY_DATE_ADDED   Location_Sorting = 1
	Location_BY_RELEASE_DATE Location_Sorting = 2
)

var Location_Sorting_name = map[int32]string{
	0: "BY_LABEL_CATNO",
	1: "BY_DATE_ADDED",
	2: "BY_RELEASE_DATE",
}
var Location_Sorting_value = map[string]int32{
	"BY_LABEL_CATNO":  0,
	"BY_DATE_ADDED":   1,
	"BY_RELEASE_DATE": 2,
}

func (x Location_Sorting) String() string {
	return proto1.EnumName(Location_Sorting_name, int32(x))
}
func (Location_Sorting) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto1.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CombinedRelease struct {
	Release  *godiscogs.Release             `protobuf:"bytes,1,opt,name=release" json:"release,omitempty"`
	Metadata *discogsserver.ReleaseMetadata `protobuf:"bytes,2,opt,name=metadata" json:"metadata,omitempty"`
}

func (m *CombinedRelease) Reset()                    { *m = CombinedRelease{} }
func (m *CombinedRelease) String() string            { return proto1.CompactTextString(m) }
func (*CombinedRelease) ProtoMessage()               {}
func (*CombinedRelease) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CombinedRelease) GetRelease() *godiscogs.Release {
	if m != nil {
		return m.Release
	}
	return nil
}

func (m *CombinedRelease) GetMetadata() *discogsserver.ReleaseMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type ReleasePlacement struct {
	// The id of the release
	ReleaseId int32 `protobuf:"varint,1,opt,name=release_id,json=releaseId" json:"release_id,omitempty"`
	// The index in the folder
	Index int32 `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	// The slot in the folder
	Slot int32 `protobuf:"varint,3,opt,name=slot" json:"slot,omitempty"`
	// The prior release
	BeforeReleaseId int32 `protobuf:"varint,4,opt,name=before_release_id,json=beforeReleaseId" json:"before_release_id,omitempty"`
	// The following release
	AfterReleaseId int32 `protobuf:"varint,5,opt,name=after_release_id,json=afterReleaseId" json:"after_release_id,omitempty"`
	// The name of the folder
	Folder string `protobuf:"bytes,6,opt,name=folder" json:"folder,omitempty"`
}

func (m *ReleasePlacement) Reset()                    { *m = ReleasePlacement{} }
func (m *ReleasePlacement) String() string            { return proto1.CompactTextString(m) }
func (*ReleasePlacement) ProtoMessage()               {}
func (*ReleasePlacement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ReleasePlacement) GetReleaseId() int32 {
	if m != nil {
		return m.ReleaseId
	}
	return 0
}

func (m *ReleasePlacement) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *ReleasePlacement) GetSlot() int32 {
	if m != nil {
		return m.Slot
	}
	return 0
}

func (m *ReleasePlacement) GetBeforeReleaseId() int32 {
	if m != nil {
		return m.BeforeReleaseId
	}
	return 0
}

func (m *ReleasePlacement) GetAfterReleaseId() int32 {
	if m != nil {
		return m.AfterReleaseId
	}
	return 0
}

func (m *ReleasePlacement) GetFolder() string {
	if m != nil {
		return m.Folder
	}
	return ""
}

type Location struct {
	// The name of the location
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// The number of units within the location
	Units int32 `protobuf:"varint,2,opt,name=units" json:"units,omitempty"`
	// The folder ids that are stored in this location
	FolderIds []int32 `protobuf:"varint,3,rep,packed,name=folder_ids,json=folderIds" json:"folder_ids,omitempty"`
	// The placement of releases in the folder
	ReleasesLocation []*ReleasePlacement `protobuf:"bytes,4,rep,name=releases_location,json=releasesLocation" json:"releases_location,omitempty"`
	Sort             Location_Sorting    `protobuf:"varint,5,opt,name=sort,enum=proto.Location_Sorting" json:"sort,omitempty"`
	// The timestamp of this given location / arrangement
	Timestamp int64 `protobuf:"varint,6,opt,name=timestamp" json:"timestamp,omitempty"`
	// The allowed quota for this location, if any
	Quota int32 `protobuf:"varint,7,opt,name=quota" json:"quota,omitempty"`
	// The type of format we expect in this location
	ExpectedFormat string `protobuf:"bytes,8,opt,name=expected_format,json=expectedFormat" json:"expected_format,omitempty"`
	// The type of label we don't expect in this location
	UnexpectedLabel string `protobuf:"bytes,9,opt,name=unexpected_label,json=unexpectedLabel" json:"unexpected_label,omitempty"`
}

func (m *Location) Reset()                    { *m = Location{} }
func (m *Location) String() string            { return proto1.CompactTextString(m) }
func (*Location) ProtoMessage()               {}
func (*Location) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Location) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Location) GetUnits() int32 {
	if m != nil {
		return m.Units
	}
	return 0
}

func (m *Location) GetFolderIds() []int32 {
	if m != nil {
		return m.FolderIds
	}
	return nil
}

func (m *Location) GetReleasesLocation() []*ReleasePlacement {
	if m != nil {
		return m.ReleasesLocation
	}
	return nil
}

func (m *Location) GetSort() Location_Sorting {
	if m != nil {
		return m.Sort
	}
	return Location_BY_LABEL_CATNO
}

func (m *Location) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Location) GetQuota() int32 {
	if m != nil {
		return m.Quota
	}
	return 0
}

func (m *Location) GetExpectedFormat() string {
	if m != nil {
		return m.ExpectedFormat
	}
	return ""
}

func (m *Location) GetUnexpectedLabel() string {
	if m != nil {
		return m.UnexpectedLabel
	}
	return ""
}

type LocationList struct {
	Locations []*Location `protobuf:"bytes,1,rep,name=locations" json:"locations,omitempty"`
}

func (m *LocationList) Reset()                    { *m = LocationList{} }
func (m *LocationList) String() string            { return proto1.CompactTextString(m) }
func (*LocationList) ProtoMessage()               {}
func (*LocationList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *LocationList) GetLocations() []*Location {
	if m != nil {
		return m.Locations
	}
	return nil
}

type Organisation struct {
	// Timestamp this organisation was made
	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	// The locations in this sorting
	Locations []*Location `protobuf:"bytes,2,rep,name=locations" json:"locations,omitempty"`
}

func (m *Organisation) Reset()                    { *m = Organisation{} }
func (m *Organisation) String() string            { return proto1.CompactTextString(m) }
func (*Organisation) ProtoMessage()               {}
func (*Organisation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Organisation) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Organisation) GetLocations() []*Location {
	if m != nil {
		return m.Locations
	}
	return nil
}

type OrganisationList struct {
	Organisations []*Organisation `protobuf:"bytes,1,rep,name=organisations" json:"organisations,omitempty"`
}

func (m *OrganisationList) Reset()                    { *m = OrganisationList{} }
func (m *OrganisationList) String() string            { return proto1.CompactTextString(m) }
func (*OrganisationList) ProtoMessage()               {}
func (*OrganisationList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *OrganisationList) GetOrganisations() []*Organisation {
	if m != nil {
		return m.Organisations
	}
	return nil
}

type LocationMove struct {
	Old      *ReleasePlacement `protobuf:"bytes,1,opt,name=old" json:"old,omitempty"`
	New      *ReleasePlacement `protobuf:"bytes,2,opt,name=new" json:"new,omitempty"`
	SlotMove bool              `protobuf:"varint,3,opt,name=slot_move,json=slotMove" json:"slot_move,omitempty"`
}

func (m *LocationMove) Reset()                    { *m = LocationMove{} }
func (m *LocationMove) String() string            { return proto1.CompactTextString(m) }
func (*LocationMove) ProtoMessage()               {}
func (*LocationMove) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *LocationMove) GetOld() *ReleasePlacement {
	if m != nil {
		return m.Old
	}
	return nil
}

func (m *LocationMove) GetNew() *ReleasePlacement {
	if m != nil {
		return m.New
	}
	return nil
}

func (m *LocationMove) GetSlotMove() bool {
	if m != nil {
		return m.SlotMove
	}
	return false
}

type OrganisationMoves struct {
	StartTimestamp int64           `protobuf:"varint,1,opt,name=start_timestamp,json=startTimestamp" json:"start_timestamp,omitempty"`
	EndTimestamp   int64           `protobuf:"varint,2,opt,name=end_timestamp,json=endTimestamp" json:"end_timestamp,omitempty"`
	Moves          []*LocationMove `protobuf:"bytes,3,rep,name=moves" json:"moves,omitempty"`
}

func (m *OrganisationMoves) Reset()                    { *m = OrganisationMoves{} }
func (m *OrganisationMoves) String() string            { return proto1.CompactTextString(m) }
func (*OrganisationMoves) ProtoMessage()               {}
func (*OrganisationMoves) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *OrganisationMoves) GetStartTimestamp() int64 {
	if m != nil {
		return m.StartTimestamp
	}
	return 0
}

func (m *OrganisationMoves) GetEndTimestamp() int64 {
	if m != nil {
		return m.EndTimestamp
	}
	return 0
}

func (m *OrganisationMoves) GetMoves() []*LocationMove {
	if m != nil {
		return m.Moves
	}
	return nil
}

type ReleaseLocation struct {
	Location *Location          `protobuf:"bytes,1,opt,name=location" json:"location,omitempty"`
	Slot     int32              `protobuf:"varint,2,opt,name=slot" json:"slot,omitempty"`
	Before   *godiscogs.Release `protobuf:"bytes,3,opt,name=before" json:"before,omitempty"`
	After    *godiscogs.Release `protobuf:"bytes,4,opt,name=after" json:"after,omitempty"`
}

func (m *ReleaseLocation) Reset()                    { *m = ReleaseLocation{} }
func (m *ReleaseLocation) String() string            { return proto1.CompactTextString(m) }
func (*ReleaseLocation) ProtoMessage()               {}
func (*ReleaseLocation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ReleaseLocation) GetLocation() *Location {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *ReleaseLocation) GetSlot() int32 {
	if m != nil {
		return m.Slot
	}
	return 0
}

func (m *ReleaseLocation) GetBefore() *godiscogs.Release {
	if m != nil {
		return m.Before
	}
	return nil
}

func (m *ReleaseLocation) GetAfter() *godiscogs.Release {
	if m != nil {
		return m.After
	}
	return nil
}

type DiffRequest struct {
	LocationName string `protobuf:"bytes,3,opt,name=location_name,json=locationName" json:"location_name,omitempty"`
	Slot         int32  `protobuf:"varint,4,opt,name=slot" json:"slot,omitempty"`
}

func (m *DiffRequest) Reset()                    { *m = DiffRequest{} }
func (m *DiffRequest) String() string            { return proto1.CompactTextString(m) }
func (*DiffRequest) ProtoMessage()               {}
func (*DiffRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *DiffRequest) GetLocationName() string {
	if m != nil {
		return m.LocationName
	}
	return ""
}

func (m *DiffRequest) GetSlot() int32 {
	if m != nil {
		return m.Slot
	}
	return 0
}

type CleanList struct {
	Entries []*godiscogs.Release `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
}

func (m *CleanList) Reset()                    { *m = CleanList{} }
func (m *CleanList) String() string            { return proto1.CompactTextString(m) }
func (*CleanList) ProtoMessage()               {}
func (*CleanList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *CleanList) GetEntries() []*godiscogs.Release {
	if m != nil {
		return m.Entries
	}
	return nil
}

func init() {
	proto1.RegisterType((*Empty)(nil), "proto.Empty")
	proto1.RegisterType((*CombinedRelease)(nil), "proto.CombinedRelease")
	proto1.RegisterType((*ReleasePlacement)(nil), "proto.ReleasePlacement")
	proto1.RegisterType((*Location)(nil), "proto.Location")
	proto1.RegisterType((*LocationList)(nil), "proto.LocationList")
	proto1.RegisterType((*Organisation)(nil), "proto.Organisation")
	proto1.RegisterType((*OrganisationList)(nil), "proto.OrganisationList")
	proto1.RegisterType((*LocationMove)(nil), "proto.LocationMove")
	proto1.RegisterType((*OrganisationMoves)(nil), "proto.OrganisationMoves")
	proto1.RegisterType((*ReleaseLocation)(nil), "proto.ReleaseLocation")
	proto1.RegisterType((*DiffRequest)(nil), "proto.DiffRequest")
	proto1.RegisterType((*CleanList)(nil), "proto.CleanList")
	proto1.RegisterEnum("proto.Location_Sorting", Location_Sorting_name, Location_Sorting_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for OrganiserService service

type OrganiserServiceClient interface {
	Organise(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*OrganisationMoves, error)
	Locate(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*ReleaseLocation, error)
	AddLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*Location, error)
	GetLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*Location, error)
	GetOrganisation(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Organisation, error)
	UpdateLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*Location, error)
	Diff(ctx context.Context, in *DiffRequest, opts ...grpc.CallOption) (*OrganisationMoves, error)
	GetQuotaViolations(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LocationList, error)
	CleanLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*CleanList, error)
}

type organiserServiceClient struct {
	cc *grpc.ClientConn
}

func NewOrganiserServiceClient(cc *grpc.ClientConn) OrganiserServiceClient {
	return &organiserServiceClient{cc}
}

func (c *organiserServiceClient) Organise(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*OrganisationMoves, error) {
	out := new(OrganisationMoves)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/Organise", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organiserServiceClient) Locate(ctx context.Context, in *godiscogs.Release, opts ...grpc.CallOption) (*ReleaseLocation, error) {
	out := new(ReleaseLocation)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/Locate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organiserServiceClient) AddLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*Location, error) {
	out := new(Location)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/AddLocation", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organiserServiceClient) GetLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*Location, error) {
	out := new(Location)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/GetLocation", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organiserServiceClient) GetOrganisation(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Organisation, error) {
	out := new(Organisation)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/GetOrganisation", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organiserServiceClient) UpdateLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*Location, error) {
	out := new(Location)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/UpdateLocation", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organiserServiceClient) Diff(ctx context.Context, in *DiffRequest, opts ...grpc.CallOption) (*OrganisationMoves, error) {
	out := new(OrganisationMoves)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/Diff", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organiserServiceClient) GetQuotaViolations(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LocationList, error) {
	out := new(LocationList)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/GetQuotaViolations", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *organiserServiceClient) CleanLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*CleanList, error) {
	out := new(CleanList)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/CleanLocation", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OrganiserService service

type OrganiserServiceServer interface {
	Organise(context.Context, *Empty) (*OrganisationMoves, error)
	Locate(context.Context, *godiscogs.Release) (*ReleaseLocation, error)
	AddLocation(context.Context, *Location) (*Location, error)
	GetLocation(context.Context, *Location) (*Location, error)
	GetOrganisation(context.Context, *Empty) (*Organisation, error)
	UpdateLocation(context.Context, *Location) (*Location, error)
	Diff(context.Context, *DiffRequest) (*OrganisationMoves, error)
	GetQuotaViolations(context.Context, *Empty) (*LocationList, error)
	CleanLocation(context.Context, *Location) (*CleanList, error)
}

func RegisterOrganiserServiceServer(s *grpc.Server, srv OrganiserServiceServer) {
	s.RegisterService(&_OrganiserService_serviceDesc, srv)
}

func _OrganiserService_Organise_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).Organise(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/Organise",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).Organise(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganiserService_Locate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(godiscogs.Release)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).Locate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/Locate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).Locate(ctx, req.(*godiscogs.Release))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganiserService_AddLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Location)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).AddLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/AddLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).AddLocation(ctx, req.(*Location))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganiserService_GetLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Location)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).GetLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/GetLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).GetLocation(ctx, req.(*Location))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganiserService_GetOrganisation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).GetOrganisation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/GetOrganisation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).GetOrganisation(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganiserService_UpdateLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Location)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).UpdateLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/UpdateLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).UpdateLocation(ctx, req.(*Location))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganiserService_Diff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiffRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).Diff(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/Diff",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).Diff(ctx, req.(*DiffRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganiserService_GetQuotaViolations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).GetQuotaViolations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/GetQuotaViolations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).GetQuotaViolations(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrganiserService_CleanLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Location)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).CleanLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/CleanLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).CleanLocation(ctx, req.(*Location))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrganiserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.OrganiserService",
	HandlerType: (*OrganiserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Organise",
			Handler:    _OrganiserService_Organise_Handler,
		},
		{
			MethodName: "Locate",
			Handler:    _OrganiserService_Locate_Handler,
		},
		{
			MethodName: "AddLocation",
			Handler:    _OrganiserService_AddLocation_Handler,
		},
		{
			MethodName: "GetLocation",
			Handler:    _OrganiserService_GetLocation_Handler,
		},
		{
			MethodName: "GetOrganisation",
			Handler:    _OrganiserService_GetOrganisation_Handler,
		},
		{
			MethodName: "UpdateLocation",
			Handler:    _OrganiserService_UpdateLocation_Handler,
		},
		{
			MethodName: "Diff",
			Handler:    _OrganiserService_Diff_Handler,
		},
		{
			MethodName: "GetQuotaViolations",
			Handler:    _OrganiserService_GetQuotaViolations_Handler,
		},
		{
			MethodName: "CleanLocation",
			Handler:    _OrganiserService_CleanLocation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "organise.proto",
}

func init() { proto1.RegisterFile("organise.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 926 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0x4b, 0x6f, 0xdb, 0x46,
	0x10, 0x16, 0x45, 0x3d, 0xc7, 0xb2, 0x44, 0x6f, 0x8a, 0x94, 0x50, 0x1f, 0x30, 0xd8, 0x43, 0xe5,
	0x36, 0x95, 0x11, 0x37, 0x08, 0xe0, 0x02, 0x3d, 0xc8, 0x96, 0x62, 0x04, 0x90, 0x93, 0x96, 0x76,
	0x0b, 0x18, 0x3d, 0x10, 0x94, 0x38, 0x52, 0x08, 0x90, 0x5c, 0x85, 0xbb, 0x72, 0x12, 0xf4, 0xd4,
	0x7b, 0xaf, 0xfd, 0x0f, 0xfd, 0x2b, 0xfd, 0x51, 0x05, 0x8a, 0x7d, 0xf0, 0x21, 0x55, 0x6a, 0xe1,
	0x93, 0xb8, 0xdf, 0x7e, 0x3b, 0xfb, 0xcd, 0xb7, 0x33, 0x23, 0xe8, 0xd2, 0x74, 0xe9, 0x27, 0x21,
	0xc3, 0xe1, 0x2a, 0xa5, 0x9c, 0x92, 0xba, 0xfc, 0xe9, 0x3f, 0x5d, 0x86, 0xfc, 0xcd, 0x7a, 0x36,
	0x9c, 0xd3, 0xf8, 0x74, 0x96, 0x52, 0xfe, 0x06, 0xd3, 0x88, 0x2e, 0xc3, 0xf9, 0xe9, 0x92, 0x06,
	0x21, 0x9b, 0xd3, 0x25, 0x2b, 0xbe, 0xd4, 0xc9, 0xfe, 0xf9, 0xbe, 0x23, 0x9a, 0xc6, 0x3e, 0x24,
	0x73, 0x4c, 0x4f, 0x19, 0xa6, 0xf7, 0xf9, 0x8f, 0x3a, 0xea, 0x34, 0xa1, 0x3e, 0x89, 0x57, 0xfc,
	0x83, 0xf3, 0x2b, 0xf4, 0x2e, 0x69, 0x3c, 0x0b, 0x13, 0x0c, 0x5c, 0x8c, 0xd0, 0x67, 0x48, 0x9e,
	0x40, 0x33, 0x55, 0x9f, 0xb6, 0x71, 0x6c, 0x0c, 0x0e, 0xce, 0xc8, 0xb0, 0xb8, 0x59, 0x93, 0xdc,
	0x8c, 0x42, 0xbe, 0x83, 0x56, 0x8c, 0xdc, 0x0f, 0x7c, 0xee, 0xdb, 0x55, 0x49, 0xff, 0x7c, 0x98,
	0xdd, 0xaf, 0x6e, 0xd4, 0x47, 0xae, 0x35, 0xcb, 0xcd, 0xf9, 0xce, 0x5f, 0x06, 0x58, 0x7a, 0xf7,
	0x87, 0xc8, 0x9f, 0x63, 0x8c, 0x09, 0x27, 0x9f, 0x01, 0xe8, 0xd8, 0x5e, 0x18, 0x48, 0x05, 0x75,
	0xb7, 0xad, 0x91, 0x97, 0x01, 0xf9, 0x08, 0xea, 0x61, 0x12, 0xe0, 0x7b, 0x79, 0x59, 0xdd, 0x55,
	0x0b, 0x42, 0xa0, 0xc6, 0x22, 0xca, 0x6d, 0x53, 0x82, 0xf2, 0x9b, 0x7c, 0x05, 0x47, 0x33, 0x5c,
	0xd0, 0x14, 0xbd, 0x52, 0xbc, 0x9a, 0x24, 0xf4, 0xd4, 0x86, 0x9b, 0x47, 0x1d, 0x80, 0xe5, 0x2f,
	0x38, 0xa6, 0x65, 0x6a, 0x5d, 0x52, 0xbb, 0x12, 0x2f, 0x98, 0x8f, 0xa1, 0xb1, 0xa0, 0x51, 0x80,
	0xa9, 0xdd, 0x38, 0x36, 0x06, 0x6d, 0x57, 0xaf, 0x9c, 0x3f, 0x4c, 0x68, 0x4d, 0xe9, 0xdc, 0xe7,
	0x21, 0x4d, 0x84, 0x9c, 0xc4, 0x8f, 0x95, 0x7f, 0x6d, 0x57, 0x7e, 0x0b, 0xe1, 0xeb, 0x24, 0xe4,
	0x2c, 0x13, 0x2e, 0x17, 0x22, 0x5b, 0x15, 0xc0, 0x0b, 0x03, 0x66, 0x9b, 0xc7, 0xa6, 0xc8, 0x56,
	0x21, 0x2f, 0x03, 0x46, 0xc6, 0x70, 0xa4, 0x15, 0x31, 0x2f, 0xd2, 0xd1, 0xed, 0xda, 0xb1, 0x39,
	0x38, 0x38, 0xfb, 0x58, 0x3d, 0xe5, 0x70, 0xdb, 0x40, 0xd7, 0xca, 0x4e, 0xe4, 0x72, 0xbe, 0x86,
	0x1a, 0xa3, 0x29, 0x97, 0x19, 0x75, 0xf3, 0x83, 0xd9, 0xf6, 0xf0, 0x86, 0xa6, 0x3c, 0x4c, 0x96,
	0xae, 0x24, 0x91, 0x4f, 0xa1, 0xcd, 0xc3, 0x18, 0x19, 0xf7, 0xe3, 0x95, 0xcc, 0xd1, 0x74, 0x0b,
	0x40, 0x64, 0xf1, 0x76, 0x4d, 0xb9, 0x6f, 0x37, 0x55, 0x16, 0x72, 0x41, 0xbe, 0x84, 0x1e, 0xbe,
	0x5f, 0xe1, 0x9c, 0x63, 0xe0, 0x2d, 0x68, 0x1a, 0xfb, 0xdc, 0x6e, 0xc9, 0xd4, 0xbb, 0x19, 0xfc,
	0x42, 0xa2, 0xe4, 0x04, 0xac, 0x75, 0x92, 0x53, 0x23, 0x7f, 0x86, 0x91, 0xdd, 0x96, 0xcc, 0x5e,
	0x81, 0x4f, 0x05, 0xec, 0x4c, 0xa0, 0xa9, 0x85, 0x11, 0x02, 0xdd, 0x8b, 0x3b, 0x6f, 0x3a, 0xba,
	0x98, 0x4c, 0xbd, 0xcb, 0xd1, 0xed, 0xab, 0xd7, 0x56, 0x85, 0x1c, 0xc1, 0xe1, 0xc5, 0x9d, 0x37,
	0x1e, 0xdd, 0x4e, 0xbc, 0xd1, 0x78, 0x3c, 0x19, 0x5b, 0x06, 0x79, 0x04, 0xbd, 0x8b, 0x3b, 0xcf,
	0x9d, 0x4c, 0x27, 0xa3, 0x9b, 0x89, 0xdc, 0xb2, 0xaa, 0xce, 0xf7, 0xd0, 0xc9, 0x12, 0x9d, 0x86,
	0x8c, 0x93, 0x6f, 0xa0, 0x9d, 0x19, 0xc9, 0x6c, 0x43, 0x3a, 0xd9, 0xdb, 0x32, 0xc4, 0x2d, 0x18,
	0xce, 0x2f, 0xd0, 0x79, 0xad, 0xfa, 0x55, 0x59, 0xb9, 0xe1, 0x8e, 0xb1, 0xed, 0xce, 0x46, 0xf0,
	0xea, 0xff, 0x06, 0xbf, 0x06, 0xab, 0x1c, 0x5c, 0xea, 0x3b, 0x87, 0x43, 0x5a, 0xc2, 0x32, 0x8d,
	0x8f, 0x74, 0x98, 0x32, 0xdf, 0xdd, 0x64, 0x3a, 0xbf, 0x19, 0x45, 0xae, 0xd7, 0xf4, 0x1e, 0xc9,
	0x09, 0x98, 0x34, 0x0a, 0x74, 0x17, 0xef, 0xad, 0x17, 0xc1, 0x11, 0xd4, 0x04, 0xdf, 0xe9, 0x0e,
	0xde, 0x4f, 0x4d, 0xf0, 0x1d, 0xf9, 0x04, 0xda, 0xa2, 0xbf, 0xbc, 0x98, 0xde, 0xa3, 0x6c, 0xb8,
	0x96, 0xdb, 0x12, 0x80, 0xb8, 0xd2, 0xf9, 0xdd, 0x80, 0xa3, 0xb2, 0x46, 0x01, 0x32, 0x51, 0x1f,
	0x8c, 0xfb, 0x29, 0xf7, 0xb6, 0xbd, 0xeb, 0x4a, 0xf8, 0x36, 0x37, 0xf0, 0x0b, 0x38, 0xc4, 0x24,
	0x28, 0xd1, 0xaa, 0x92, 0xd6, 0xc1, 0x24, 0x28, 0x48, 0x27, 0x50, 0x17, 0x77, 0xab, 0x76, 0x29,
	0xac, 0x29, 0xa7, 0xee, 0x2a, 0x86, 0xf3, 0xa7, 0x01, 0x3d, 0x9d, 0x45, 0xa9, 0x1b, 0x5a, 0x79,
	0x2b, 0x29, 0x6b, 0xfe, 0xf5, 0x46, 0x39, 0x21, 0x1f, 0x2c, 0xd5, 0x8d, 0xc1, 0xd2, 0x50, 0xf3,
	0x43, 0x66, 0xbf, 0x7b, 0x3e, 0x6a, 0x06, 0x19, 0x40, 0x5d, 0x0e, 0x10, 0x39, 0x78, 0x76, 0x53,
	0x15, 0xc1, 0x79, 0x01, 0x07, 0xe3, 0x70, 0xb1, 0x70, 0xf1, 0xed, 0x1a, 0x19, 0x17, 0x4e, 0x64,
	0x22, 0x3c, 0x39, 0x4b, 0x4c, 0xd9, 0x26, 0x9d, 0x0c, 0x7c, 0x25, 0x66, 0x4a, 0xa6, 0xae, 0x56,
	0xa8, 0x73, 0xce, 0xa1, 0x7d, 0x19, 0xa1, 0xaf, 0xaa, 0xe9, 0x09, 0x34, 0x31, 0xe1, 0x69, 0x88,
	0x59, 0x1d, 0xed, 0x9c, 0xe5, 0x9a, 0x72, 0xf6, 0xb7, 0x99, 0x17, 0x24, 0xa6, 0x37, 0x98, 0xde,
	0x87, 0x73, 0x24, 0xcf, 0xa0, 0x95, 0x61, 0xa4, 0xa3, 0x8d, 0x92, 0xff, 0x1d, 0x7d, 0x7b, 0x47,
	0x4d, 0xca, 0xf7, 0x76, 0x2a, 0xe4, 0x39, 0x34, 0xa4, 0x9b, 0x48, 0x76, 0xdc, 0xd8, 0x7f, 0xbc,
	0x59, 0x60, 0x99, 0xef, 0x4e, 0x85, 0x3c, 0x85, 0x83, 0x51, 0x10, 0xe4, 0x6f, 0xb5, 0xfd, 0x32,
	0xfd, 0x6d, 0x40, 0x1d, 0xb9, 0x42, 0xfe, 0xa0, 0x23, 0xcf, 0xa1, 0x77, 0x85, 0x7c, 0xa3, 0xb1,
	0x37, 0x53, 0xdb, 0xd5, 0x6e, 0x4e, 0x85, 0x3c, 0x83, 0xee, 0x4f, 0xab, 0xc0, 0xe7, 0xf8, 0xc0,
	0xdb, 0x6a, 0xe2, 0x65, 0x09, 0xd1, 0x5b, 0xa5, 0x67, 0xfe, 0x4f, 0x0f, 0xcf, 0x81, 0x5c, 0x21,
	0xff, 0x51, 0x4c, 0xd8, 0x9f, 0x43, 0x1a, 0xa9, 0x2e, 0xdf, 0x23, 0xb4, 0x3c, 0xe3, 0xa4, 0xd0,
	0x43, 0x55, 0x04, 0x7b, 0x75, 0x5a, 0x1a, 0xc8, 0x6b, 0xc5, 0xa9, 0xcc, 0x1a, 0x12, 0xfa, 0xf6,
	0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x42, 0x20, 0x8a, 0x55, 0xa3, 0x08, 0x00, 0x00,
}
