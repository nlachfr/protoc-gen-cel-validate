package authorize

import (
	"context"
	"net/http"
	"testing"

	"github.com/Neakxs/protocel/testdata/authorize"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

func TestAuthzInterceptor(t *testing.T) {
	env, _ := cel.NewEnv(
		cel.Types(&authorize.PingRequest{}),
		cel.Declarations(
			decls.NewVar(
				"headers",
				decls.NewMapType(
					decls.String,
					decls.NewListType(decls.String),
				),
			),
			decls.NewVar(
				"request",
				decls.NewObjectType(string((&authorize.PingRequest{}).ProtoReflect().Descriptor().FullName())),
			),
		),
	)
	astBool, _ := env.Compile(`request.ping == "ping" && "hdr" in headers`)
	pgrBool, _ := env.Program(astBool)
	astString, _ := env.Compile(`request.ping`)
	pgrString, _ := env.Program(astString)
	tests := []struct {
		Name    string
		Mapping map[string]cel.Program
		Request *authorize.PingRequest
		WantErr bool
	}{
		{
			Name: "Permission denied (bool)",
			Mapping: map[string]cel.Program{
				"": pgrBool,
			},
			Request: &authorize.PingRequest{Ping: ""},
			WantErr: true,
		},
		{
			Name: "OK (bool)",
			Mapping: map[string]cel.Program{
				"": pgrBool,
			},
			Request: &authorize.PingRequest{Ping: "ping"},
			WantErr: false,
		},
		{
			Name: "Unknown (str)",
			Mapping: map[string]cel.Program{
				"": pgrString,
			},
			Request: &authorize.PingRequest{Ping: "ping"},
			WantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := NewAuthzInterceptor(tt.Mapping).Authorize(context.Background(), "", http.Header{"hdr": []string{}}, tt.Request)
			if (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
