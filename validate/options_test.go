package validate

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestBuildEnvOption(t *testing.T) {
	tests := []struct {
		Name    string
		Expr    string
		Desc    protoreflect.MessageDescriptor
		EnvOpt  cel.EnvOption
		Config  *Options
		WantErr bool
	}{
		{
			Name:    "No options",
			Expr:    "1 == 2",
			WantErr: false,
		},
		{
			Name:    "Constant (error)",
			Expr:    `name == "name"`,
			WantErr: true,
		},
		{
			Name: "Constant",
			Expr: `name == "name"`,
			Config: &Options{
				Globals: &Options_Globals{
					Constants: map[string]string{"name": "name"},
				},
			},
			WantErr: false,
		},
		{
			Name: "Overload variable (error)",
			Expr: `name == "name"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Variables: map[string]*Options_Overloads_Type{"name": {
						Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}},
					},
				},
			},
			WantErr: true,
		},
		{
			Name: "Overload variable",
			Expr: `name == "name"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Variables: map[string]*Options_Overloads_Type{"name": {
						Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}},
					},
				},
			},
			EnvOpt: cel.Lib(&Library{
				PgrOpts: []cel.ProgramOption{cel.Globals(map[string]interface{}{"name": "name"})},
			}),
			WantErr: false,
		},
		{
			Name: "Overload function (error)",
			Expr: `getName() == "name"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"getName()": {Result: &Options_Overloads_Type{Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}}},
					},
				},
			},
			WantErr: true,
		},
		{
			Name: "Overload function",
			Expr: `getName() == "name"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"getName": {Result: &Options_Overloads_Type{Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}}},
					},
				},
			},
			EnvOpt: cel.Lib(&Library{
				EnvOpts: []cel.EnvOption{cel.Function("getName", cel.Overload("getName", []*cel.Type{}, cel.StringType, cel.FunctionBinding(func(values ...ref.Val) ref.Val { return types.String("name") })))},
				PgrOpts: []cel.ProgramOption{},
			}),
			WantErr: false,
		},
		{
			Name:    "Stdlib override (error)",
			Expr:    `type == "type"`,
			WantErr: true,
		},
		{
			Name: "Stdlib override (const)",
			Expr: `type == "type"`,
			Config: &Options{
				Globals: &Options_Globals{
					Constants: map[string]string{"type": "name"},
				},
				StdlibOverridingEnabled: true,
			},
			WantErr: false,
		},
		{
			Name: "Stdlib override (variable)",
			Expr: `type == "type"`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Variables: map[string]*Options_Overloads_Type{"type": {
						Type: &Options_Overloads_Type_Primitive_{Primitive: Options_Overloads_Type_STRING}},
					},
				},
				StdlibOverridingEnabled: true,
			},
			EnvOpt: cel.Lib(&Library{
				PgrOpts: []cel.ProgramOption{cel.Globals(map[string]interface{}{"type": "name"})},
			}),
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			opts := []cel.EnvOption{BuildEnvOption(tt.Config, tt.Desc)}
			if tt.EnvOpt != nil {
				opts = append(opts, tt.EnvOpt)
			}
			env, err := cel.NewCustomEnv(opts...)
			if err == nil {
				ast, issues := env.Compile(tt.Expr)
				if issues != nil {
					err = issues.Err()
				} else {
					pgr, perr := env.Program(ast, cel.EvalOptions(cel.OptOptimize))
					if perr != nil {
						err = perr
					} else {
						_, _, err = pgr.Eval(map[string]interface{}{})
					}
				}
			}
			if (err != nil && !tt.WantErr) || (err == nil && tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

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

func TestBuildMacros(t *testing.T) {
	tests := []struct {
		Name    string
		Rule    string
		Options *Options
		WantErr bool
	}{
		{
			Name:    "None (no options)",
			Rule:    "1 == 1",
			WantErr: false,
		},
		{
			Name:    "None (err undefined function)",
			Rule:    "1 == err()",
			WantErr: true,
			Options: &Options{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"macro": "1 == 1",
					},
				},
			},
		},
		{
			Name:    "One",
			Rule:    "true == macro()",
			WantErr: false,
			Options: &Options{
				Globals: &Options_Globals{
					Functions: map[string]string{
						"macro": "1 == 1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := BuildMacros(tt.Options, tt.Rule, []cel.EnvOption{BuildEnvOption(tt.Options)})
			if (err == nil && tt.WantErr) || (err != nil && !tt.WantErr) {
				t.Errorf("wantErr %v, got %v", tt.WantErr, err)
			}
		})
	}
}

func TestFindMacrosExpr(t *testing.T) {
	m := map[string]string{
		"myFn": "1 == 1",
	}
	tests := []struct {
		Name   string
		Expr   *v1alpha1.Expr
		Result []string
	}{
		{
			Name: "None",
			Expr: &v1alpha1.Expr{
				ExprKind: &v1alpha1.Expr_CallExpr{
					CallExpr: &v1alpha1.Expr_Call{
						Function: operators.Equals,
						Args: []*v1alpha1.Expr{
							{
								ExprKind: &v1alpha1.Expr_ConstExpr{
									ConstExpr: &v1alpha1.Constant{
										ConstantKind: &v1alpha1.Constant_Int64Value{
											Int64Value: 1,
										},
									},
								},
							},
							{
								ExprKind: &v1alpha1.Expr_ConstExpr{
									ConstExpr: &v1alpha1.Constant{
										ConstantKind: &v1alpha1.Constant_Int64Value{
											Int64Value: 1,
										},
									},
								},
							},
						},
					},
				},
			},
			Result: []string{},
		},
		{
			Name: "Simple",
			Expr: &v1alpha1.Expr{
				ExprKind: &v1alpha1.Expr_CallExpr{
					CallExpr: &v1alpha1.Expr_Call{
						Function: "myFn",
						Args:     []*v1alpha1.Expr{},
					},
				},
			},
			Result: []string{"myFn"},
		},
		{
			Name: "Embed",
			Expr: &v1alpha1.Expr{
				ExprKind: &v1alpha1.Expr_CallExpr{
					CallExpr: &v1alpha1.Expr_Call{
						Function: operators.In,
						Args: []*v1alpha1.Expr{
							{
								ExprKind: &v1alpha1.Expr_ConstExpr{
									ConstExpr: &v1alpha1.Constant{
										ConstantKind: &v1alpha1.Constant_BoolValue{
											BoolValue: true,
										},
									},
								},
							},
							{
								ExprKind: &v1alpha1.Expr_CallExpr{
									CallExpr: &v1alpha1.Expr_Call{
										Function: "myFn",
										Args:     []*v1alpha1.Expr{},
									},
								},
							},
						},
					},
				},
			},
			Result: []string{"myFn"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			res := findMacrosExpr(tt.Expr, m)
			if !cmp.Equal(res, tt.Result, cmpopts.SortSlices(func(l, r string) bool {
				return l < r
			})) {
				t.Errorf("want %v, got %v", tt.Result, res)
			}
		})
	}
}

func compareExpr(l, r *v1alpha1.Expr) bool {
	if l == nil && r == nil {
		return true
	}
	switch ll := l.ExprKind.(type) {
	case *v1alpha1.Expr_ConstExpr:
		rrr := r.GetConstExpr()
		switch lll := ll.ConstExpr.ConstantKind.(type) {
		case *v1alpha1.Constant_BoolValue:
			return lll.BoolValue == rrr.GetBoolValue()
		case *v1alpha1.Constant_Int64Value:
			return lll.Int64Value == rrr.GetInt64Value()
		case *v1alpha1.Constant_Uint64Value:
			return lll.Uint64Value == rrr.GetUint64Value()
		case *v1alpha1.Constant_DoubleValue:
			return lll.DoubleValue == rrr.GetDoubleValue()
		case *v1alpha1.Constant_StringValue:
			return lll.StringValue == rrr.GetStringValue()
		case *v1alpha1.Constant_BytesValue:
			return cmp.Equal(lll.BytesValue, rrr.GetBytesValue())
		}
	case *v1alpha1.Expr_IdentExpr:
		return ll.IdentExpr.Name == r.GetIdentExpr().Name
	case *v1alpha1.Expr_SelectExpr:
		rr := r.GetSelectExpr()
		return ll.SelectExpr.TestOnly == rr.TestOnly && compareExpr(ll.SelectExpr.Operand, rr.Operand) && ll.SelectExpr.Field == rr.Field
	case *v1alpha1.Expr_CallExpr:
		rr := r.GetCallExpr()
		for i := 0; i < len(ll.CallExpr.Args); i++ {
			if !compareExpr(ll.CallExpr.Args[i], rr.Args[i]) {
				return false
			}
		}
		return ll.CallExpr.Function == rr.Function && compareExpr(ll.CallExpr.Target, rr.Target)
	case *v1alpha1.Expr_ListExpr:
		rr := r.GetListExpr()
		for i := 0; i < len(ll.ListExpr.Elements); i++ {
			if !compareExpr(ll.ListExpr.Elements[i], rr.Elements[i]) {
				return false
			}
		}
		return true
	case *v1alpha1.Expr_StructExpr:
		rr := r.GetStructExpr()
		for i := 0; i < len(ll.StructExpr.Entries); i++ {
			lll := ll.StructExpr.Entries[i]
			rrr := rr.Entries[i]
			if lll.Id != rrr.Id || !compareExpr(lll.Value, rrr.Value) {
				return false
			}
			switch llll := lll.KeyKind.(type) {
			case *v1alpha1.Expr_CreateStruct_Entry_FieldKey:
				rrrr := rrr.GetFieldKey()
				if llll.FieldKey != rrrr {
					return false
				}
			case *v1alpha1.Expr_CreateStruct_Entry_MapKey:
				rrrr := rrr.GetMapKey()
				if !compareExpr(llll.MapKey, rrrr) {
					return false
				}
			}
		}
		return true
	case *v1alpha1.Expr_ComprehensionExpr:
		rr := r.GetComprehensionExpr()
		if compareExpr(ll.ComprehensionExpr.AccuInit, rr.AccuInit) && ll.ComprehensionExpr.AccuVar == rr.AccuVar && compareExpr(ll.ComprehensionExpr.IterRange, rr.IterRange) && ll.ComprehensionExpr.IterVar == rr.IterVar && compareExpr(ll.ComprehensionExpr.LoopCondition, rr.LoopCondition) && compareExpr(ll.ComprehensionExpr.LoopStep, rr.LoopStep) && compareExpr(ll.ComprehensionExpr.Result, rr.Result) {
			return true
		}
	}
	return false
}

func TestTranslateMacroExpr(t *testing.T) {
	tests := []struct {
		Name string
		Expr *v1alpha1.Expr
	}{
		{
			Name: "Int64 constant",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
				ConstExpr: &v1alpha1.Constant{
					ConstantKind: &v1alpha1.Constant_Int64Value{
						Int64Value: 5,
					},
				},
			}},
		},
		{
			Name: "Bool constant",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
				ConstExpr: &v1alpha1.Constant{
					ConstantKind: &v1alpha1.Constant_BoolValue{
						BoolValue: true,
					},
				},
			}},
		},
		{
			Name: "Double constant",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
				ConstExpr: &v1alpha1.Constant{
					ConstantKind: &v1alpha1.Constant_DoubleValue{
						DoubleValue: 6,
					},
				},
			}},
		},
		{
			Name: "Uint64 constant",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
				ConstExpr: &v1alpha1.Constant{
					ConstantKind: &v1alpha1.Constant_Uint64Value{
						Uint64Value: 7,
					},
				},
			}},
		},
		{
			Name: "String constant",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
				ConstExpr: &v1alpha1.Constant{
					ConstantKind: &v1alpha1.Constant_StringValue{
						StringValue: "true",
					},
				},
			}},
		},
		{
			Name: "Bytes constant",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
				ConstExpr: &v1alpha1.Constant{
					ConstantKind: &v1alpha1.Constant_BytesValue{
						BytesValue: []byte{0, 2, 4},
					},
				},
			}},
		},
		{
			Name: "Ident",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_IdentExpr{
				IdentExpr: &v1alpha1.Expr_Ident{
					Name: "ident",
				},
			}},
		},
		{
			Name: "Select",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_SelectExpr{
				SelectExpr: &v1alpha1.Expr_Select{
					Field: "field",
				},
			}},
		},
		{
			Name: "Call",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_CallExpr{
				CallExpr: &v1alpha1.Expr_Call{
					Function: "function",
				},
			}},
		},
		{
			Name: "List",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ListExpr{
				ListExpr: &v1alpha1.Expr_CreateList{
					Elements: []*v1alpha1.Expr{{ExprKind: &v1alpha1.Expr_ConstExpr{
						ConstExpr: &v1alpha1.Constant{
							ConstantKind: &v1alpha1.Constant_BoolValue{
								BoolValue: true,
							},
						},
					}}},
				},
			}},
		},
		{
			Name: "Struct field",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_StructExpr{
				StructExpr: &v1alpha1.Expr_CreateStruct{
					MessageName: "message",
					Entries: []*v1alpha1.Expr_CreateStruct_Entry{{
						KeyKind: &v1alpha1.Expr_CreateStruct_Entry_FieldKey{
							FieldKey: "key",
						},
					}},
				},
			}},
		},
		{
			Name: "Struct map",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_StructExpr{
				StructExpr: &v1alpha1.Expr_CreateStruct{
					MessageName: "message",
					Entries: []*v1alpha1.Expr_CreateStruct_Entry{{
						KeyKind: &v1alpha1.Expr_CreateStruct_Entry_MapKey{
							MapKey: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
								ConstExpr: &v1alpha1.Constant{
									ConstantKind: &v1alpha1.Constant_BoolValue{
										BoolValue: true,
									},
								},
							}},
						},
					}},
				},
			}},
		},
		{
			Name: "Comprehension",
			Expr: &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ComprehensionExpr{
				ComprehensionExpr: &v1alpha1.Expr_Comprehension{
					IterVar: "iter",
					AccuVar: "accu",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			res := translateMacroExpr(tt.Expr, &exprHelper{})
			if !cmp.Equal(tt.Expr, res, cmp.Comparer(compareExpr)) {
				t.Errorf("want %v, got %v", tt.Expr, res)
			}
		})
	}
}

