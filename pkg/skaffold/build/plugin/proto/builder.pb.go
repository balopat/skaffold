// Code generated by protoc-gen-go. DO NOT EDIT.
// source: builder.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	builder.proto

It has these top-level messages:
	LabelsRequest
	LabelsResponse
	DockerArtifact
	CustomArtifact
	Artifact
	BuildRequest
	BuildResult
	BuildResponse
	ValidationRequest
	ValidationResponse
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

type LabelsRequest struct {
}

func (m *LabelsRequest) Reset()                    { *m = LabelsRequest{} }
func (m *LabelsRequest) String() string            { return proto1.CompactTextString(m) }
func (*LabelsRequest) ProtoMessage()               {}
func (*LabelsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type LabelsResponse struct {
	Labels map[string]string `protobuf:"bytes,1,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *LabelsResponse) Reset()                    { *m = LabelsResponse{} }
func (m *LabelsResponse) String() string            { return proto1.CompactTextString(m) }
func (*LabelsResponse) ProtoMessage()               {}
func (*LabelsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *LabelsResponse) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type DockerArtifact struct {
	DockerfilePath string            `protobuf:"bytes,1,opt,name=DockerfilePath" json:"DockerfilePath,omitempty"`
	BuildArgs      map[string]string `protobuf:"bytes,2,rep,name=BuildArgs" json:"BuildArgs,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	CacheFrom      []string          `protobuf:"bytes,3,rep,name=CacheFrom" json:"CacheFrom,omitempty"`
	Target         string            `protobuf:"bytes,4,opt,name=Target" json:"Target,omitempty"`
}

func (m *DockerArtifact) Reset()                    { *m = DockerArtifact{} }
func (m *DockerArtifact) String() string            { return proto1.CompactTextString(m) }
func (*DockerArtifact) ProtoMessage()               {}
func (*DockerArtifact) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *DockerArtifact) GetDockerfilePath() string {
	if m != nil {
		return m.DockerfilePath
	}
	return ""
}

func (m *DockerArtifact) GetBuildArgs() map[string]string {
	if m != nil {
		return m.BuildArgs
	}
	return nil
}

func (m *DockerArtifact) GetCacheFrom() []string {
	if m != nil {
		return m.CacheFrom
	}
	return nil
}

func (m *DockerArtifact) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

type CustomArtifact struct {
	Configuration []byte `protobuf:"bytes,1,opt,name=configuration,proto3" json:"configuration,omitempty"`
}

func (m *CustomArtifact) Reset()                    { *m = CustomArtifact{} }
func (m *CustomArtifact) String() string            { return proto1.CompactTextString(m) }
func (*CustomArtifact) ProtoMessage()               {}
func (*CustomArtifact) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CustomArtifact) GetConfiguration() []byte {
	if m != nil {
		return m.Configuration
	}
	return nil
}

