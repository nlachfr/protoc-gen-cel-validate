package gateway

import (
	"context"
	"fmt"
	"testing"

	"github.com/bufbuild/connect-go"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
)

type mockServiceRuleValidater struct {
	validate func(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error
}

func (v *mockServiceRuleValidater) Validate(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error {
	return v.validate(ctx, attr, m)
}

type mockConnectClient struct {
	unary        func(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.Response[*dynamicpb.Message], error)
	clientStream func(ctx context.Context) *connect.ClientStreamForClient[*dynamicpb.Message, *dynamicpb.Message]
	serverStream func(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.ServerStreamForClient[*dynamicpb.Message], error)
	bidiStream   func(ctx context.Context) *connect.BidiStreamForClient[*dynamicpb.Message, *dynamicpb.Message]
}

func (c *mockConnectClient) CallUnary(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.Response[*dynamicpb.Message], error) {
	return c.unary(ctx, request)
}

func (c *mockConnectClient) CallClientStream(ctx context.Context) *connect.ClientStreamForClient[*dynamicpb.Message, *dynamicpb.Message] {
	return c.clientStream(ctx)
}

func (c *mockConnectClient) CallServerStream(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.ServerStreamForClient[*dynamicpb.Message], error) {
	return c.serverStream(ctx, request)
}

func (c *mockConnectClient) CallBidiStream(ctx context.Context) *connect.BidiStreamForClient[*dynamicpb.Message, *dynamicpb.Message] {
	return c.bidiStream(ctx)
}

func TestMethodHandlerUnary(t *testing.T) {
	m := dynamicpb.NewMessage(nil)
	req := connect.NewRequest(&m)
	tests := []struct {
		Name    string
		Handler *methodHandler
		WantErr bool
	}{
		{
			Name: "Validation error",
			Handler: &methodHandler{
				srv: &mockServiceRuleValidater{
					validate: func(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error {
						return fmt.Errorf("error")
					},
				},
				client: &mockConnectClient{
					unary: func(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.Response[*dynamicpb.Message], error) {
						return connect.NewResponse[*dynamicpb.Message](nil), nil
					},
				},
			},
			WantErr: true,
		},
		{
			Name: "Client error",
			Handler: &methodHandler{
				srv: &mockServiceRuleValidater{
					validate: func(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error {
						return nil
					},
				},
				client: &mockConnectClient{
					unary: func(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.Response[*dynamicpb.Message], error) {
						return connect.NewResponse[*dynamicpb.Message](nil), fmt.Errorf("error")
					},
				},
			},
			WantErr: true,
		},
		{
			Name: "OK",
			Handler: &methodHandler{
				srv: &mockServiceRuleValidater{
					validate: func(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error {
						return nil
					},
				},
				client: &mockConnectClient{
					unary: func(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.Response[*dynamicpb.Message], error) {
						return connect.NewResponse[*dynamicpb.Message](nil), nil
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := tt.Handler.unary(context.Background(), req)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
