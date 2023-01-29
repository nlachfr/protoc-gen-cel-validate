package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/nlachfr/protocel/testdata/validate"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewUpstream(t *testing.T) {
	tests := []struct {
		Name    string
		Config  *Configuration_Server_Upstream
		Server  func() *httptest.Server
		WantErr bool
	}{
		{
			Name:    "Nil config",
			WantErr: true,
		},
		{
			Name: "Wrong address",
			Config: &Configuration_Server_Upstream{
				Address: "upstream",
			},
			WantErr: true,
		},
		{
			Name: "Server with wrong scheme",
			Config: &Configuration_Server_Upstream{
				Address: "http://upstream",
				Server:  "http://upstream",
			},
			WantErr: true,
		},
		{
			Name: "HTTP address",
			Config: &Configuration_Server_Upstream{
				Address: "http://upstream",
			},
		},
		{
			Name: "HTTPS address",
			Config: &Configuration_Server_Upstream{
				Address: "https://upstream",
			},
		},
		{
			Name: "Server with tcp scheme",
			Config: &Configuration_Server_Upstream{
				Address: "http://upstream",
				Server:  "tcp://upstream",
			},
		},
		{
			Name: "Server with unix scheme",
			Config: &Configuration_Server_Upstream{
				Address: "http://upstream",
				Server:  "unix://upstream",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := NewUpstream(tt.Config)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

func TestUpstream(t *testing.T) {
	mux := http.NewServeMux()
	md := validate.File_testdata_validate_service_proto.Services().ByName("Service").Methods().ByName("Rpc")
	rpc := fmt.Sprintf("/%s/%s", md.Parent().FullName(), md.Name())
	mux.Handle(rpc, connect.NewUnaryHandler(rpc, func(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
		return &connect.Response[emptypb.Empty]{}, nil
	}))
	request := dynamicpb.NewMessage(md.Input())
	tests := []struct {
		Name    string
		Args    func() (*Upstream, *httptest.Server)
		WantErr bool
	}{
		{
			Name: "HTTP Server (TCP)",
			Args: func() (*Upstream, *httptest.Server) {
				srv := httptest.NewUnstartedServer(h2c.NewHandler(mux, &http2.Server{}))
				srv.EnableHTTP2 = true
				srv.Start()
				upstream, err := NewUpstream(&Configuration_Server_Upstream{
					Address: fmt.Sprintf("http://%s", srv.Listener.Addr().String()),
				})
				if err != nil {
					t.Error(err)
				}
				return upstream, srv
			},
			WantErr: false,
		},
		{
			Name: "HTTP Server (Unix)",
			Args: func() (*Upstream, *httptest.Server) {
				srv := httptest.NewUnstartedServer(h2c.NewHandler(mux, &http2.Server{}))
				srv.EnableHTTP2 = true
				l, err := net.Listen("unix", filepath.Join(t.TempDir(), "socket.unix"))
				if err != nil {
					t.Error(err)
				}
				srv.Listener = l
				srv.Start()
				upstream, err := NewUpstream(&Configuration_Server_Upstream{
					Address: "http://localhost",
					Server:  fmt.Sprintf("unix://%s", srv.Listener.Addr().String()),
				})
				if err != nil {
					t.Error(err)
				}
				return upstream, srv
			},
			WantErr: false,
		},
		{
			Name: "HTTPS Server (TCP)",
			Args: func() (*Upstream, *httptest.Server) {
				srv := httptest.NewUnstartedServer(h2c.NewHandler(mux, &http2.Server{}))
				srv.EnableHTTP2 = true
				srv.StartTLS()
				upstream, err := NewUpstream(&Configuration_Server_Upstream{
					Address: fmt.Sprintf("https://%s", srv.Listener.Addr().String()),
				})
				if err != nil {
					t.Error(err)
				}
				upstream.httpClient.Transport.(*http2.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
				return upstream, srv
			},
			WantErr: false,
		},
		{
			Name: "HTTPS Server (Unix)",
			Args: func() (*Upstream, *httptest.Server) {
				srv := httptest.NewUnstartedServer(h2c.NewHandler(mux, &http2.Server{}))
				srv.EnableHTTP2 = true
				l, err := net.Listen("unix", filepath.Join(t.TempDir(), "socket.unix"))
				if err != nil {
					t.Error(err)
				}
				srv.Listener = l
				srv.StartTLS()
				upstream, err := NewUpstream(&Configuration_Server_Upstream{
					Address: "https://localhost",
					Server:  fmt.Sprintf("unix://%s", srv.Listener.Addr().String()),
				})
				if err != nil {
					t.Error(err)
				}
				upstream.httpClient.Transport.(*http2.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
				return upstream, srv
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			upstream, srv := tt.Args()
			_, err := upstream.NewClient(md).CallUnary(context.Background(), connect.NewRequest(&request))
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
			srv.Close()
		})
	}
}