type exprHelper struct{}

func (*exprHelper) Copy(e *v1alpha1.Expr) *v1alpha1.Expr {
	return e
}

func (*exprHelper) LiteralBool(value bool) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
		ConstExpr: &v1alpha1.Constant{
			ConstantKind: &v1alpha1.Constant_BoolValue{
				BoolValue: value,
			},
		},
	}}
}

func (*exprHelper) LiteralBytes(value []byte) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
		ConstExpr: &v1alpha1.Constant{
			ConstantKind: &v1alpha1.Constant_BytesValue{
				BytesValue: value,
			},
		},
	}}
}
func (*exprHelper) LiteralDouble(value float64) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
		ConstExpr: &v1alpha1.Constant{
			ConstantKind: &v1alpha1.Constant_DoubleValue{
				DoubleValue: value,
			},
		},
	}}
}

func (*exprHelper) LiteralInt(value int64) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
		ConstExpr: &v1alpha1.Constant{
			ConstantKind: &v1alpha1.Constant_Int64Value{
				Int64Value: value,
			},
		},
	}}
}

func (*exprHelper) LiteralString(value string) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
		ConstExpr: &v1alpha1.Constant{
			ConstantKind: &v1alpha1.Constant_StringValue{
				StringValue: value,
			},
		},
	}}
}

