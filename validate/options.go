package validate

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/parser"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type Library struct {
	EnvOpts []cel.EnvOption
	PgrOpts []cel.ProgramOption
}

func (l *Library) CompileOptions() []cel.EnvOption     { return l.EnvOpts }
func (l *Library) ProgramOptions() []cel.ProgramOption { return l.PgrOpts }

func BuildEnvOption(options *Options, descs ...protoreflect.MessageDescriptor) cel.EnvOption {
	if options != nil {
		decls := []*v1alpha1.Decl{}
		if options.Globals != nil {
			decls = append(decls, buildDeclsFromOptionsGlobals(options.Globals)...)
		}
		if options.Overloads != nil {
			decls = append(decls, buildDeclsFromOptionsOverloads(options.Overloads)...)
		}
		reservedNames := map[string]bool{}
		if options.StdlibOverridingEnabled {
			if options.Globals != nil {
				for k := range options.Globals.Constants {
					reservedNames[k] = true
				}
				for k := range options.Globals.Functions {
					reservedNames[k] = true
				}
			}
			if options.Overloads != nil {
				for k := range options.Overloads.Functions {
					reservedNames[k] = true
				}
				for k := range options.Overloads.Variables {
					reservedNames[k] = true
				}
			}
			for _, desc := range descs {
				if desc != nil {
					for i := 0; i < desc.Fields().Len(); i++ {
						reservedNames[desc.Fields().Get(i).TextName()] = true
					}
				}
			}
		}
		for _, decl := range checker.StandardDeclarations() {
			if _, ok := reservedNames[decl.Name]; !ok {
				decls = append(decls, decl)
			}
		}
		return cel.Lib(&Library{
			EnvOpts: []cel.EnvOption{cel.Declarations(decls...), cel.Macros(cel.StandardMacros...)},
			PgrOpts: []cel.ProgramOption{cel.Functions(functions.StandardOverloads()...)},
		})
	}
	return cel.StdLib()
}

func buildDeclsFromOptionsGlobals(globals *Options_Globals) []*v1alpha1.Decl {
	dcls := []*v1alpha1.Decl{}
	if globals != nil {
		for k, v := range globals.Constants {
			dcls = append(dcls, decls.NewConst(
				k,
				decls.String,
				&v1alpha1.Constant{ConstantKind: &v1alpha1.Constant_StringValue{StringValue: v}},
			))
		}
	}
	return dcls
}

func buildDeclsFromOptionsOverloads(overloads *Options_Overloads) []*v1alpha1.Decl {
	dcls := []*v1alpha1.Decl{}
	if overloads != nil {
		for name, v := range overloads.Functions {
			args := []*v1alpha1.Type{}
			overload := name
			for i := 0; i < len(v.Args); i++ {
				args = append(args, TypeFromOverloadType(v.Args[i]))
			}
			dcls = append(dcls, decls.NewFunction(
				name, decls.NewOverload(
					overload,
					args,
					TypeFromOverloadType(v.Result),
				),
			))
		}
		for k, v := range overloads.Variables {
			dcls = append(dcls, decls.NewVar(k, TypeFromOverloadType(v)))
		}
	}
	return dcls
}

func TypeFromOverloadType(t *Options_Overloads_Type) *v1alpha1.Type {
	switch v := t.Type.(type) {
	case *Options_Overloads_Type_Primitive_:
		switch v.Primitive {
		case Options_Overloads_Type_BOOL:
			return decls.Bool
		case Options_Overloads_Type_INT:
			return decls.Int
		case Options_Overloads_Type_UINT:
			return decls.Uint
		case Options_Overloads_Type_DOUBLE:
			return decls.Double
		case Options_Overloads_Type_BYTES:
			return decls.Bytes
		case Options_Overloads_Type_STRING:
			return decls.String
		case Options_Overloads_Type_DURATION:
			return decls.Duration
		case Options_Overloads_Type_TIMESTAMP:
			return decls.Timestamp
		case Options_Overloads_Type_ERROR:
			return decls.Error
		case Options_Overloads_Type_DYN:
			return decls.Dyn
		case Options_Overloads_Type_ANY:
			return decls.Any
		}
	case *Options_Overloads_Type_Object:
		return decls.NewObjectType(v.Object)
	case *Options_Overloads_Type_Array_:
		return decls.NewListType(TypeFromOverloadType(v.Array.Type))
	case *Options_Overloads_Type_Map_:
		return decls.NewMapType(TypeFromOverloadType(v.Map.Key), TypeFromOverloadType(v.Map.Value))
	}
	return decls.Null
}

func BuildMacros(options *Options, expr string, envOpts []cel.EnvOption) ([]parser.Macro, error) {
	macros := []parser.Macro{}
	if rawMacros, err := findMacros(options, expr, envOpts); err != nil {
		return nil, fmt.Errorf("find macros error: %v", err)
	} else {
		env, err := cel.NewCustomEnv(envOpts...)
		if err != nil {
			return nil, fmt.Errorf("new env error: %w", err)
		}
		for _, macro := range rawMacros {
			ast, issues := env.Compile(options.Globals.Functions[macro])
			if issues != nil && issues.Err() != nil {
				return nil, fmt.Errorf("macro error: %w", issues.Err())
			}
			macros = append(macros, parser.NewGlobalMacro(macro, 0, buildMacroExpander(ast)))
		}
	}
	return macros, nil
}

