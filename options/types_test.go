package options

import (
	"testing"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/go-cmp/cmp"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestTypeFromOverloadType(t *testing.T) {
	tests := []struct {
		Name string
		In   *Options_Overloads_Type
		Out  *v1alpha1.Type
	}{
		{
			Name: "Primitive bool",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_BOOL,
				},
			},
			Out: decls.Bool,
		},
		{
			Name: "Primitive int",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_INT,
				},
			},
			Out: decls.Int,
		},
		{
			Name: "Primitive uint",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_UINT,
				},
			},
			Out: decls.Uint,
		},
		{
			Name: "Primitive double",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_DOUBLE,
				},
			},
			Out: decls.Double,
		},
		{
			Name: "Primitive bytes",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_BYTES,
				},
			},
			Out: decls.Bytes,
		},
		{
			Name: "Primitive string",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_STRING,
				},
			},
			Out: decls.String,
		},
		{
			Name: "Primitive duration",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_DURATION,
				},
			},
			Out: decls.Duration,
		},
		{
			Name: "Primitive timestamp",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_TIMESTAMP,
				},
			},
			Out: decls.Timestamp,
		},
		{
			Name: "Primitive error",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_ERROR,
				},
			},
			Out: decls.Error,
		},
		{
			Name: "Primitive dyn",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_DYN,
				},
			},
			Out: decls.Dyn,
		},
		{
			Name: "Primitive any",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Primitive_{
					Primitive: Options_Overloads_Type_ANY,
				},
			},
			Out: decls.Any,
		},
		{
			Name: "Object",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Object{
					Object: "object",
				},
			},
			Out: decls.NewObjectType("object"),
		},
		{
			Name: "Array",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Array_{
					Array: &Options_Overloads_Type_Array{
						Type: &Options_Overloads_Type{
							Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_BOOL,
							},
						},
					},
				},
			},
			Out: decls.NewListType(decls.Bool),
		},
		{
			Name: "Map",
			In: &Options_Overloads_Type{
				Type: &Options_Overloads_Type_Map_{
					Map: &Options_Overloads_Type_Map{
						Key: &Options_Overloads_Type{
							Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							},
						},
						Value: &Options_Overloads_Type{
							Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Out: decls.NewMapType(decls.String, decls.String),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			res := TypeFromOverloadType(tt.In)
			if !cmp.Equal(res, tt.Out, protocmp.Transform()) {
				t.Errorf("want %v, got %v", tt.Out, res)
			}
		})
	}
}