func (*exprHelper) LiteralUint(value uint64) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ConstExpr{
		ConstExpr: &v1alpha1.Constant{
			ConstantKind: &v1alpha1.Constant_Uint64Value{
				Uint64Value: value,
			},
		},
	}}
}

func (*exprHelper) NewList(elems ...*v1alpha1.Expr) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ListExpr{
		ListExpr: &v1alpha1.Expr_CreateList{
			Elements: elems,
		},
	}}
}

func (*exprHelper) NewMap(entries ...*v1alpha1.Expr_CreateStruct_Entry) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_StructExpr{
		StructExpr: &v1alpha1.Expr_CreateStruct{
			Entries: entries,
		},
	}}
}

func (*exprHelper) NewMapEntry(key *v1alpha1.Expr, val *v1alpha1.Expr, optional bool) *v1alpha1.Expr_CreateStruct_Entry {
	return &v1alpha1.Expr_CreateStruct_Entry{
		KeyKind: &v1alpha1.Expr_CreateStruct_Entry_MapKey{
			MapKey: key,
		},
		Value:         val,
		OptionalEntry: optional,
	}
}

func (*exprHelper) NewObject(typeName string, fieldInits ...*v1alpha1.Expr_CreateStruct_Entry) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_StructExpr{
		StructExpr: &v1alpha1.Expr_CreateStruct{
			MessageName: typeName,
			Entries:     fieldInits,
		},
	}}
}

