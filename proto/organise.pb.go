// Code generated by protoc-gen-go.
// source: organise.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	organise.proto

It has these top-level messages:
	Empty
	ReleasePlacement
	Location
	Organisation
	OrganisationList
	LocationMove
	OrganisationMoves
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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
	Location_BY_LABEL_CATNO Location_Sorting = 0
)

var Location_Sorting_name = map[int32]string{
	0: "BY_LABEL_CATNO",
}
var Location_Sorting_value = map[string]int32{
	"BY_LABEL_CATNO": 0,
}

func (x Location_Sorting) String() string {
	return proto1.EnumName(Location_Sorting_name, int32(x))
}
func (Location_Sorting) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto1.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ReleasePlacement struct {
	// The id of the release
	ReleaseId int32 `protobuf:"varint,1,opt,name=release_id,json=releaseId" json:"release_id,omitempty"`
	// The index in the folder
	Index int32 `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	// The slot in the folder
	Slot int32 `protobuf:"varint,3,opt,name=slot" json:"slot,omitempty"`
}

func (m *ReleasePlacement) Reset()                    { *m = ReleasePlacement{} }
func (m *ReleasePlacement) String() string            { return proto1.CompactTextString(m) }
func (*ReleasePlacement) ProtoMessage()               {}
func (*ReleasePlacement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Location struct {
	// The name of the location
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// The number of units within the location
	Units int32 `protobuf:"varint,2,opt,name=units" json:"units,omitempty"`
	// The folder ids that are stored in this location
	FolderIds []int32 `protobuf:"varint,3,rep,name=folder_ids,json=folderIds" json:"folder_ids,omitempty"`
	// The placement of releases in the folder
	ReleasesLocation []*ReleasePlacement `protobuf:"bytes,4,rep,name=releases_location,json=releasesLocation" json:"releases_location,omitempty"`
	Sort             Location_Sorting    `protobuf:"varint,5,opt,name=sort,enum=proto.Location_Sorting" json:"sort,omitempty"`
}

func (m *Location) Reset()                    { *m = Location{} }
func (m *Location) String() string            { return proto1.CompactTextString(m) }
func (*Location) ProtoMessage()               {}
func (*Location) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Location) GetReleasesLocation() []*ReleasePlacement {
	if m != nil {
		return m.ReleasesLocation
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
func (*Organisation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

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
func (*OrganisationList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *OrganisationList) GetOrganisations() []*Organisation {
	if m != nil {
		return m.Organisations
	}
	return nil
}

type LocationMove struct {
	Old *ReleasePlacement `protobuf:"bytes,1,opt,name=old" json:"old,omitempty"`
	New *ReleasePlacement `protobuf:"bytes,2,opt,name=new" json:"new,omitempty"`
}

func (m *LocationMove) Reset()                    { *m = LocationMove{} }
func (m *LocationMove) String() string            { return proto1.CompactTextString(m) }
func (*LocationMove) ProtoMessage()               {}
func (*LocationMove) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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

type OrganisationMoves struct {
	StartTimestamp int64           `protobuf:"varint,1,opt,name=start_timestamp,json=startTimestamp" json:"start_timestamp,omitempty"`
	EndTimestamp   int64           `protobuf:"varint,2,opt,name=end_timestamp,json=endTimestamp" json:"end_timestamp,omitempty"`
	Moves          []*LocationMove `protobuf:"bytes,3,rep,name=moves" json:"moves,omitempty"`
}

func (m *OrganisationMoves) Reset()                    { *m = OrganisationMoves{} }
func (m *OrganisationMoves) String() string            { return proto1.CompactTextString(m) }
func (*OrganisationMoves) ProtoMessage()               {}
func (*OrganisationMoves) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *OrganisationMoves) GetMoves() []*LocationMove {
	if m != nil {
		return m.Moves
	}
	return nil
}

func init() {
	proto1.RegisterType((*Empty)(nil), "proto.Empty")
	proto1.RegisterType((*ReleasePlacement)(nil), "proto.ReleasePlacement")
	proto1.RegisterType((*Location)(nil), "proto.Location")
	proto1.RegisterType((*Organisation)(nil), "proto.Organisation")
	proto1.RegisterType((*OrganisationList)(nil), "proto.OrganisationList")
	proto1.RegisterType((*LocationMove)(nil), "proto.LocationMove")
	proto1.RegisterType((*OrganisationMoves)(nil), "proto.OrganisationMoves")
	proto1.RegisterEnum("proto.Location_Sorting", Location_Sorting_name, Location_Sorting_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for OrganiserService service

type OrganiserServiceClient interface {
	// rpc Organise (Empty) returns (OrganisationMoves) {};
	// rpc Locate (godiscogs.Release) returns (Location) {};
	AddLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*Location, error)
	GetLocation(ctx context.Context, in *Location, opts ...grpc.CallOption) (*Location, error)
	GetOrganisations(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*OrganisationList, error)
}

type organiserServiceClient struct {
	cc *grpc.ClientConn
}

func NewOrganiserServiceClient(cc *grpc.ClientConn) OrganiserServiceClient {
	return &organiserServiceClient{cc}
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

func (c *organiserServiceClient) GetOrganisations(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*OrganisationList, error) {
	out := new(OrganisationList)
	err := grpc.Invoke(ctx, "/proto.OrganiserService/GetOrganisations", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OrganiserService service

type OrganiserServiceServer interface {
	// rpc Organise (Empty) returns (OrganisationMoves) {};
	// rpc Locate (godiscogs.Release) returns (Location) {};
	AddLocation(context.Context, *Location) (*Location, error)
	GetLocation(context.Context, *Location) (*Location, error)
	GetOrganisations(context.Context, *Empty) (*OrganisationList, error)
}

func RegisterOrganiserServiceServer(s *grpc.Server, srv OrganiserServiceServer) {
	s.RegisterService(&_OrganiserService_serviceDesc, srv)
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

func _OrganiserService_GetOrganisations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrganiserServiceServer).GetOrganisations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.OrganiserService/GetOrganisations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrganiserServiceServer).GetOrganisations(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _OrganiserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.OrganiserService",
	HandlerType: (*OrganiserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddLocation",
			Handler:    _OrganiserService_AddLocation_Handler,
		},
		{
			MethodName: "GetLocation",
			Handler:    _OrganiserService_GetLocation_Handler,
		},
		{
			MethodName: "GetOrganisations",
			Handler:    _OrganiserService_GetOrganisations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto1.RegisterFile("organise.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 471 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x53, 0x5d, 0x6b, 0xd4, 0x40,
	0x14, 0xdd, 0x74, 0x1b, 0xeb, 0xde, 0x6e, 0xb7, 0xe9, 0x28, 0x34, 0x14, 0x0b, 0x32, 0x3e, 0x68,
	0x11, 0x17, 0xac, 0x4f, 0xe2, 0xd3, 0x56, 0x8b, 0x14, 0xb6, 0x56, 0xa6, 0x7d, 0x11, 0x1f, 0x42,
	0xdc, 0xb9, 0x96, 0x81, 0x64, 0x66, 0x99, 0x19, 0xab, 0xfe, 0x07, 0xff, 0x8e, 0x3f, 0xc9, 0xff,
	0xe1, 0xcc, 0xe4, 0xd3, 0x20, 0x85, 0x3e, 0x6d, 0x72, 0xee, 0xb9, 0xe7, 0xdc, 0x7b, 0x6e, 0x16,
	0x66, 0x4a, 0x5f, 0xe7, 0x52, 0x18, 0x9c, 0xaf, 0xb5, 0xb2, 0x8a, 0xc4, 0xe1, 0x87, 0x6e, 0x41,
	0x7c, 0x5a, 0xae, 0xed, 0x4f, 0xfa, 0x19, 0x12, 0x86, 0x05, 0xe6, 0x06, 0x3f, 0x16, 0xf9, 0x0a,
	0x4b, 0x94, 0x96, 0x1c, 0x02, 0xe8, 0x0a, 0xcb, 0x04, 0x4f, 0xa3, 0xc7, 0xd1, 0xb3, 0x98, 0x4d,
	0x6a, 0xe4, 0x8c, 0x93, 0x87, 0x10, 0x0b, 0xc9, 0xf1, 0x47, 0xba, 0x11, 0x2a, 0xd5, 0x0b, 0x21,
	0xb0, 0x69, 0x0a, 0x65, 0xd3, 0x71, 0x00, 0xc3, 0x33, 0xfd, 0x13, 0xc1, 0xfd, 0xa5, 0x5a, 0xe5,
	0x56, 0x28, 0xe9, 0x09, 0x32, 0x2f, 0x31, 0xe8, 0x4d, 0x58, 0x78, 0xf6, 0x52, 0xdf, 0xa4, 0xb0,
	0xa6, 0x91, 0x0a, 0x2f, 0xde, 0xff, 0xab, 0x2a, 0x38, 0x6a, 0x67, 0x6f, 0x9c, 0xe0, 0xd8, 0xfb,
	0x57, 0xc8, 0x19, 0x37, 0xe4, 0x1d, 0xec, 0xd5, 0xc3, 0x98, 0xac, 0xa8, 0xd5, 0xd3, 0x4d, 0xc7,
	0xda, 0x3e, 0xde, 0xaf, 0xb6, 0x9c, 0x0f, 0x57, 0x62, 0x49, 0xd3, 0xd1, 0x8e, 0xf3, 0xdc, 0xcd,
	0xab, 0xb4, 0x4d, 0x63, 0xe7, 0x3c, 0x6b, 0x1b, 0x9b, 0xf2, 0xfc, 0xd2, 0xd5, 0x84, 0xbc, 0x66,
	0x81, 0x44, 0x0f, 0x61, 0xab, 0x06, 0xdc, 0x1a, 0xb3, 0x93, 0x4f, 0xd9, 0x72, 0x71, 0x72, 0xba,
	0xcc, 0xde, 0x2e, 0xae, 0x3e, 0x5c, 0x24, 0x23, 0x17, 0xe2, 0xf4, 0xa2, 0x8a, 0xb9, 0xd2, 0x7e,
	0x04, 0x13, 0x2b, 0x4a, 0x34, 0x36, 0x2f, 0xd7, 0x61, 0xdf, 0x31, 0xeb, 0x00, 0xf2, 0x02, 0x26,
	0xcd, 0xd8, 0x7e, 0x71, 0x3f, 0xf7, 0xee, 0xc0, 0x9e, 0x75, 0x0c, 0x7a, 0x0e, 0x49, 0x5f, 0x7c,
	0x29, 0x8c, 0x25, 0xaf, 0x61, 0x47, 0xf5, 0x30, 0xe3, 0x4c, 0xbc, 0xcc, 0x83, 0x5a, 0xa6, 0xcf,
	0x67, 0xff, 0x32, 0x29, 0x87, 0x69, 0xe3, 0x72, 0xae, 0x6e, 0x90, 0x1c, 0xc1, 0xd8, 0x25, 0x1b,
	0xa6, 0xbc, 0x25, 0x3f, 0xcf, 0xf1, 0x54, 0x89, 0xdf, 0xc3, 0xad, 0x6e, 0xa3, 0x3a, 0x0e, 0xfd,
	0x15, 0xc1, 0x5e, 0x7f, 0x0a, 0x6f, 0x65, 0xc8, 0x53, 0xd8, 0x75, 0x11, 0x68, 0x9b, 0x0d, 0xd3,
	0x99, 0x05, 0xf8, 0xaa, 0x8d, 0xe8, 0x09, 0xec, 0xa0, 0xe4, 0x3d, 0xda, 0x46, 0xa0, 0x4d, 0x1d,
	0xd8, 0x91, 0x8e, 0x20, 0x2e, 0xbd, 0x6c, 0xf8, 0x42, 0xba, 0xe5, 0xfb, 0xdb, 0xb1, 0x8a, 0x71,
	0xfc, 0x3b, 0x6a, 0x43, 0x44, 0x7d, 0x89, 0xfa, 0x46, 0xac, 0x90, 0xbc, 0x84, 0xed, 0x05, 0xe7,
	0xed, 0x07, 0x31, 0xbc, 0xc1, 0xc1, 0x10, 0xa0, 0x23, 0xdf, 0xf2, 0x1e, 0xed, 0x9d, 0x5a, 0xde,
	0x40, 0xe2, 0x5a, 0xfa, 0x59, 0x18, 0x32, 0xad, 0x69, 0xe1, 0x2f, 0x78, 0xb0, 0xff, 0x9f, 0xab,
	0xf9, 0x2b, 0xd3, 0xd1, 0x97, 0x7b, 0xa1, 0xf2, 0xea, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe1,
	0xf5, 0xf6, 0x18, 0xc6, 0x03, 0x00, 0x00,
}
