package grpc

import (
	"context"
	"net"
	"testing"

	tvalidate "github.com/Neakxs/protocel/testdata/validate"
	"github.com/Neakxs/protocel/validate"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
)

type validateInterceptor struct {
	mapping map[string]*validate.Program
}

func (v *validateInterceptor) Validate(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error {
	return validate.NewValidateInterceptor(v.mapping).Validate(ctx, attr, m)
}

func TestNewGRPCUnaryInterceptor(t *testing.T) {
	tests := []struct {
		Name    string
		Expr    string
		Method  string
		Context context.Context
		Request *tvalidate.TestRpcRequest
		Info    *grpc.UnaryServerInfo
		WantErr bool
	}{
		{
			Name:    "Method provided correctly",
			Context: context.Background(),
			Expr:    `attribute_context.api.operation == "full.method"`,
			Method:  "full.method",
			Request: &tvalidate.TestRpcRequest{},
			Info:    &grpc.UnaryServerInfo{FullMethod: "full.method"},
			WantErr: false,
		},
		{
			Name: "Metadata provided correctly",
			Expr: `attribute_context.request.headers["x-metadata"] == "passed"`,
			Context: metadata.NewIncomingContext(context.Background(), metadata.MD{
				"x-metadata": []string{"passed"},
			}),
			Request: &tvalidate.TestRpcRequest{},
			Info:    &grpc.UnaryServerInfo{},
			WantErr: false,
		},
		{
			Name: "Peer provided correctly",
			Expr: `attribute_context.source.ip == "127.0.0.42" && attribute_context.source.port == 4242`,
			Context: peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.42"), Port: 4242},
			}),
			Request: &tvalidate.TestRpcRequest{},
			Info:    &grpc.UnaryServerInfo{},
			WantErr: false,
		},
		{
			Name:    "Request message provided correctly",
			Expr:    `ref == "ref"`,
			Context: context.Background(),
			Request: &tvalidate.TestRpcRequest{
				Ref: "ref",
			},
			Info:    &grpc.UnaryServerInfo{},
			WantErr: false,
		},
		{
			Name:    "Validatio error",
			Expr:    `ref != "ref"`,
			Context: context.Background(),
			Request: &tvalidate.TestRpcRequest{
				Ref: "ref",
			},
			Info:    &grpc.UnaryServerInfo{},
			WantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if pgr, err := validate.BuildMethodValidateProgram([]string{tt.Expr}, nil, tt.Request.ProtoReflect().Descriptor(), nil); err != nil {
				t.Error(err)
			} else if _, err = NewGRPCUnaryInterceptor(&validateInterceptor{mapping: map[string]*validate.Program{
				tt.Method: pgr,
			}})(tt.Context, tt.Request, tt.Info, func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil }); (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
