package gateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/nlachfr/protocel/validate"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type connectClient interface {
	CallUnary(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.Response[*dynamicpb.Message], error)
	CallClientStream(ctx context.Context) *connect.ClientStreamForClient[*dynamicpb.Message, *dynamicpb.Message]
	CallServerStream(ctx context.Context, request *connect.Request[*dynamicpb.Message]) (*connect.ServerStreamForClient[*dynamicpb.Message], error)
	CallBidiStream(ctx context.Context) *connect.BidiStreamForClient[*dynamicpb.Message, *dynamicpb.Message]
}

type methodHandler struct {
	srv    validate.ServiceRuleValidater
	client connectClient
}

func (h *methodHandler) unary(ctx context.Context, r *connect.Request[*dynamicpb.Message]) (*connect.Response[*dynamicpb.Message], error) {
	if err := h.srv.Validate(ctx, BuildAttributeContext(r.Spec(), r.Peer(), r.Header()), *r.Msg); err != nil {
		return nil, err
	} else if res, err := h.client.CallUnary(ctx, r); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (h *methodHandler) clientStream(ctx context.Context, cs *connect.ClientStream[*dynamicpb.Message]) (*connect.Response[dynamicpb.Message], error) {
	csfc := h.client.CallClientStream(ctx)
	defer csfc.CloseAndReceive()
	attr := BuildAttributeContext(cs.Spec(), cs.Peer(), cs.RequestHeader())
	for cs.Receive() {
		msg := cs.Msg()
		if err := h.srv.Validate(ctx, attr, *msg); err != nil {
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
}

func (h *methodHandler) serverStream(ctx context.Context, r *connect.Request[*dynamicpb.Message], ss *connect.ServerStream[dynamicpb.Message]) error {
	if err := h.srv.Validate(ctx, BuildAttributeContext(r.Spec(), r.Peer(), r.Header()), *r.Msg); err != nil {
		return err
	} else if res, err := h.client.CallServerStream(ctx, r); err != nil {
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
}

func (h *methodHandler) bidiStream(ctx context.Context, bs *connect.BidiStream[*dynamicpb.Message, dynamicpb.Message]) error {
	return nil
}

func NewServiceHandler(sd protoreflect.ServiceDescriptor, srv validate.ServiceRuleValidater, upstream *Upstream, opts ...connect.HandlerOption) (string, http.Handler) {
	sroot := fmt.Sprintf("/%s/", sd.FullName())
	mux := http.NewServeMux()
	for i := 0; i < sd.Methods().Len(); i++ {
		md := sd.Methods().Get(i)
		mroot := sroot + string(md.Name())
		handler := &methodHandler{
			srv:    srv,
			client: upstream.NewClient(md),
		}
		opt := connect.WithHandlerOptions(append(opts, newCodecs(md.Input()))...)
		if md.IsStreamingClient() {
			if md.IsStreamingServer() {
				mux.Handle(mroot, connect.NewBidiStreamHandler(mroot, handler.bidiStream, opt))
			} else {
				mux.Handle(mroot, connect.NewClientStreamHandler(mroot, handler.clientStream, opt))
			}
		} else if md.IsStreamingServer() {
			mux.Handle(mroot, connect.NewServerStreamHandler(mroot, handler.serverStream, opt))
		} else {
			mux.Handle(mroot, connect.NewUnaryHandler(mroot, handler.unary, opt))
		}
	}
	return sroot, mux
}
