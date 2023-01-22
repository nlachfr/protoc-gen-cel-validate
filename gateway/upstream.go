package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func NewUpstream(cfg *Configuration_Server_Upstream) (*Upstream, error) {
	if cfg == nil || cfg.Address == "" {
		return nil, fmt.Errorf("")
	}
	addrUrl, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, err
	}
	var defaultPort bool
	var srvnet, srvaddr string
	if cfg.Server != "" {
		if srvUrl, err := url.Parse(cfg.Server); err != nil {
			return nil, err
		} else {
			switch srvUrl.Scheme {
			case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
				srvnet = srvUrl.Scheme
				srvaddr = srvUrl.Host
				if srvUrl.Port() == "" {
					defaultPort = true
				}
			case "unix", "unixgram", "unixpacket":
				srvnet = srvUrl.Scheme
				srvaddr = srvUrl.Host + srvUrl.RequestURI()
			default:
				return nil, fmt.Errorf(`unknown scheme: "%s"`, srvUrl.Scheme)
			}
		}
	}
	dialer := &net.Dialer{
		Timeout: time.Second * 2,
	}
	var httpTransport *http2.Transport
	switch addrUrl.Scheme {
	case "http":
		if srvnet != "" {
			if defaultPort {
				srvaddr = srvaddr + ":80"
			}
			httpTransport = &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					return dialer.DialContext(ctx, srvnet, srvaddr)
				},
			}
		} else {
			httpTransport = &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					return dialer.DialContext(ctx, network, addr)
				},
			}
		}
	case "https":
		if srvnet != "" {
			if defaultPort {
				srvaddr = srvaddr + ":443"
			}
			httpTransport = &http2.Transport{
				DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					return tls.DialWithDialer(dialer, srvnet, srvaddr, cfg)
				},
			}
		} else {
			httpTransport = &http2.Transport{
				DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					return tls.DialWithDialer(dialer, network, addr, cfg)
				},
			}
		}
	default:
		return nil, fmt.Errorf(`unknown scheme: "%s"`, addrUrl.Scheme)
	}
	var opt connect.ClientOption
	switch cfg.Protocol {
	case Configuration_Server_Upstream_GRPC:
		opt = connect.WithGRPC()
	case Configuration_Server_Upstream_GRPC_WEB:
		opt = connect.WithGRPCWeb()
	case Configuration_Server_Upstream_CONNECT:
		// opt = connect.WithProtoJSON()
	}
	return &Upstream{
		target: addrUrl,
		httpClient: &http.Client{
			Transport: httpTransport,
		},
		opt: opt,
	}, nil
}

type Upstream struct {
	target     *url.URL
	httpClient *http.Client
	opt        connect.ClientOption
}

func (u *Upstream) NewClient(md protoreflect.MethodDescriptor) *connect.Client[*dynamicpb.Message, *dynamicpb.Message] {
	opts := []connect.ClientOption{newCodecs(md.Output())}
	if u.opt != nil {
		opts = append(opts, u.opt)
	}
	return connect.NewClient[*dynamicpb.Message, *dynamicpb.Message](
		u.httpClient,
		u.target.JoinPath(fmt.Sprintf("/%s/%s", md.Parent().FullName(), md.Name())).String(),
		opts...,
	)
}
