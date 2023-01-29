package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	tvalidate "github.com/nlachfr/protocel/testdata/validate"
	"github.com/nlachfr/protocel/validate"
)

func TestFindServiceUpstream(t *testing.T) {
	sd := tvalidate.File_testdata_validate_service_proto.Services().ByName("Service")
	var up1, up2, up3 *Upstream = &Upstream{}, &Upstream{}, &Upstream{}
	tests := []struct {
		Name         string
		Upstreams    map[string]*Upstream
		WantUpstream *Upstream
	}{
		{
			Name: "Wildcard pattern",
			Upstreams: map[string]*Upstream{
				"*": up1,
			},
			WantUpstream: up1,
		},
		{
			Name: "Pattern not matching",
			Upstreams: map[string]*Upstream{
				"not.matching": up2,
			},
			WantUpstream: nil,
		},
		{
			Name: "Wildcard and subwildcard patterns",
			Upstreams: map[string]*Upstream{
				"*":          up1,
				"testdata.*": up2,
			},
			WantUpstream: up2,
		},
		{
			Name: "Subwildcard and subsubwildcard patterns",
			Upstreams: map[string]*Upstream{
				"testdata.validate.*": up1,
				"testdata.*":          up2,
			},
			WantUpstream: up1,
		},
		{
			Name: "Local wildcard pattern and exact pattern",
			Upstreams: map[string]*Upstream{
				"testdata.validate.*":       up1,
				"testdata.validate.Service": up2,
				"testdata.validate.Ser*":    up3,
			},
			WantUpstream: up2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			up, err := findServiceUpstream(sd, tt.Upstreams)
			if err != nil {
				t.Error(err)
			}
			if tt.WantUpstream != up {
				t.Errorf("wantUpstream %v, got %v", tt.WantUpstream, up)
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	linker, err := NewLinker(context.Background(), &Configuration_Files{
		Sources: []string{"testdata/validate/service.proto"},
		Imports: []string{".."},
	})
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		Name          string
		Linker        *Linker
		Config        *Configuration_Server
		Options       *validate.Options
		WantNewErr    bool
		WantListenErr bool
	}{
		{
			Name:       "Nil server config",
			WantNewErr: true,
		},
		{
			Name: "Nil upstream",
			Config: &Configuration_Server{
				Upstreams: map[string]*Configuration_Server_Upstream{
					"*": nil,
				},
			},
			WantNewErr: true,
		},
		{
			Name:          "Nil linker",
			Config:        &Configuration_Server{},
			WantListenErr: true,
		},
		{
			Name:   "Error in validation rule",
			Linker: linker,
			Config: &Configuration_Server{
				Listen: []string{"127.0.0.1:0"},
				Upstreams: map[string]*Configuration_Server_Upstream{
					"testdata.validate.ServiceOptions": {Address: "http://localhost"},
				},
			},
			WantNewErr: true,
		},
		{
			Name:   "OK (tcp)",
			Linker: linker,
			Config: &Configuration_Server{
				Listen: []string{"127.0.0.1:0"},
				Upstreams: map[string]*Configuration_Server_Upstream{
					"testdata.validate.Service": {Address: "http://localhost"},
				},
			},
		},
		{
			Name:   "OK (unix)",
			Linker: linker,
			Config: &Configuration_Server{
				Listen: []string{fmt.Sprintf("unix://%s/unix.sock", t.TempDir())},
				Upstreams: map[string]*Configuration_Server_Upstream{
					"testdata.validate.Service": {Address: "http://localhost"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			server, err := NewServer(tt.Linker, tt.Config, tt.Options)
			if (tt.WantNewErr && err == nil) || (!tt.WantNewErr && err != nil) {
				t.Errorf("wantNewErr %v, got %v", tt.WantNewErr, err)
			} else if err == nil {
				errChan := make(chan error, 1)
				go func() {
					if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
						errChan <- err
					}
					errChan <- nil
				}()
				time.Sleep(time.Second / 100)
				server.Shutdown(context.TODO())
				select {
				case cerr := <-errChan:
					err = cerr
				default:
				}
				if (tt.WantListenErr && err == nil) || (!tt.WantListenErr && err != nil) {
					t.Errorf("wantListenErr %v, got %v", tt.WantListenErr, err)
				}
			}
		})
	}
}