func (*exprHelper) NewObjectFieldInit(field string, init *v1alpha1.Expr, optional bool) *v1alpha1.Expr_CreateStruct_Entry {
	return &v1alpha1.Expr_CreateStruct_Entry{
		KeyKind: &v1alpha1.Expr_CreateStruct_Entry_FieldKey{
			FieldKey: field,
		},
		Value:         init,
		OptionalEntry: optional,
	}
}

func (*exprHelper) Fold(iterVar string, iterRange *v1alpha1.Expr, accuVar string, accuInit *v1alpha1.Expr, condition *v1alpha1.Expr, step *v1alpha1.Expr, result *v1alpha1.Expr) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_ComprehensionExpr{
		ComprehensionExpr: &v1alpha1.Expr_Comprehension{
			IterVar:       iterVar,
			IterRange:     iterRange,
			AccuVar:       accuVar,
			AccuInit:      accuInit,
			LoopCondition: condition,
			LoopStep:      step,
			Result:        result,
		},
	}}
}

func (*exprHelper) Ident(name string) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_IdentExpr{
		IdentExpr: &v1alpha1.Expr_Ident{
			Name: name,
		},
	}}
}

func (*exprHelper) AccuIdent() *v1alpha1.Expr {
	return nil
}

func (*exprHelper) GlobalCall(function string, args ...*v1alpha1.Expr) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_CallExpr{
		CallExpr: &v1alpha1.Expr_Call{
			Function: function,
			Args:     args,
		},
	}}
}

func (*exprHelper) ReceiverCall(function string, target *v1alpha1.Expr, args ...*v1alpha1.Expr) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_CallExpr{
		CallExpr: &v1alpha1.Expr_Call{
			Target:   target,
			Function: function,
			Args:     args,
		},
	}}
}

func (*exprHelper) PresenceTest(operand *v1alpha1.Expr, field string) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_SelectExpr{
		SelectExpr: &v1alpha1.Expr_Select{
			Operand:  operand,
			Field:    field,
			TestOnly: true,
		},
	}}
}

func (*exprHelper) Select(operand *v1alpha1.Expr, field string) *v1alpha1.Expr {
	return &v1alpha1.Expr{ExprKind: &v1alpha1.Expr_SelectExpr{
		SelectExpr: &v1alpha1.Expr_Select{
			Operand: operand,
			Field:   field,
		},
	}}
}

func (*exprHelper) OffsetLocation(exprID int64) common.Location { return nil }
