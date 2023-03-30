package gateway

import (
	"bytes"
	"testing"

	"github.com/nlachfr/protoc-gen-cel-validate/testdata/validate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestDynamicpbCodecMarshal(t *testing.T) {
	tests := []struct {
		Name          string
		Desc          protoreflect.MessageDescriptor
		In            any
		WantedResults func() ([]byte, bool)
	}{
		{
			Name: "Invalid type",
			Desc: nil,
			In:   "string",
			WantedResults: func() ([]byte, bool) {
				return nil, true
			},
		},
		{
			Name: "Proto message",
			Desc: nil,
			In:   &validate.TestRpcRequest{},
			WantedResults: func() ([]byte, bool) {
				return []byte{}, false
			},
		},
		{
			Name: "Dynamicpb ptr with value",
			Desc: (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			In: func() **dynamicpb.Message {
				msg := dynamicpb.NewMessage((&validate.TestRpcRequest{}).ProtoReflect().Descriptor())
				msg.Set((&validate.TestRpcRequest{}).ProtoReflect().Descriptor().Fields().Get(0), protoreflect.ValueOf("myref"))
				return &msg
			}(),
			WantedResults: func() ([]byte, bool) {
				data, err := proto.Marshal(&validate.TestRpcRequest{Ref: "myref"})
				if err != nil {
					t.Error(err)
				}
				return data, false
			},
		},
		{
			Name: "Dynamicpb ptr with nil",
			Desc: (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			In: func() **dynamicpb.Message {
				var msg *dynamicpb.Message
				return &msg
			}(),
			WantedResults: func() ([]byte, bool) {
				return []byte{}, false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			wantRes, wantErr := tt.WantedResults()
			res, err := (&dynamicpbCodec{desc: tt.Desc, marshal: proto.Marshal}).Marshal(tt.In)
			if (wantErr && err == nil) || (!wantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", wantErr, err)
			} else if err == nil && !bytes.Equal(wantRes, res) {
				t.Errorf("wantRes %v, got %v", wantRes, res)
			}
		})
	}
}

func TestDynamicpbCodecUnmarshal(t *testing.T) {
	tests := []struct {
		Name    string
		Desc    protoreflect.MessageDescriptor
		Data    []byte
		Message any
		WantErr bool
	}{
		{
			Name:    "Invalid type",
			Desc:    nil,
			Message: "string",
			WantErr: true,
		},
		{
			Name:    "Proto message",
			Desc:    nil,
			Message: &validate.TestRpcRequest{},
			WantErr: false,
		},
		{
			Name: "Dynamicpb ptr with value",
			Desc: (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Data: func() []byte {
				data, err := proto.Marshal(&validate.TestRpcRequest{Ref: "myref"})
				if err != nil {
					t.Error(err)
				}
				return data
			}(),
			Message: func() interface{} {
				var msg *dynamicpb.Message
				return &msg
			}(),
		},
		{
			Name: "Dynamicpb ptr with nil",
			Desc: (&validate.TestRpcRequest{}).ProtoReflect().Descriptor(),
			Message: func() interface{} {
				var msg *dynamicpb.Message
				return &msg
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := (&dynamicpbCodec{desc: tt.Desc, unmarshal: proto.Unmarshal}).Unmarshal(tt.Data, tt.Message)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
