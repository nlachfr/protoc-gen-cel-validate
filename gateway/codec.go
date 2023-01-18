package gateway

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func buildClientCodecs(desc protoreflect.MessageDescriptor) []connect.ClientOption {
	return []connect.ClientOption{
		connect.WithCodec(&protoBinaryCodec{desc: desc}),
		connect.WithCodec(&protoJSONCodec{desc: desc}),
		connect.WithCodec(&protoJSONUTF8Codec{protoJSONCodec{desc: desc}}),
	}
}

func buildHandlerCodecs(desc protoreflect.MessageDescriptor) []connect.HandlerOption {
	return []connect.HandlerOption{
		connect.WithCodec(&protoBinaryCodec{desc: desc}),
		connect.WithCodec(&protoJSONCodec{desc: desc}),
		connect.WithCodec(&protoJSONUTF8Codec{protoJSONCodec{desc: desc}}),
	}
}

type protoBinaryCodec struct {
	desc protoreflect.MessageDescriptor
}

func (c *protoBinaryCodec) Name() string { return "proto" }

func (c *protoBinaryCodec) Marshal(message any) ([]byte, error) {
	if msg, ok := message.(*dynamicpb.Message); ok && msg == nil {
		return nil, nil
	}
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("%T doesn't implement proto.Message", message)
	}
	return proto.Marshal(protoMessage)
}

func (c *protoBinaryCodec) Unmarshal(data []byte, message any) error {
	if msg, ok := message.(**dynamicpb.Message); ok {
		new := dynamicpb.NewMessage(c.desc)
		*msg = new
		return proto.Unmarshal(data, *msg)
	}
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return fmt.Errorf("%T doesn't implement proto.Message", message)
	}
	return proto.Unmarshal(data, protoMessage)
}

type protoJSONCodec struct {
	desc protoreflect.MessageDescriptor
}

func (c *protoJSONCodec) Name() string { return "json" }

func (c *protoJSONCodec) Marshal(message any) ([]byte, error) {
	if msg, ok := message.(*dynamicpb.Message); ok && msg == nil {
		return nil, nil
	}
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("%T doesn't implement proto.Message", message)
	}
	return protojson.Marshal(protoMessage)
}

func (c *protoJSONCodec) Unmarshal(binary []byte, message any) error {
	if msg, ok := message.(**dynamicpb.Message); ok {
		new := dynamicpb.NewMessage(c.desc)
		*msg = new
		return protojson.Unmarshal(binary, *msg)
	}
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return fmt.Errorf("%T doesn't implement proto.Message", message)
	}
	return protojson.Unmarshal(binary, protoMessage)
}

type protoJSONUTF8Codec struct{ protoJSONCodec }

func (c *protoJSONUTF8Codec) Name() string { return "json; charset=utf-8" }
