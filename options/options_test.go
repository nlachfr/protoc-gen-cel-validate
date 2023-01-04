package options

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestJoin(t *testing.T) {
	tests := []struct {
		Name    string
		Options []*Options
		Want    *Options
	}{
		{
			Name:    "nil option",
			Options: []*Options{},
			Want:    nil,
		},
		{
			Name: "single option",
			Options: []*Options{{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn": "func()",
					},
					Constants: map[string]string{
						"const": "const",
					},
				},
				StdlibOverridingEnabled: true,
			}},
			Want: &Options{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn": "func()",
					},
					Constants: map[string]string{
						"const": "const",
					},
				},
				StdlibOverridingEnabled: true,
			},
		},
		{
			Name: "merge with nil",
			Options: []*Options{nil, {
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn": "func()",
					},
				},
			}},
			Want: &Options{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn": "func()",
					},
				},
			},
		},
		{
			Name: "merge without conflict",
			Options: []*Options{{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn1": "func()",
					},
					Constants: map[string]string{
						"const1": "const",
					},
				},
			}, {
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn2": "func()",
					},
					Constants: map[string]string{
						"const2": "const",
					},
				},
			}},
			Want: &Options{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn1": "func()",
						"fn2": "func()",
					},
					Constants: map[string]string{
						"const1": "const",
						"const2": "const",
					},
				},
			},
		},
		{
			Name: "merge with conflict",
			Options: []*Options{{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn": "func()",
					},
				},
			}, {
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn": "func2()",
					},
				},
			}},
			Want: &Options{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"fn": "func2()",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			res := Join(tt.Options...)
			if !cmp.Equal(tt.Want, res, protocmp.Transform()) {
				t.Errorf("want %v, got %v", tt.Want, res)
			}
		})
	}
}
