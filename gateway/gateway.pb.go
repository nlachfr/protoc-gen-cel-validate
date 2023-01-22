// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: gateway/gateway.proto

package gateway

import (
	validate "github.com/Neakxs/protocel/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Configuration_Server_Upstream_Protocol int32

const (
	Configuration_Server_Upstream_GRPC     Configuration_Server_Upstream_Protocol = 0
	Configuration_Server_Upstream_GRPC_WEB Configuration_Server_Upstream_Protocol = 1
	Configuration_Server_Upstream_CONNECT  Configuration_Server_Upstream_Protocol = 2
)

// Enum value maps for Configuration_Server_Upstream_Protocol.
var (
	Configuration_Server_Upstream_Protocol_name = map[int32]string{
		0: "GRPC",
		1: "GRPC_WEB",
		2: "CONNECT",
	}
	Configuration_Server_Upstream_Protocol_value = map[string]int32{
		"GRPC":     0,
		"GRPC_WEB": 1,
		"CONNECT":  2,
	}
)

func (x Configuration_Server_Upstream_Protocol) Enum() *Configuration_Server_Upstream_Protocol {
	p := new(Configuration_Server_Upstream_Protocol)
	*p = x
	return p
}

func (x Configuration_Server_Upstream_Protocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Configuration_Server_Upstream_Protocol) Descriptor() protoreflect.EnumDescriptor {
	return file_gateway_gateway_proto_enumTypes[0].Descriptor()
}

func (Configuration_Server_Upstream_Protocol) Type() protoreflect.EnumType {
	return &file_gateway_gateway_proto_enumTypes[0]
}

func (x Configuration_Server_Upstream_Protocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Configuration_Server_Upstream_Protocol.Descriptor instead.
func (Configuration_Server_Upstream_Protocol) EnumDescriptor() ([]byte, []int) {
	return file_gateway_gateway_proto_rawDescGZIP(), []int{0, 0, 0, 0}
}

type Configuration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Servers  []*Configuration_Server `protobuf:"bytes,1,rep,name=servers,proto3" json:"servers,omitempty"`
	Files    *Configuration_Files    `protobuf:"bytes,2,opt,name=files,proto3" json:"files,omitempty"`
	Validate *validate.Options       `protobuf:"bytes,10,opt,name=validate,proto3" json:"validate,omitempty"`
}

func (x *Configuration) Reset() {
	*x = Configuration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_gateway_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration) ProtoMessage() {}

func (x *Configuration) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_gateway_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration.ProtoReflect.Descriptor instead.
func (*Configuration) Descriptor() ([]byte, []int) {
	return file_gateway_gateway_proto_rawDescGZIP(), []int{0}
}

func (x *Configuration) GetServers() []*Configuration_Server {
	if x != nil {
		return x.Servers
	}
	return nil
}

func (x *Configuration) GetFiles() *Configuration_Files {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *Configuration) GetValidate() *validate.Options {
	if x != nil {
		return x.Validate
	}
	return nil
}

type Configuration_Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Listen    []string                                  `protobuf:"bytes,1,rep,name=listen,proto3" json:"listen,omitempty"`
	Upstreams map[string]*Configuration_Server_Upstream `protobuf:"bytes,2,rep,name=upstreams,proto3" json:"upstreams,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Configuration_Server) Reset() {
	*x = Configuration_Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_gateway_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration_Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration_Server) ProtoMessage() {}

func (x *Configuration_Server) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_gateway_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration_Server.ProtoReflect.Descriptor instead.
func (*Configuration_Server) Descriptor() ([]byte, []int) {
	return file_gateway_gateway_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Configuration_Server) GetListen() []string {
	if x != nil {
		return x.Listen
	}
	return nil
}

func (x *Configuration_Server) GetUpstreams() map[string]*Configuration_Server_Upstream {
	if x != nil {
		return x.Upstreams
	}
	return nil
}

type Configuration_Files struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sources []string `protobuf:"bytes,1,rep,name=sources,proto3" json:"sources,omitempty"`
	Imports []string `protobuf:"bytes,2,rep,name=imports,proto3" json:"imports,omitempty"`
}

func (x *Configuration_Files) Reset() {
	*x = Configuration_Files{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_gateway_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration_Files) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration_Files) ProtoMessage() {}

func (x *Configuration_Files) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_gateway_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration_Files.ProtoReflect.Descriptor instead.
func (*Configuration_Files) Descriptor() ([]byte, []int) {
	return file_gateway_gateway_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Configuration_Files) GetSources() []string {
	if x != nil {
		return x.Sources
	}
	return nil
}

func (x *Configuration_Files) GetImports() []string {
	if x != nil {
		return x.Imports
	}
	return nil
}

type Configuration_Server_Upstream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address  string                                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Server   string                                 `protobuf:"bytes,2,opt,name=server,proto3" json:"server,omitempty"`
	Protocol Configuration_Server_Upstream_Protocol `protobuf:"varint,3,opt,name=protocol,proto3,enum=protocel.gateway.Configuration_Server_Upstream_Protocol" json:"protocol,omitempty"`
}

func (x *Configuration_Server_Upstream) Reset() {
	*x = Configuration_Server_Upstream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_gateway_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration_Server_Upstream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration_Server_Upstream) ProtoMessage() {}

func (x *Configuration_Server_Upstream) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_gateway_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration_Server_Upstream.ProtoReflect.Descriptor instead.
func (*Configuration_Server_Upstream) Descriptor() ([]byte, []int) {
	return file_gateway_gateway_proto_rawDescGZIP(), []int{0, 0, 0}
}

func (x *Configuration_Server_Upstream) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Configuration_Server_Upstream) GetServer() string {
	if x != nil {
		return x.Server
	}
	return ""
}

func (x *Configuration_Server_Upstream) GetProtocol() Configuration_Server_Upstream_Protocol {
	if x != nil {
		return x.Protocol
	}
	return Configuration_Server_Upstream_GRPC
}

var File_gateway_gateway_proto protoreflect.FileDescriptor

var file_gateway_gateway_proto_rawDesc = []byte{
	0x0a, 0x15, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x65,
	0x6c, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xb0, 0x05, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x65, 0x6c,
	0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x07, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x12, 0x3b, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x65, 0x6c,
	0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x05, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x65, 0x6c,
	0x2e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x1a, 0xaa, 0x03, 0x0a, 0x06,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x12, 0x53,
	0x0a, 0x09, 0x75, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x35, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x65, 0x6c, 0x2e, 0x67, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x75, 0x70, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x73, 0x1a, 0xc3, 0x01, 0x0a, 0x08, 0x55, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x12, 0x54, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x38, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x65, 0x6c, 0x2e,
	0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x52, 0x08,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x2f, 0x0a, 0x08, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x08, 0x0a, 0x04, 0x47, 0x52, 0x50, 0x43, 0x10, 0x00, 0x12, 0x0c,
	0x0a, 0x08, 0x47, 0x52, 0x50, 0x43, 0x5f, 0x57, 0x45, 0x42, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07,
	0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x02, 0x1a, 0x6d, 0x0a, 0x0e, 0x55, 0x70, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x45, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x65, 0x6c, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3b, 0x0a, 0x05, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x69,
	0x6d, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x69, 0x6d,
	0x70, 0x6f, 0x72, 0x74, 0x73, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x4e, 0x65, 0x61, 0x6b, 0x78, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x65, 0x6c, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_gateway_gateway_proto_rawDescOnce sync.Once
	file_gateway_gateway_proto_rawDescData = file_gateway_gateway_proto_rawDesc
)

func file_gateway_gateway_proto_rawDescGZIP() []byte {
	file_gateway_gateway_proto_rawDescOnce.Do(func() {
		file_gateway_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(file_gateway_gateway_proto_rawDescData)
	})
	return file_gateway_gateway_proto_rawDescData
}

var file_gateway_gateway_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gateway_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_gateway_gateway_proto_goTypes = []interface{}{
	(Configuration_Server_Upstream_Protocol)(0), // 0: protocel.gateway.Configuration.Server.Upstream.Protocol
	(*Configuration)(nil),                       // 1: protocel.gateway.Configuration
	(*Configuration_Server)(nil),                // 2: protocel.gateway.Configuration.Server
	(*Configuration_Files)(nil),                 // 3: protocel.gateway.Configuration.Files
	(*Configuration_Server_Upstream)(nil),       // 4: protocel.gateway.Configuration.Server.Upstream
	nil,                                         // 5: protocel.gateway.Configuration.Server.UpstreamsEntry
	(*validate.Options)(nil),                    // 6: protocel.validate.Options
}
var file_gateway_gateway_proto_depIdxs = []int32{
	2, // 0: protocel.gateway.Configuration.servers:type_name -> protocel.gateway.Configuration.Server
	3, // 1: protocel.gateway.Configuration.files:type_name -> protocel.gateway.Configuration.Files
	6, // 2: protocel.gateway.Configuration.validate:type_name -> protocel.validate.Options
	5, // 3: protocel.gateway.Configuration.Server.upstreams:type_name -> protocel.gateway.Configuration.Server.UpstreamsEntry
	0, // 4: protocel.gateway.Configuration.Server.Upstream.protocol:type_name -> protocel.gateway.Configuration.Server.Upstream.Protocol
	4, // 5: protocel.gateway.Configuration.Server.UpstreamsEntry.value:type_name -> protocel.gateway.Configuration.Server.Upstream
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_gateway_gateway_proto_init() }
func file_gateway_gateway_proto_init() {
	if File_gateway_gateway_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gateway_gateway_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration); i {
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
		file_gateway_gateway_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration_Server); i {
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
		file_gateway_gateway_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration_Files); i {
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
		file_gateway_gateway_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration_Server_Upstream); i {
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
			RawDescriptor: file_gateway_gateway_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gateway_gateway_proto_goTypes,
		DependencyIndexes: file_gateway_gateway_proto_depIdxs,
		EnumInfos:         file_gateway_gateway_proto_enumTypes,
		MessageInfos:      file_gateway_gateway_proto_msgTypes,
	}.Build()
	File_gateway_gateway_proto = out.File
	file_gateway_gateway_proto_rawDesc = nil
	file_gateway_gateway_proto_goTypes = nil
	file_gateway_gateway_proto_depIdxs = nil
}