type Artifact struct {
	ImageName string            `protobuf:"bytes,1,opt,name=ImageName" json:"ImageName,omitempty" yaml`
	Workspace string            `protobuf:"bytes,2,opt,name=Workspace" json:"Workspace,omitempty"`
	Sync      map[string]string `protobuf:"bytes,3,rep,name=Sync" json:"Sync,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// Types that are valid to be assigned to Type:
	//	*Artifact_Docker
	//	*Artifact_Custom
	Type isArtifact_Type `protobuf_oneof:"type"`
}

func (m *Artifact) Reset()                    { *m = Artifact{} }
func (m *Artifact) String() string            { return proto1.CompactTextString(m) }
func (*Artifact) ProtoMessage()               {}
func (*Artifact) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type isArtifact_Type interface{ isArtifact_Type() }

type Artifact_Docker struct {
	Docker *DockerArtifact `protobuf:"bytes,4,opt,name=docker,oneof"`
}
type Artifact_Custom struct {
	Custom *CustomArtifact `protobuf:"bytes,5,opt,name=custom,oneof"`
}

func (*Artifact_Docker) isArtifact_Type() {}
func (*Artifact_Custom) isArtifact_Type() {}

func (m *Artifact) GetType() isArtifact_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Artifact) GetImageName() string {
	if m != nil {
		return m.ImageName
	}
	return ""
}

func (m *Artifact) GetWorkspace() string {
	if m != nil {
		return m.Workspace
	}
	return ""
}

func (m *Artifact) GetSync() map[string]string {
	if m != nil {
		return m.Sync
	}
	return nil
}

func (m *Artifact) GetDocker() *DockerArtifact {
	if x, ok := m.GetType().(*Artifact_Docker); ok {
		return x.Docker
	}
	return nil
}

func (m *Artifact) GetCustom() *CustomArtifact {
	if x, ok := m.GetType().(*Artifact_Custom); ok {
		return x.Custom
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Artifact) XXX_OneofFuncs() (func(msg proto1.Message, b *proto1.Buffer) error, func(msg proto1.Message, tag, wire int, b *proto1.Buffer) (bool, error), func(msg proto1.Message) (n int), []interface{}) {
	return _Artifact_OneofMarshaler, _Artifact_OneofUnmarshaler, _Artifact_OneofSizer, []interface{}{
		(*Artifact_Docker)(nil),
		(*Artifact_Custom)(nil),
	}
}

func _Artifact_OneofMarshaler(msg proto1.Message, b *proto1.Buffer) error {
	m := msg.(*Artifact)
	// type
	switch x := m.Type.(type) {
	case *Artifact_Docker:
		b.EncodeVarint(4<<3 | proto1.WireBytes)
		if err := b.EncodeMessage(x.Docker); err != nil {
			return err
		}
	case *Artifact_Custom:
		b.EncodeVarint(5<<3 | proto1.WireBytes)
		if err := b.EncodeMessage(x.Custom); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Artifact.Type has unexpected type %T", x)
	}
	return nil
}

func _Artifact_OneofUnmarshaler(msg proto1.Message, tag, wire int, b *proto1.Buffer) (bool, error) {
	m := msg.(*Artifact)
	switch tag {
	case 4: // type.docker
		if wire != proto1.WireBytes {
			return true, proto1.ErrInternalBadWireType
		}
		msg := new(DockerArtifact)
		err := b.DecodeMessage(msg)
		m.Type = &Artifact_Docker{msg}
		return true, err
	case 5: // type.custom
		if wire != proto1.WireBytes {
			return true, proto1.ErrInternalBadWireType
		}
		msg := new(CustomArtifact)
		err := b.DecodeMessage(msg)
		m.Type = &Artifact_Custom{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Artifact_OneofSizer(msg proto1.Message) (n int) {
	m := msg.(*Artifact)
	// type
	switch x := m.Type.(type) {
	case *Artifact_Docker:
		s := proto1.Size(x.Docker)
		n += proto1.SizeVarint(4<<3 | proto1.WireBytes)
		n += proto1.SizeVarint(uint64(s))
		n += s
	case *Artifact_Custom:
		s := proto1.Size(x.Custom)
		n += proto1.SizeVarint(5<<3 | proto1.WireBytes)
		n += proto1.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type BuildRequest struct {
	Artifacts []*Artifact `protobuf:"bytes,1,rep,name=artifacts" json:"artifacts,omitempty"`
}

func (m *BuildRequest) Reset()                    { *m = BuildRequest{} }
func (m *BuildRequest) String() string            { return proto1.CompactTextString(m) }
func (*BuildRequest) ProtoMessage()               {}
func (*BuildRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *BuildRequest) GetArtifacts() []*Artifact {
	if m != nil {
		return m.Artifacts
	}
	return nil
}

type BuildResult struct {
	ImageName string `protobuf:"bytes,1,opt,name=ImageName" json:"ImageName,omitempty"`
	Tag       string `protobuf:"bytes,2,opt,name=Tag" json:"Tag,omitempty"`
}

func (m *BuildResult) Reset()                    { *m = BuildResult{} }
func (m *BuildResult) String() string            { return proto1.CompactTextString(m) }
func (*BuildResult) ProtoMessage()               {}
func (*BuildResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *BuildResult) GetImageName() string {
	if m != nil {
		return m.ImageName
	}
	return ""
}

func (m *BuildResult) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

type BuildResponse struct {
	BuildResults []*BuildResult `protobuf:"bytes,1,rep,name=buildResults" json:"buildResults,omitempty"`
	Error        string         `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
}

