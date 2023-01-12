package gateway

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bufbuild/connect-go"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
)

func BuildAttributeContext(s connect.Spec, p connect.Peer, h http.Header) *attribute_context.AttributeContext {
	peer := &attribute_context.AttributeContext_Peer{}
	ss := strings.SplitN(p.Addr, ":", 2)
	switch len(ss) {
	case 1:
		peer.Ip = ss[0]
	case 2:
		if port, err := strconv.ParseInt(ss[1], 10, 32); err == nil {
			peer.Ip = p.Addr
		} else {
			peer.Ip = ss[0]
			peer.Port = port
		}
	}
	attr := &attribute_context.AttributeContext{
		Api: &attribute_context.AttributeContext_Api{
			Operation: s.Procedure,
		},
		Origin:  peer,
		Source:  peer,
		Request: &attribute_context.AttributeContext_Request{},
	}
	switch p.Protocol {
	case connect.ProtocolConnect:
		attr.Api.Protocol = "http"
	case connect.ProtocolGRPC, connect.ProtocolGRPCWeb:
		attr.Api.Protocol = "grpc"
	}
	for k, v := range h {
		attr.Request.Headers[k] = strings.Join(v, ",")
	}
	return attr
}