func buildMacroExpander(ast *cel.Ast) parser.MacroExpander {
	return func(eh parser.ExprHelper, target *v1alpha1.Expr, args []*v1alpha1.Expr) (*v1alpha1.Expr, *common.Error) {
		return translateMacroExpr(ast.Expr(), eh), nil
	}
}

func translateMacroExpr(e *v1alpha1.Expr, eh parser.ExprHelper) *v1alpha1.Expr {
	if e == nil {
		return nil
	}
	switch exp := e.ExprKind.(type) {
	case *v1alpha1.Expr_ConstExpr:
		switch k := exp.ConstExpr.ConstantKind.(type) {
		case *v1alpha1.Constant_BoolValue:
			return eh.LiteralBool(k.BoolValue)
		case *v1alpha1.Constant_Int64Value:
			return eh.LiteralInt(k.Int64Value)
		case *v1alpha1.Constant_Uint64Value:
			return eh.LiteralUint(k.Uint64Value)
		case *v1alpha1.Constant_DoubleValue:
			return eh.LiteralDouble(k.DoubleValue)
		case *v1alpha1.Constant_StringValue:
			return eh.LiteralString(k.StringValue)
		case *v1alpha1.Constant_BytesValue:
			return eh.LiteralBytes(k.BytesValue)
		default:
			return e
		}
	case *v1alpha1.Expr_IdentExpr:
		return eh.Ident(exp.IdentExpr.GetName())
	case *v1alpha1.Expr_SelectExpr:
		return eh.Select(translateMacroExpr(exp.SelectExpr.GetOperand(), eh), exp.SelectExpr.GetField())
	case *v1alpha1.Expr_CallExpr:
		args := []*v1alpha1.Expr{}
		for i := 0; i < len(exp.CallExpr.Args); i++ {
			args = append(args, translateMacroExpr(exp.CallExpr.Args[i], eh))
		}
		if exp.CallExpr.Target != nil {
			return eh.ReceiverCall(exp.CallExpr.GetFunction(), translateMacroExpr(exp.CallExpr.Target, eh), args...)
		}
		return eh.GlobalCall(exp.CallExpr.GetFunction(), args...)
	case *v1alpha1.Expr_ListExpr:
		args := []*v1alpha1.Expr{}
		for i := 0; i < len(exp.ListExpr.GetElements()); i++ {
			args = append(args, translateMacroExpr(exp.ListExpr.Elements[i], eh))
		}
		return eh.NewList(args...)
	case *v1alpha1.Expr_StructExpr:
		fieldInits := []*v1alpha1.Expr_CreateStruct_Entry{}
		for i := 0; i < len(exp.StructExpr.Entries); i++ {
			entry := exp.StructExpr.Entries[i]
			switch eexp := entry.KeyKind.(type) {
			case *v1alpha1.Expr_CreateStruct_Entry_FieldKey:
				fieldInits = append(fieldInits, eh.NewObjectFieldInit(eexp.FieldKey, entry.Value, entry.OptionalEntry))
			case *v1alpha1.Expr_CreateStruct_Entry_MapKey:
				fieldInits = append(fieldInits, eh.NewMapEntry(eexp.MapKey, entry.Value, entry.OptionalEntry))
			}
		}
		return eh.NewObject(exp.StructExpr.MessageName, fieldInits...)
	case *v1alpha1.Expr_ComprehensionExpr:
		return eh.Fold(
			exp.ComprehensionExpr.IterVar,
			translateMacroExpr(exp.ComprehensionExpr.IterRange, eh),
			exp.ComprehensionExpr.AccuVar,
			translateMacroExpr(exp.ComprehensionExpr.AccuInit, eh),
			translateMacroExpr(exp.ComprehensionExpr.LoopCondition, eh),
			translateMacroExpr(exp.ComprehensionExpr.LoopStep, eh),
			translateMacroExpr(exp.ComprehensionExpr.Result, eh),
		)
	}
	return nil
}

func findMacros(options *Options, expr string, opts []cel.EnvOption) ([]string, error) {
	if options == nil || options.Globals == nil {
		return nil, nil
	}
	envOpts := opts
	for k := range options.Globals.Functions {
		envOpts = append(envOpts, cel.Declarations(decls.NewFunction(k, decls.NewOverload(k, []*v1alpha1.Type{}, &v1alpha1.Type{TypeKind: &v1alpha1.Type_Dyn{}}))))
	}
	env, err := cel.NewCustomEnv(envOpts...)
	if err != nil {
		return nil, fmt.Errorf("new env error: %w", err)
	}
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("compile error: %w", issues.Err())
	}
	return findMacrosExpr(ast.Expr(), options.Globals.Functions), nil
}

func findMacrosExpr(e *v1alpha1.Expr, m map[string]string) []string {
	res := []string{}
	switch exp := e.ExprKind.(type) {
	case *v1alpha1.Expr_ConstExpr:
	case *v1alpha1.Expr_IdentExpr:
	case *v1alpha1.Expr_SelectExpr:
	case *v1alpha1.Expr_CallExpr:
		if _, ok := m[exp.CallExpr.Function]; ok {
			res = append(res, exp.CallExpr.Function)
		} else {
			for _, i := range exp.CallExpr.Args {
				res = append(res, findMacrosExpr(i, m)...)
			}
		}
	case *v1alpha1.Expr_ListExpr:
		for _, i := range exp.ListExpr.Elements {
			res = append(res, findMacrosExpr(i, m)...)
		}
	case *v1alpha1.Expr_StructExpr:
	case *v1alpha1.Expr_ComprehensionExpr:
	}
	return res
}
