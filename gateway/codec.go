package gateway

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func newCodecs(desc protoreflect.MessageDescriptor) connect.Option {
	return connect.WithOptions(
		connect.WithCodec(&dynamicpbCodec{
			name:      "proto",
			marshal:   proto.Marshal,
			unmarshal: proto.Unmarshal,
			desc:      desc,
		}),
		connect.WithCodec(&dynamicpbCodec{
			name:      "json",
			marshal:   protojson.Marshal,
			unmarshal: protojson.Unmarshal,
			desc:      desc,
		}),
		connect.WithCodec(&dynamicpbCodec{
			name:      "json; charset=utf-8",
			marshal:   protojson.Marshal,
			unmarshal: protojson.Unmarshal,
			desc:      desc,
		}),
	)
}

type dynamicpbCodec struct {
	name      string
	marshal   func(proto.Message) ([]byte, error)
	unmarshal func([]byte, proto.Message) error

	desc protoreflect.MessageDescriptor
}

func (c *dynamicpbCodec) Name() string { return c.name }

func (c *dynamicpbCodec) Marshal(message any) ([]byte, error) {
	if msg, ok := message.(**dynamicpb.Message); ok {
		if *msg == nil {
			return nil, nil
		}
		return c.marshal(*msg)
	} else if msg, ok := message.(proto.Message); ok {
		return c.marshal(msg)
	}
	return nil, fmt.Errorf("marshal error: invalid message type: %T", message)
}

func (c *dynamicpbCodec) Unmarshal(data []byte, message any) error {
	if msg, ok := message.(**dynamicpb.Message); ok {
		*msg = dynamicpb.NewMessage(c.desc)
		return c.unmarshal(data, *msg)
	} else if msg, ok := message.(proto.Message); ok {
		return c.unmarshal(data, msg)
	}
	return fmt.Errorf("unmarshal error: invalid message type: %T", message)
}
