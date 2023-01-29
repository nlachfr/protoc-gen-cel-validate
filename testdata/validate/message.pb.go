// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: testdata/validate/message.proto

package validate

import (
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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_validate_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_validate_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_testdata_validate_message_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type MessageExpr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *MessageExpr) Reset() {
	*x = MessageExpr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_validate_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageExpr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageExpr) ProtoMessage() {}

func (x *MessageExpr) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_validate_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageExpr.ProtoReflect.Descriptor instead.
func (*MessageExpr) Descriptor() ([]byte, []int) {
	return file_testdata_validate_message_proto_rawDescGZIP(), []int{1}
}

func (x *MessageExpr) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type MessageNested struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageExpr *MessageExpr `protobuf:"bytes,1,opt,name=message_expr,json=messageExpr,proto3" json:"message_expr,omitempty"`
}

func (x *MessageNested) Reset() {
	*x = MessageNested{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_validate_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageNested) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageNested) ProtoMessage() {}

func (x *MessageNested) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_validate_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageNested.ProtoReflect.Descriptor instead.
func (*MessageNested) Descriptor() ([]byte, []int) {
	return file_testdata_validate_message_proto_rawDescGZIP(), []int{2}
}

func (x *MessageNested) GetMessageExpr() *MessageExpr {
	if x != nil {
		return x.MessageExpr
	}
	return nil
}

type MessageNestedExpr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageExpr *MessageExpr `protobuf:"bytes,1,opt,name=message_expr,json=messageExpr,proto3" json:"message_expr,omitempty"`
}

func (x *MessageNestedExpr) Reset() {
	*x = MessageNestedExpr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_validate_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageNestedExpr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageNestedExpr) ProtoMessage() {}

func (x *MessageNestedExpr) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_validate_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageNestedExpr.ProtoReflect.Descriptor instead.
func (*MessageNestedExpr) Descriptor() ([]byte, []int) {
	return file_testdata_validate_message_proto_rawDescGZIP(), []int{3}
}

func (x *MessageNestedExpr) GetMessageExpr() *MessageExpr {
	if x != nil {
		return x.MessageExpr
	}
	return nil
}

type MessageOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *MessageOptions) Reset() {
	*x = MessageOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_validate_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOptions) ProtoMessage() {}

func (x *MessageOptions) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_validate_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOptions.ProtoReflect.Descriptor instead.
func (*MessageOptions) Descriptor() ([]byte, []int) {
	return file_testdata_validate_message_proto_rawDescGZIP(), []int{4}
}

func (x *MessageOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type MessageLocalOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *MessageLocalOptions) Reset() {
	*x = MessageLocalOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testdata_validate_message_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageLocalOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageLocalOptions) ProtoMessage() {}

func (x *MessageLocalOptions) ProtoReflect() protoreflect.Message {
	mi := &file_testdata_validate_message_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageLocalOptions.ProtoReflect.Descriptor instead.
func (*MessageLocalOptions) Descriptor() ([]byte, []int) {
	return file_testdata_validate_message_proto_rawDescGZIP(), []int{5}
}

func (x *MessageLocalOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_testdata_validate_message_proto protoreflect.FileDescriptor

var file_testdata_validate_message_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x11, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1d, 0x0a,
	0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x36, 0x0a, 0x0b,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x78, 0x70, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x3a,
	0x13, 0xd2, 0x49, 0x10, 0x12, 0x0e, 0x12, 0x0c, 0x12, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x20, 0x21,
	0x3d, 0x20, 0x22, 0x22, 0x22, 0x52, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e,
	0x65, 0x73, 0x74, 0x65, 0x64, 0x12, 0x41, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x5f, 0x65, 0x78, 0x70, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x74, 0x65,
	0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x78, 0x70, 0x72, 0x52, 0x0b, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x45, 0x78, 0x70, 0x72, 0x22, 0x78, 0x0a, 0x11, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x45, 0x78, 0x70, 0x72, 0x12, 0x41, 0x0a,
	0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x65, 0x78, 0x70, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45,
	0x78, 0x70, 0x72, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x78, 0x70, 0x72,
	0x3a, 0x20, 0xd2, 0x49, 0x1d, 0x12, 0x1b, 0x12, 0x19, 0x12, 0x17, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x65, 0x78, 0x70, 0x72, 0x2e, 0x6e, 0x61, 0x6d, 0x65, 0x20, 0x21, 0x3d, 0x20,
	0x22, 0x22, 0x22, 0x40, 0x0a, 0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x3a, 0x1a, 0xd2, 0x49, 0x17, 0x12, 0x15, 0x12,
	0x13, 0x12, 0x11, 0x6e, 0x61, 0x6d, 0x65, 0x20, 0x21, 0x3d, 0x20, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x58, 0x0a, 0x13, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4c,
	0x6f, 0x63, 0x61, 0x6c, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x3a,
	0x2d, 0xd2, 0x49, 0x2a, 0x12, 0x28, 0x0a, 0x11, 0x0a, 0x0f, 0x12, 0x0d, 0x0a, 0x09, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x00, 0x12, 0x13, 0x12, 0x11, 0x6e, 0x61, 0x6d,
	0x65, 0x20, 0x21, 0x3d, 0x20, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x2f,
	0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x6c, 0x61,
	0x63, 0x68, 0x66, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x65, 0x6c, 0x2f, 0x74, 0x65,
	0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_testdata_validate_message_proto_rawDescOnce sync.Once
	file_testdata_validate_message_proto_rawDescData = file_testdata_validate_message_proto_rawDesc
)

func file_testdata_validate_message_proto_rawDescGZIP() []byte {
	file_testdata_validate_message_proto_rawDescOnce.Do(func() {
		file_testdata_validate_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_testdata_validate_message_proto_rawDescData)
	})
	return file_testdata_validate_message_proto_rawDescData
}

var file_testdata_validate_message_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_testdata_validate_message_proto_goTypes = []interface{}{
	(*Message)(nil),             // 0: testdata.validate.Message
	(*MessageExpr)(nil),         // 1: testdata.validate.MessageExpr
	(*MessageNested)(nil),       // 2: testdata.validate.MessageNested
	(*MessageNestedExpr)(nil),   // 3: testdata.validate.MessageNestedExpr
	(*MessageOptions)(nil),      // 4: testdata.validate.MessageOptions
	(*MessageLocalOptions)(nil), // 5: testdata.validate.MessageLocalOptions
}
var file_testdata_validate_message_proto_depIdxs = []int32{
	1, // 0: testdata.validate.MessageNested.message_expr:type_name -> testdata.validate.MessageExpr
	1, // 1: testdata.validate.MessageNestedExpr.message_expr:type_name -> testdata.validate.MessageExpr
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_testdata_validate_message_proto_init() }
func file_testdata_validate_message_proto_init() {
	if File_testdata_validate_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_testdata_validate_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_testdata_validate_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageExpr); i {
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
		file_testdata_validate_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageNested); i {
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
		file_testdata_validate_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageNestedExpr); i {
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
		file_testdata_validate_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOptions); i {
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
		file_testdata_validate_message_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageLocalOptions); i {
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
			RawDescriptor: file_testdata_validate_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_testdata_validate_message_proto_goTypes,
		DependencyIndexes: file_testdata_validate_message_proto_depIdxs,
		MessageInfos:      file_testdata_validate_message_proto_msgTypes,
	}.Build()
	File_testdata_validate_message_proto = out.File
	file_testdata_validate_message_proto_rawDesc = nil
	file_testdata_validate_message_proto_goTypes = nil
	file_testdata_validate_message_proto_depIdxs = nil
}
