package gateway

import (
	"context"
	"path/filepath"
	"testing"
)

func TestNewLinker(t *testing.T) {
	validatePath, _ := filepath.Abs("../testdata/validate/")
	tests := []struct {
		Name    string
		Config  *Configuration_Files
		WantErr bool
	}{
		{
			Name: "Nil config",
		},
		{
			Name: "Source path does not match",
			Config: &Configuration_Files{
				Sources: []string{"my.path.proto"},
			},
			WantErr: true,
		},
		{
			Name: "Import path pattern does not match",
			Config: &Configuration_Files{
				Imports: []string{"my.path.proto"},
			},
			WantErr: true,
		},
		{
			Name: "Missing import",
			Config: &Configuration_Files{
				Sources: []string{
					"test.proto",
				},
				Imports: []string{
					validatePath,
				},
			},
			WantErr: true,
		},
		{
			Name: "OK",
			Config: &Configuration_Files{
				Sources: []string{"message.proto"},
				Imports: []string{
					"../testdata/validate",
					"..",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := NewLinker(context.Background(), tt.Config)
			if (tt.WantErr && err == nil) || (!tt.WantErr && err != nil) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}
