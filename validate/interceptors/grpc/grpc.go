package grpc

import (
	"context"
	"net"
	"strings"

	"github.com/Neakxs/protocel/validate"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func NewGRPCUnaryInterceptor(validateHandler validate.ServiceValidateProgram, errorHandler func(err error) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		attr := &attribute_context.AttributeContext{
			Api: &attribute_context.AttributeContext_Api{
				Operation: info.FullMethod,
				Protocol:  "grpc",
			},
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			rq := &attribute_context.AttributeContext_Request{Headers: map[string]string{}}
			for k, v := range md {
				rq.Headers[strings.ToLower(k)] = strings.Join(v, ", ")
			}
			attr.Request = rq
		}
		if p, ok := peer.FromContext(ctx); ok {
			pr := &attribute_context.AttributeContext_Peer{}
			if p.Addr != nil {
				switch addr := p.Addr.(type) {
				case *net.IPAddr:
					pr.Ip = addr.IP.String()
				case *net.TCPAddr:
					pr.Ip = addr.IP.String()
					pr.Port = int64(addr.Port)
				case *net.UDPAddr:
					pr.Ip = addr.IP.String()
					pr.Port = int64(addr.Port)
				case *net.UnixAddr:
					pr.Ip = addr.Name
				}
			}
			attr.Origin = pr
			attr.Source = pr
		}
		if err := validateHandler.Validate(ctx, attr, req.(proto.Message)); err != nil {
			if errorHandler != nil {
				return nil, errorHandler(err)
			}
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return handler(ctx, req)
	}
}
