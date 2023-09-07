// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.20.1
// source: internal/stores/pluginproto/plugin.proto

package pluginproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Objects struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Objects []string `protobuf:"bytes,1,rep,name=objects,proto3" json:"objects,omitempty"`
}

func (x *Objects) Reset() {
	*x = Objects{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_stores_pluginproto_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Objects) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Objects) ProtoMessage() {}

func (x *Objects) ProtoReflect() protoreflect.Message {
	mi := &file_internal_stores_pluginproto_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Objects.ProtoReflect.Descriptor instead.
func (*Objects) Descriptor() ([]byte, []int) {
	return file_internal_stores_pluginproto_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *Objects) GetObjects() []string {
	if x != nil {
		return x.Objects
	}
	return nil
}

type Options struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BackendAddress        string `protobuf:"bytes,1,opt,name=backend_address,json=backendAddress,proto3" json:"backend_address,omitempty"`
	PluginAddress         string `protobuf:"bytes,2,opt,name=plugin_address,json=pluginAddress,proto3" json:"plugin_address,omitempty"`
	MetadataFileExtension string `protobuf:"bytes,3,opt,name=metadata_file_extension,json=metadataFileExtension,proto3" json:"metadata_file_extension,omitempty"`
	Region                string `protobuf:"bytes,4,opt,name=region,proto3" json:"region,omitempty"`
}

func (x *Options) Reset() {
	*x = Options{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_stores_pluginproto_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Options) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Options) ProtoMessage() {}

func (x *Options) ProtoReflect() protoreflect.Message {
	mi := &file_internal_stores_pluginproto_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Options.ProtoReflect.Descriptor instead.
func (*Options) Descriptor() ([]byte, []int) {
	return file_internal_stores_pluginproto_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *Options) GetBackendAddress() string {
	if x != nil {
		return x.BackendAddress
	}
	return ""
}

func (x *Options) GetPluginAddress() string {
	if x != nil {
		return x.PluginAddress
	}
	return ""
}

func (x *Options) GetMetadataFileExtension() string {
	if x != nil {
		return x.MetadataFileExtension
	}
	return ""
}

func (x *Options) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

var File_internal_stores_pluginproto_plugin_proto protoreflect.FileDescriptor

var file_internal_stores_pluginproto_plugin_proto_rawDesc = []byte{
	0x0a, 0x28, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x73, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x23, 0x0a, 0x07, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x73, 0x22, 0xa9, 0x01, 0x0a, 0x07, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x27, 0x0a, 0x0f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x36, 0x0a, 0x17, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x15, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x32, 0xe6, 0x01, 0x0a, 0x06, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x33, 0x0a, 0x06, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x0f, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x35, 0x0a, 0x08, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x12, 0x0f, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x00,
	0x12, 0x37, 0x0a, 0x0a, 0x53, 0x65, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x0f,
	0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x65, 0x6e, 0x74, 0x65,
	0x6d, 0x2f, 0x63, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x73, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_stores_pluginproto_plugin_proto_rawDescOnce sync.Once
	file_internal_stores_pluginproto_plugin_proto_rawDescData = file_internal_stores_pluginproto_plugin_proto_rawDesc
)

func file_internal_stores_pluginproto_plugin_proto_rawDescGZIP() []byte {
	file_internal_stores_pluginproto_plugin_proto_rawDescOnce.Do(func() {
		file_internal_stores_pluginproto_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_stores_pluginproto_plugin_proto_rawDescData)
	})
	return file_internal_stores_pluginproto_plugin_proto_rawDescData
}