func (m *BuildResponse) Reset()                    { *m = BuildResponse{} }
func (m *BuildResponse) String() string            { return proto1.CompactTextString(m) }
func (*BuildResponse) ProtoMessage()               {}
func (*BuildResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *BuildResponse) GetBuildResults() []*BuildResult {
	if m != nil {
		return m.BuildResults
	}
	return nil
}

func (m *BuildResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type ValidationRequest struct {
}

func (m *ValidationRequest) Reset()                    { *m = ValidationRequest{} }
func (m *ValidationRequest) String() string            { return proto1.CompactTextString(m) }
func (*ValidationRequest) ProtoMessage()               {}
func (*ValidationRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type ValidationResponse struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
}

func (m *ValidationResponse) Reset()                    { *m = ValidationResponse{} }
func (m *ValidationResponse) String() string            { return proto1.CompactTextString(m) }
func (*ValidationResponse) ProtoMessage()               {}
func (*ValidationResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ValidationResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto1.RegisterType((*LabelsRequest)(nil), "proto.LabelsRequest")
	proto1.RegisterType((*LabelsResponse)(nil), "proto.LabelsResponse")
	proto1.RegisterType((*DockerArtifact)(nil), "proto.DockerArtifact")
	proto1.RegisterType((*CustomArtifact)(nil), "proto.CustomArtifact")
	proto1.RegisterType((*Artifact)(nil), "proto.Artifact")
	proto1.RegisterType((*BuildRequest)(nil), "proto.BuildRequest")
	proto1.RegisterType((*BuildResult)(nil), "proto.BuildResult")
	proto1.RegisterType((*BuildResponse)(nil), "proto.BuildResponse")
	proto1.RegisterType((*ValidationRequest)(nil), "proto.ValidationRequest")
	proto1.RegisterType((*ValidationResponse)(nil), "proto.ValidationResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Builder service

type BuilderClient interface {
	Labels(ctx context.Context, in *LabelsRequest, opts ...grpc.CallOption) (*LabelsResponse, error)
	Build(ctx context.Context, in *BuildRequest, opts ...grpc.CallOption) (*BuildResponse, error)
	Validate(ctx context.Context, in *ValidationRequest, opts ...grpc.CallOption) (*ValidationResponse, error)
}

type builderClient struct {
	cc *grpc.ClientConn
}

func NewBuilderClient(cc *grpc.ClientConn) BuilderClient {
	return &builderClient{cc}
}

func (c *builderClient) Labels(ctx context.Context, in *LabelsRequest, opts ...grpc.CallOption) (*LabelsResponse, error) {
	out := new(LabelsResponse)
	err := grpc.Invoke(ctx, "/proto.Builder/Labels", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *builderClient) Build(ctx context.Context, in *BuildRequest, opts ...grpc.CallOption) (*BuildResponse, error) {
	out := new(BuildResponse)
	err := grpc.Invoke(ctx, "/proto.Builder/Build", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *builderClient) Validate(ctx context.Context, in *ValidationRequest, opts ...grpc.CallOption) (*ValidationResponse, error) {
	out := new(ValidationResponse)
	err := grpc.Invoke(ctx, "/proto.Builder/Validate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Builder service

type BuilderServer interface {
	Labels(context.Context, *LabelsRequest) (*LabelsResponse, error)
	Build(context.Context, *BuildRequest) (*BuildResponse, error)
	Validate(context.Context, *ValidationRequest) (*ValidationResponse, error)
}

func RegisterBuilderServer(s *grpc.Server, srv BuilderServer) {
	s.RegisterService(&_Builder_serviceDesc, srv)
}

func _Builder_Labels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LabelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BuilderServer).Labels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Builder/Labels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BuilderServer).Labels(ctx, req.(*LabelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Builder_Build_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuildRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BuilderServer).Build(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Builder/Build",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BuilderServer).Build(ctx, req.(*BuildRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Builder_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BuilderServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Builder/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BuilderServer).Validate(ctx, req.(*ValidationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Builder_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Builder",
	HandlerType: (*BuilderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Labels",
			Handler:    _Builder_Labels_Handler,
		},
		{
			MethodName: "Build",
			Handler:    _Builder_Build_Handler,
		},
		{
			MethodName: "Validate",
			Handler:    _Builder_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "builder.proto",
}

func init() { proto1.RegisterFile("builder.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 544 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x25, 0x4d, 0x1b, 0x96, 0x9b, 0xb6, 0x1b, 0x5e, 0x87, 0xb2, 0x68, 0x0f, 0x25, 0x9a, 0x50,
	0x85, 0xb4, 0x22, 0x05, 0x31, 0x18, 0x62, 0x42, 0xeb, 0x00, 0x81, 0x84, 0x10, 0x0a, 0x13, 0x3c,
	0xf1, 0xe0, 0xa6, 0x6e, 0x16, 0x35, 0x8d, 0x8b, 0xe3, 0x20, 0xf5, 0x07, 0xf8, 0x06, 0x7e, 0x84,
	0xff, 0xe2, 0x13, 0x50, 0x6c, 0x27, 0x69, 0x42, 0x25, 0xd4, 0xa7, 0xc5, 0xc7, 0xf7, 0x9c, 0x7b,
	0xce, 0x9d, 0x6f, 0xa1, 0x37, 0xcd, 0xa2, 0x78, 0x46, 0xd8, 0x78, 0xc5, 0x28, 0xa7, 0xa8, 0x23,
	0xfe, 0xb8, 0xfb, 0xd0, 0xfb, 0x80, 0xa7, 0x24, 0x4e, 0x7d, 0xf2, 0x3d, 0x23, 0x29, 0x77, 0x7f,
	0x6a, 0xd0, 0x2f, 0x90, 0x74, 0x45, 0x93, 0x94, 0xa0, 0x0b, 0x30, 0x62, 0x81, 0xd8, 0xda, 0x50,
	0x1f, 0x59, 0xde, 0x03, 0x29, 0x31, 0xae, 0x97, 0xa9, 0xe3, 0x9b, 0x84, 0xb3, 0xb5, 0xaf, 0x08,
	0xce, 0x05, 0x58, 0x1b, 0x30, 0x3a, 0x00, 0x7d, 0x41, 0xd6, 0xb6, 0x36, 0xd4, 0x46, 0xa6, 0x9f,
	0x7f, 0xa2, 0x01, 0x74, 0x7e, 0xe0, 0x38, 0x23, 0x76, 0x4b, 0x60, 0xf2, 0xf0, 0xa2, 0xf5, 0x5c,
	0x73, 0xff, 0x68, 0xd0, 0x7f, 0x4d, 0x83, 0x05, 0x61, 0x57, 0x8c, 0x47, 0x73, 0x1c, 0x70, 0xf4,
	0xb0, 0x40, 0xe6, 0x51, 0x4c, 0x3e, 0x61, 0x7e, 0xab, 0x94, 0x1a, 0x28, 0x9a, 0x80, 0x39, 0xc9,
	0xc3, 0x5e, 0xb1, 0x30, 0xb5, 0x5b, 0xc2, 0xf3, 0xa9, 0xf2, 0x5c, 0x57, 0x1c, 0x97, 0x65, 0xd2,
	0x76, 0x45, 0x43, 0x27, 0x60, 0x5e, 0xe3, 0xe0, 0x96, 0xbc, 0x65, 0x74, 0x69, 0xeb, 0x43, 0x7d,
	0x64, 0xfa, 0x15, 0x80, 0xee, 0x83, 0x71, 0x83, 0x59, 0x48, 0xb8, 0xdd, 0x16, 0x0e, 0xd4, 0xc9,
	0x79, 0x09, 0xfd, 0xba, 0xe4, 0x4e, 0x91, 0xcf, 0xa1, 0x7f, 0x9d, 0xa5, 0x9c, 0x2e, 0xcb, 0xc4,
	0xa7, 0xd0, 0x0b, 0x68, 0x32, 0x8f, 0xc2, 0x8c, 0x61, 0x1e, 0xd1, 0x44, 0xe8, 0x74, 0xfd, 0x3a,
	0xe8, 0xfe, 0x6a, 0xc1, 0x5e, 0x49, 0x39, 0x01, 0xf3, 0xfd, 0x12, 0x87, 0xe4, 0x23, 0x5e, 0x12,
	0xd5, 0xb6, 0x02, 0xf2, 0xdb, 0xaf, 0x94, 0x2d, 0xd2, 0x15, 0x0e, 0x0a, 0x03, 0x15, 0x80, 0xce,
	0xa0, 0xfd, 0x79, 0x9d, 0x04, 0x22, 0xaf, 0xe5, 0x1d, 0xab, 0x99, 0x95, 0xd3, 0xca, 0xef, 0xe4,
	0xa0, 0x44, 0x19, 0x7a, 0x0c, 0xc6, 0x4c, 0xcc, 0x53, 0x4c, 0xc1, 0xf2, 0x8e, 0xb6, 0x0e, 0xf9,
	0xdd, 0x1d, 0x5f, 0x95, 0xe5, 0x84, 0x40, 0x04, 0xb4, 0x3b, 0x35, 0x42, 0x3d, 0x75, 0x4e, 0x90,
	0x65, 0xce, 0x33, 0x30, 0xcb, 0xa6, 0xbb, 0x8c, 0x72, 0x62, 0x40, 0x9b, 0xaf, 0x57, 0xc4, 0xbd,
	0x84, 0xae, 0xf8, 0x87, 0xa8, 0xe7, 0x8d, 0xce, 0xc0, 0xc4, 0xaa, 0x4d, 0xf1, 0x9c, 0xf7, 0x1b,
	0x31, 0xfd, 0xaa, 0xc2, 0xbd, 0x04, 0x4b, 0xd1, 0xd3, 0x2c, 0xfe, 0xdf, 0x6c, 0x0f, 0x40, 0xbf,
	0xc1, 0xa1, 0xf2, 0x92, 0x7f, 0xba, 0xdf, 0xa0, 0x57, 0xd0, 0xe5, 0x2a, 0x9d, 0x43, 0x77, 0x5a,
	0xe9, 0x15, 0x0e, 0x90, 0x72, 0xb0, 0xd1, 0xca, 0xaf, 0xd5, 0xe5, 0x41, 0x09, 0x63, 0x94, 0x15,
	0x41, 0xc5, 0xc1, 0x3d, 0x84, 0x7b, 0x5f, 0x70, 0x1c, 0xcd, 0xc4, 0x2b, 0x28, 0x16, 0xf8, 0x11,
	0xa0, 0x4d, 0x50, 0x35, 0x2e, 0x05, 0xb4, 0x0d, 0x01, 0xef, 0xb7, 0x06, 0x77, 0x27, 0xf2, 0x67,
	0x01, 0x3d, 0x05, 0x43, 0xae, 0x2a, 0x1a, 0x34, 0xf6, 0x5b, 0xe8, 0x3a, 0x47, 0x5b, 0xb7, 0x1e,
	0x79, 0xd0, 0x11, 0x0a, 0xe8, 0xb0, 0x1e, 0x42, 0x92, 0x06, 0x8d, 0x64, 0x92, 0xf3, 0x0a, 0xf6,
	0x94, 0x45, 0x82, 0x6c, 0x55, 0xf1, 0x4f, 0x10, 0xe7, 0x78, 0xcb, 0x8d, 0x14, 0x98, 0x1a, 0xe2,
	0xe6, 0xc9, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xae, 0xf9, 0xd9, 0x98, 0xd4, 0x04, 0x00, 0x00,
}
