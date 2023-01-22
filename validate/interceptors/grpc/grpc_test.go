package grpc

import (
	"context"
	"net"
	"testing"

	testdata "github.com/nlachfr/protocel/testdata/validate"
	"github.com/nlachfr/protocel/validate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewGRPCUnaryInterceptor(t *testing.T) {
	tests := []struct {
		Name    string
		Desc    protoreflect.MethodDescriptor
		Context context.Context
		Request proto.Message
		Info    *grpc.UnaryServerInfo
		WantErr bool
	}{
		{
			Name:    "Missing header",
			Desc:    testdata.File_testdata_validate_service_proto.Services().ByName("ServiceExpr").Methods().ByName("Rpc"),
			Context: context.Background(),
			Request: &emptypb.Empty{},
			Info: &grpc.UnaryServerInfo{
				FullMethod: "/testdata.validate.ServiceExpr/Rpc",
			},
			WantErr: true,
		},
		{
			Name: "Wrong header value",
			Desc: testdata.File_testdata_validate_service_proto.Services().ByName("ServiceExpr").Methods().ByName("Rpc"),
			Context: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"x-is-admin": "false",
			})),
			Request: &emptypb.Empty{},
			Info: &grpc.UnaryServerInfo{
				FullMethod: "/testdata.validate.ServiceExpr/Rpc",
			},
			WantErr: true,
		},
		{
			Name: "Good header value",
			Desc: testdata.File_testdata_validate_service_proto.Services().ByName("ServiceExpr").Methods().ByName("Rpc"),
			Context: metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
				"x-is-admin": "true",
			})),
			Request: &emptypb.Empty{},
			Info: &grpc.UnaryServerInfo{
				FullMethod: "/testdata.validate.ServiceExpr/Rpc",
			},
			WantErr: false,
		},
		{
			Name:    "Missing peer context",
			Desc:    testdata.File_testdata_validate_service_proto.Services().ByName("ServicePeerExpr").Methods().ByName("Rpc"),
			Context: context.Background(),
			Request: &emptypb.Empty{},
			Info: &grpc.UnaryServerInfo{
				FullMethod: "/testdata.validate.ServicePeerExpr/Rpc",
			},
			WantErr: true,
		},
		{
			Name: "Wrong peer context",
			Desc: testdata.File_testdata_validate_service_proto.Services().ByName("ServicePeerExpr").Methods().ByName("Rpc"),
			Context: peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.IPAddr{IP: net.ParseIP("1.1.1.1")},
			}),
			Request: &emptypb.Empty{},
			Info: &grpc.UnaryServerInfo{
				FullMethod: "/testdata.validate.ServicePeerExpr/Rpc",
			},
			WantErr: true,
		},
		{
			Name: "Wrong peer context (tcp)",
			Desc: testdata.File_testdata_validate_service_proto.Services().ByName("ServicePeerExpr").Methods().ByName("Rpc"),
			Context: peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.IPAddr{IP: net.ParseIP("1.1.1.1")},
			}),
			Request: &emptypb.Empty{},
			Info: &grpc.UnaryServerInfo{
				FullMethod: "/testdata.validate.ServicePeerExpr/Rpc",
			},
			WantErr: true,
		},
		{
			Name: "Wrong peer context (udp)",
			Desc: testdata.File_testdata_validate_service_proto.Services().ByName("ServicePeerExpr").Methods().ByName("Rpc"),
			Context: peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.TCPAddr{IP: net.ParseIP("1.1.1.1")},
			}),
			Request: &emptypb.Empty{},
			Info: &grpc.UnaryServerInfo{
				FullMethod: "/testdata.validate.ServicePeerExpr/Rpc",
			},
			WantErr: true,
		},
		{
			Name: "Good peer context",
			Desc: testdata.File_testdata_validate_service_proto.Services().ByName("ServicePeerExpr").Methods().ByName("Rpc"),
			Context: peer.NewContext(context.Background(), &peer.Peer{
				Addr: &net.UDPAddr{IP: net.ParseIP("127.0.0.1")},
			}),
			Request: &emptypb.Empty{},
			Info: &grpc.UnaryServerInfo{
				FullMethod: "/testdata.validate.ServicePeerExpr/Rpc",
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			manager, err := validate.NewManager(tt.Desc.ParentFile())
			if err != nil {
				t.Error(err)
			}
			validater, err := manager.GetServiceRuleValidater(tt.Desc.Parent().(protoreflect.ServiceDescriptor))
			if err != nil {
				t.Error(err)
			}
			if _, err = NewGRPCUnaryInterceptor(validater, nil)(tt.Context, tt.Request, tt.Info, func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil }); (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
