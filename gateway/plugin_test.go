package gateway

import (
	"testing"
)

func TestLoadPlugins(t *testing.T) {
	tests := []struct {
		Name    string
		Config  *Configuration_Plugin
		WantErr bool
	}{
		{
			Name: "Invalid path",
			Config: &Configuration_Plugin{
				Path: "",
			},
			WantErr: true,
		},
		{
			Name: "File not ELF",
			Config: &Configuration_Plugin{
				Path: "plugin_test.go",
			},
			WantErr: true,
		},
		{
			Name: "Missing symbol",
			Config: &Configuration_Plugin{
				Path: "../testdata/gateway/plugin/dummy.so",
			},
			WantErr: true,
		},
		{
			Name: "Invalid symbol",
			Config: &Configuration_Plugin{
				Path: "../testdata/gateway/plugin/invalid_symbol.so",
			},
			WantErr: true,
		},
		{
			Name: "Invalid plugin",
			Config: &Configuration_Plugin{
				Path: "../testdata/gateway/plugin/invalid_plugin.so",
			},
			WantErr: true,
		},
		{
			Name: "OK",
			Config: &Configuration_Plugin{
				Path: "../testdata/gateway/plugin/valid.so",
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := LoadPlugins(tt.Config)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
