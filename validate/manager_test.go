package validate

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestManagerRegistry(t *testing.T) {
	m, err := NewManager(validate.File_testdata_validate_test_proto)
	if err != nil {
		t.Error(err)
	}
	mm, err := NewManager(validate.File_testdata_validate_test_proto)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(mm.BuildValidaters())
	tests := []struct {
		Name               string
		BuildRegistry      func() *managerRegistry
		Manager            *Manager
		Pattern            string
		Lib                cel.Library
		WantRegisterErr    bool
		WantLoadLibraryErr bool
	}{
		{
			Name: "Nil manager",
			BuildRegistry: func() *managerRegistry {
				r := &managerRegistry{registry: &sync.Map{}}
				return r
			},
			WantRegisterErr: true,
		},
		{
			Name: "Register twice",
			BuildRegistry: func() *managerRegistry {
				r := &managerRegistry{registry: &sync.Map{}}
				r.Register(m)
				return r
			},
			Manager:         m,
			WantRegisterErr: true,
		},
		{
			Name: "Load library on used manager",
			BuildRegistry: func() *managerRegistry {
				r := &managerRegistry{registry: &sync.Map{}}
				return r
			},
			Pattern:            "*",
			Lib:                &options.Library{},
			Manager:            mm,
			WantLoadLibraryErr: true,
		},
		{
			Name: "OK",
			BuildRegistry: func() *managerRegistry {
				r := &managerRegistry{registry: &sync.Map{}}
				return r
			},
			Manager: m,
		},
	}
	for _, tt := range tests {
		registry = &managerRegistry{registry: &sync.Map{}}
		t.Run(tt.Name, func(t *testing.T) {
			registry = tt.BuildRegistry()
			err := registry.Register(tt.Manager)
			if (tt.WantRegisterErr && err == nil) || (!tt.WantRegisterErr && err != nil) {
				t.Errorf("wantRegisterErr %v, got %v", tt.WantRegisterErr, err)
			} else if !tt.WantRegisterErr {
				err = registry.LoadLibrary(tt.Pattern, tt.Lib)
				if (tt.WantLoadLibraryErr && err == nil) || (!tt.WantLoadLibraryErr && err != nil) {
					t.Errorf("wantLoadLibraryErr %v, got %v", tt.WantLoadLibraryErr, err)
				}
			}
		})
	}
}

func TestManager(t *testing.T) {
	tests := []struct {
		Name               string
		File               protoreflect.FileDescriptor
		Lib                cel.Library
		Opts               []ManagerOption
		WantNewManagerErr  bool
		WantLoadLibraryErr bool
		WantBuildErr       bool
	}{
		{
			Name:              "Nil descriptor",
			WantNewManagerErr: true,
		},
		{
			Name:         "Missing const declaration",
			File:         validate.File_testdata_validate_manager_proto,
			WantBuildErr: true,
		},
		{
			Name: "OK (const declared in options)",
			File: validate.File_testdata_validate_manager_proto,
			Opts: []ManagerOption{WithOptions(&Options{
				Rule: &FileRule{
					Options: &options.Options{
						Globals: &options.Options_Globals{
							Constants: map[string]string{
								"name_const": "name",
							},
						},
					},
				},
			})},
			WantBuildErr: false,
		},
		{
			Name: "OK (const declared in lib)",
			File: validate.File_testdata_validate_manager_proto,
			Lib: &options.Library{
				EnvOpts: []cel.EnvOption{
					cel.Declarations(decls.NewConst(
						"name_const", decls.String, &expr.Constant{ConstantKind: &expr.Constant_StringValue{StringValue: "name"}},
					)),
				},
			},
			WantBuildErr: false,
		},
		{
			Name: "OK",
			File: validate.File_testdata_validate_test_proto,
		},
	}
	for _, tt := range tests {
		registry = &managerRegistry{registry: &sync.Map{}}
		t.Run(tt.Name, func(t *testing.T) {
			m, err := NewManager(tt.File, tt.Opts...)
			if (tt.WantNewManagerErr && err == nil) || (!tt.WantNewManagerErr && err != nil) {
				t.Errorf("wantNewManagerErr %v, got %v", tt.WantNewManagerErr, err)
			} else if !tt.WantNewManagerErr {
				err = m.LoadLibrary(tt.Lib)
				if (tt.WantLoadLibraryErr && err == nil) || (!tt.WantLoadLibraryErr && err != nil) {
					t.Errorf("wantLoadLibraryErr %v, got %v", tt.WantLoadLibraryErr, err)
				} else if !tt.WantLoadLibraryErr {
					err = m.BuildValidaters()
					if (tt.WantBuildErr && err == nil) || (!tt.WantBuildErr && err != nil) {
						t.Errorf("wantBuildErr %v, got %v", tt.WantBuildErr, err)
					}
				}
			}
		})
	}
}
