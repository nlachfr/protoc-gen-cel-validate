package authorize

import (
	"testing"

	options "github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/authorize"
)

func TestBuildAuthzProgramFromDesc(t *testing.T) {
	tests := []struct {
		Name    string
		Expr    string
		Config  *AuthorizeOptions
		WantErr bool
	}{
		{
			Name:    "Unknown field",
			Expr:    `request.pong`,
			Config:  nil,
			WantErr: true,
		},
		{
			Name:    "Invalid return type",
			Expr:    `request.ping`,
			Config:  nil,
			WantErr: true,
		},
		{
			Name:    "OK",
			Expr:    `request.ping == "ping"`,
			Config:  nil,
			WantErr: false,
		},
		{
			Name:    "OK (get metadata)",
			Expr:    `headers.get("x-user") == ""`,
			WantErr: false,
		},
		{
			Name:    "OK (values metadata)",
			Expr:    `size(headers.values("x-user")) == 0`,
			WantErr: false,
		},
		{
			Name: "OK (with constant)",
			Expr: "request.ping == constPing",
			Config: &AuthorizeOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"constPing": "ping",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name: "OK (with bool macro)",
			Expr: `rule()`,
			Config: &AuthorizeOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Functions: map[string]string{
							"rule": `request.ping == "ping"`,
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name: "OK (with str macro)",
			Expr: `rule() == "ping"`,
			Config: &AuthorizeOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Functions: map[string]string{
							"rule": `request.ping`,
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name: "OK (array with str macro)",
			Expr: `"ping" in [rule()]`,
			Config: &AuthorizeOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Functions: map[string]string{
							"rule": `request.ping`,
						},
					},
				},
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := BuildAuthzProgramFromDesc(tt.Expr, nil, authorize.File_testdata_authorize_test_proto.Messages().Get(0), tt.Config)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

func TestBuildAuthProgram(t *testing.T) {
	tests := []struct {
		Name    string
		Expr    string
		Config  *AuthorizeOptions
		WantErr bool
	}{
		{
			Name:    "Unknown field",
			Expr:    `request.pong`,
			Config:  nil,
			WantErr: true,
		},
		{
			Name:    "Invalid return type",
			Expr:    `request.ping`,
			Config:  nil,
			WantErr: true,
		},
		{
			Name:    "OK",
			Expr:    `request.ping == "ping"`,
			Config:  nil,
			WantErr: false,
		},
		{
			Name:    "OK (get metadata)",
			Expr:    `headers.get("x-user") == ""`,
			WantErr: false,
		},
		{
			Name: "OK (with constant)",
			Expr: "request.ping == constPing",
			Config: &AuthorizeOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Constants: map[string]string{
							"constPing": "ping",
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name: "OK (with bool macro)",
			Expr: `rule()`,
			Config: &AuthorizeOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Functions: map[string]string{
							"rule": `request.ping == "ping"`,
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name: "OK (with str macro)",
			Expr: `rule() == "ping"`,
			Config: &AuthorizeOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Functions: map[string]string{
							"rule": `request.ping`,
						},
					},
				},
			},
			WantErr: false,
		},
		{
			Name: "OK (array with str macro)",
			Expr: `"ping" in [rule()]`,
			Config: &AuthorizeOptions{
				Options: &options.Options{
					Globals: &options.Options_Globals{
						Functions: map[string]string{
							"rule": `request.ping`,
						},
					},
				},
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := BuildAuthzProgram(tt.Expr, &authorize.PingRequest{}, tt.Config)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
