package options

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

type testOpt struct {
	fns  []*FunctionOverload
	vars []*VariableOverload
}

func (o *testOpt) GetFunctionOverloads() []*FunctionOverload { return o.fns }
func (o *testOpt) GetVariableOverloads() []*VariableOverload { return o.vars }

func TestBuildRuntimeLibrary(t *testing.T) {
	tests := []struct {
		Name    string
		Rule    string
		Config  *Options
		Options RuntimeOptions
		WantErr bool
	}{
		{
			Name: "Missing function overload",
			Rule: `myVariable == myFunction("")`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"myFunction": {
							Args: []*Options_Overloads_Type{{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							}},
							Result: &Options_Overloads_Type{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*Options_Overloads_Type{
						"myVariable": {
							Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				vars: []*VariableOverload{{
					Name:  "myVariable",
					Value: "ok",
				}},
			},
			WantErr: true,
		},
		{
			Name: "Missing variable overload",
			Rule: `myVariable == myFunction("")`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"myFunction": {
							Args: []*Options_Overloads_Type{{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							}},
							Result: &Options_Overloads_Type{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*Options_Overloads_Type{
						"myVariable": {
							Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				fns: []*FunctionOverload{{
					Name: "myFunction",
					Function: func(v ...ref.Val) ref.Val {
						return types.String("ok")
					},
				}},
			},
			WantErr: true,
		},
		{
			Name: "OK (1 arg)",
			Rule: `myVariable == myFunction("")`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"myFunction": {
							Args: []*Options_Overloads_Type{{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							}},
							Result: &Options_Overloads_Type{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*Options_Overloads_Type{
						"myVariable": {
							Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				fns: []*FunctionOverload{{
					Name: "myFunction",
					Function: func(v ...ref.Val) ref.Val {
						return types.String("ok")
					},
				}},
				vars: []*VariableOverload{{
					Name:  "myVariable",
					Value: "ok",
				}},
			},
			WantErr: false,
		},
		{
			Name: "OK (2 args)",
			Rule: `myVariable == myFunction("", "")`,
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"myFunction": {
							Args: []*Options_Overloads_Type{{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							}, {
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							}},
							Result: &Options_Overloads_Type{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*Options_Overloads_Type{
						"myVariable": {
							Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				fns: []*FunctionOverload{{
					Name: "myFunction",
					Function: func(v ...ref.Val) ref.Val {
						return types.String("ok")
					},
				}},
				vars: []*VariableOverload{{
					Name:  "myVariable",
					Value: "ok",
				}},
			},
			WantErr: false,
		},
		{
			Name: "OK (any args)",
			Rule: "myVariable == myFunction()",
			Config: &Options{
				Overloads: &Options_Overloads{
					Functions: map[string]*Options_Overloads_Function{
						"myFunction": {
							Args: []*Options_Overloads_Type{},
							Result: &Options_Overloads_Type{
								Type: &Options_Overloads_Type_Primitive_{
									Primitive: Options_Overloads_Type_STRING,
								},
							},
						},
					},
					Variables: map[string]*Options_Overloads_Type{
						"myVariable": {
							Type: &Options_Overloads_Type_Primitive_{
								Primitive: Options_Overloads_Type_STRING,
							},
						},
					},
				},
			},
			Options: &testOpt{
				fns: []*FunctionOverload{{
					Name: "myFunction",
					Function: func(v ...ref.Val) ref.Val {
						return types.String("ok")
					},
				}},
				vars: []*VariableOverload{{
					Name:  "myVariable",
					Value: "ok",
				}},
			},
			WantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			envOpts := []cel.EnvOption{BuildEnvOption(tt.Config), cel.Lib(BuildRuntimeLibrary(tt.Config, tt.Options))}
			if macros, err := BuildMacros(tt.Config, tt.Rule, envOpts); err != nil {
				if !tt.WantErr {
					t.Errorf("wantErr %v, got %v", tt.WantErr, err)
				}
			} else {
				envOpts = append(envOpts, cel.Macros(macros...))
				if env, err := cel.NewCustomEnv(envOpts...); err != nil {
					if !tt.WantErr {
						t.Errorf("wantErr %v, got %v", tt.WantErr, err)
					}
				} else if ast, issues := env.Compile(tt.Rule); issues != nil && issues.Err() != nil {
					if !tt.WantErr {
						t.Errorf("wantErr %v, got %v", tt.WantErr, issues.Err())
					}
				} else if pgr, err := env.Program(ast); err != nil {
					if !tt.WantErr {
						t.Errorf("wantErr %v, got %v", tt.WantErr, err)
					}
				} else if _, _, err = pgr.Eval(map[string]interface{}{}); (err == nil && tt.WantErr) || (!tt.WantErr && err != nil) {
					t.Errorf("wantErr %v, got %v", tt.WantErr, err)
				}
			}
		})
	}
}
