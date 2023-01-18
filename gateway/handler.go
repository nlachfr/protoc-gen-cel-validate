package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Neakxs/protocel/validate"
	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type MethodHandlerBuilder func(rpc string, srv validate.ServiceRuleValidater, client *connect.Client[dynamicpb.Message, *dynamicpb.Message], opts ...connect.HandlerOption) (string, http.Handler)

func NewServiceHandler(target string, httpClient connect.HTTPClient, srv validate.ServiceRuleValidater, desc protoreflect.ServiceDescriptor, opts ...connect.HandlerOption) (string, http.Handler) {
	pattern := fmt.Sprintf("/%s/", desc.FullName())
	mux := http.NewServeMux()
	for i := 0; i < desc.Methods().Len(); i++ {
		methodDesc := desc.Methods().Get(i)
		var handlerBuilder MethodHandlerBuilder
		if methodDesc.IsStreamingClient() {
			if methodDesc.IsStreamingServer() {
				handlerBuilder = NewBidiStreamMethodHandler
			} else {
				handlerBuilder = NewClientStreamMethodHandler
			}
		} else if methodDesc.IsStreamingServer() {
			handlerBuilder = NewServerStreamMethodHandler
		} else {
			handlerBuilder = NewUnaryMethodHandler
		}
		mux.Handle(handlerBuilder(
			pattern+string(methodDesc.Name()),
			srv,
			connect.NewClient[dynamicpb.Message, *dynamicpb.Message](httpClient, target, buildClientCodecs(methodDesc.Output())...),
			append(opts, buildHandlerCodecs(methodDesc.Input())...)...,
		))
	}
	return pattern, mux
}

func NewUnaryMethodHandler(rpc string, srv validate.ServiceRuleValidater, client *connect.Client[dynamicpb.Message, *dynamicpb.Message], opts ...connect.HandlerOption) (string, http.Handler) {
	return rpc, connect.NewUnaryHandler(rpc, func(ctx context.Context, r *connect.Request[*dynamicpb.Message]) (*connect.Response[dynamicpb.Message], error) {
		if err := srv.Validate(ctx, BuildAttributeContext(r.Spec(), r.Peer(), r.Header()), *r.Msg); err != nil {
			fmt.Printf("%t\n", err)
			return nil, err
		} else if res, err := client.CallUnary(ctx, connect.NewRequest(*r.Msg)); err != nil {
			return nil, err
		} else {
			return connect.NewResponse(*res.Msg), nil
		}
	}, opts...)
}

func NewClientStreamMethodHandler(rpc string, srv validate.ServiceRuleValidater, client *connect.Client[dynamicpb.Message, *dynamicpb.Message], opts ...connect.HandlerOption) (string, http.Handler) {
	return rpc, connect.NewClientStreamHandler(rpc, func(ctx context.Context, cs *connect.ClientStream[*dynamicpb.Message]) (*connect.Response[dynamicpb.Message], error) {
		csfc := client.CallClientStream(ctx)
		defer csfc.CloseAndReceive()
		attr := BuildAttributeContext(cs.Spec(), cs.Peer(), cs.RequestHeader())
		for cs.Receive() {
			msg := *cs.Msg()
			if err := srv.Validate(ctx, attr, msg); err != nil {
				return nil, err
			} else if err = csfc.Send(msg); err != nil {
				return nil, err
			}
		}
		if err := cs.Err(); err != nil {
			return nil, err
		} else if res, err := csfc.CloseAndReceive(); err != nil {
			return nil, err
		} else {
			return connect.NewResponse(*res.Msg), nil
		}
	})
}

func NewServerStreamMethodHandler(rpc string, srv validate.ServiceRuleValidater, client *connect.Client[dynamicpb.Message, *dynamicpb.Message], opts ...connect.HandlerOption) (string, http.Handler) {
	return rpc, connect.NewServerStreamHandler(rpc, func(ctx context.Context, r *connect.Request[*dynamicpb.Message], ss *connect.ServerStream[dynamicpb.Message]) error {
		if err := srv.Validate(ctx, BuildAttributeContext(r.Spec(), r.Peer(), r.Header()), *r.Msg); err != nil {
			return err
		} else if res, err := client.CallServerStream(ctx, connect.NewRequest(*r.Msg)); err != nil {
			return err
		} else {
			defer res.Close()
			for k, v := range res.ResponseHeader() {
				for _, vv := range v {
					ss.ResponseHeader().Add(k, vv)
				}
			}
			for res.Receive() {
				if err := ss.Send(*res.Msg()); err != nil {
					return err
				}
			}
			if err := res.Err(); err != nil {
				if errors.Is(err, io.EOF) {
					for k, v := range res.ResponseTrailer() {
						for _, vv := range v {
							ss.ResponseTrailer().Add(k, vv)
						}
					}
					return nil
				}
				return err
			}
		}
		return nil
	}, opts...)
}

func NewBidiStreamMethodHandler(rpc string, srv validate.ServiceRuleValidater, client *connect.Client[dynamicpb.Message, *dynamicpb.Message], opts ...connect.HandlerOption) (string, http.Handler) {
	return rpc, connect.NewBidiStreamHandler(rpc, func(ctx context.Context, bs *connect.BidiStream[*dynamicpb.Message, dynamicpb.Message]) error {
		return nil
	}, opts...)
}