var file_internal_stores_pluginproto_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_stores_pluginproto_plugin_proto_goTypes = []interface{}{
	(*Objects)(nil),       // 0: plugin.Objects
	(*Options)(nil),       // 1: plugin.Options
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_internal_stores_pluginproto_plugin_proto_depIdxs = []int32{
	0, // 0: plugin.Plugin.Upload:input_type -> plugin.Objects
	0, // 1: plugin.Plugin.Retrieve:input_type -> plugin.Objects
	2, // 2: plugin.Plugin.GetOptions:input_type -> google.protobuf.Empty
	1, // 3: plugin.Plugin.SetOptions:input_type -> plugin.Options
	2, // 4: plugin.Plugin.Upload:output_type -> google.protobuf.Empty
	2, // 5: plugin.Plugin.Retrieve:output_type -> google.protobuf.Empty
	1, // 6: plugin.Plugin.GetOptions:output_type -> plugin.Options
	2, // 7: plugin.Plugin.SetOptions:output_type -> google.protobuf.Empty
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_stores_pluginproto_plugin_proto_init() }
func file_internal_stores_pluginproto_plugin_proto_init() {
	if File_internal_stores_pluginproto_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_stores_pluginproto_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Objects); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_stores_pluginproto_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Options); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_stores_pluginproto_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_stores_pluginproto_plugin_proto_goTypes,
		DependencyIndexes: file_internal_stores_pluginproto_plugin_proto_depIdxs,
		MessageInfos:      file_internal_stores_pluginproto_plugin_proto_msgTypes,
	}.Build()
	File_internal_stores_pluginproto_plugin_proto = out.File
	file_internal_stores_pluginproto_plugin_proto_rawDesc = nil
	file_internal_stores_pluginproto_plugin_proto_goTypes = nil
	file_internal_stores_pluginproto_plugin_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PluginClient is the client API for Plugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PluginClient interface {
	Upload(ctx context.Context, in *Objects, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Retrieve(ctx context.Context, in *Objects, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetOptions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Options, error)
	SetOptions(ctx context.Context, in *Options, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type pluginClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginClient(cc grpc.ClientConnInterface) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) Upload(ctx context.Context, in *Objects, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/plugin.Plugin/Upload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Retrieve(ctx context.Context, in *Objects, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/plugin.Plugin/Retrieve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetOptions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Options, error) {
	out := new(Options)
	err := c.cc.Invoke(ctx, "/plugin.Plugin/GetOptions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) SetOptions(ctx context.Context, in *Options, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/plugin.Plugin/SetOptions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
type PluginServer interface {
	Upload(context.Context, *Objects) (*emptypb.Empty, error)
	Retrieve(context.Context, *Objects) (*emptypb.Empty, error)
	GetOptions(context.Context, *emptypb.Empty) (*Options, error)
	SetOptions(context.Context, *Options) (*emptypb.Empty, error)
}

// UnimplementedPluginServer can be embedded to have forward compatible implementations.
type UnimplementedPluginServer struct {
}

func (*UnimplementedPluginServer) Upload(context.Context, *Objects) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (*UnimplementedPluginServer) Retrieve(context.Context, *Objects) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Retrieve not implemented")
}
func (*UnimplementedPluginServer) GetOptions(context.Context, *emptypb.Empty) (*Options, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOptions not implemented")
}
func (*UnimplementedPluginServer) SetOptions(context.Context, *Options) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOptions not implemented")
}

func RegisterPluginServer(s *grpc.Server, srv PluginServer) {
	s.RegisterService(&_Plugin_serviceDesc, srv)
}

func _Plugin_Upload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Objects)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Upload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Plugin/Upload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Upload(ctx, req.(*Objects))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Retrieve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Objects)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Retrieve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Plugin/Retrieve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Retrieve(ctx, req.(*Objects))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetOptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetOptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Plugin/GetOptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetOptions(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_SetOptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Options)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).SetOptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/plugin.Plugin/SetOptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).SetOptions(ctx, req.(*Options))
	}
	return interceptor(ctx, in, info, handler)
}

var _Plugin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "plugin.Plugin",
	HandlerType: (*PluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Upload",
			Handler:    _Plugin_Upload_Handler,
		},
		{
			MethodName: "Retrieve",
			Handler:    _Plugin_Retrieve_Handler,
		},
		{
			MethodName: "GetOptions",
			Handler:    _Plugin_GetOptions_Handler,
		},
		{
			MethodName: "SetOptions",
			Handler:    _Plugin_SetOptions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/stores/pluginproto/plugin.proto",
}
