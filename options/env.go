package options

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func BuildEnvOption(options *Options) cel.EnvOption {
	decls := []*v1alpha1.Decl{}
	if options != nil {
		if options.Globals != nil {
			decls = append(decls, buildDeclsFromOptionsGlobals(options.Globals)...)
		}
		if options.Overloads != nil {
			decls = append(decls, buildDeclsFromOptionsOverloads(options.Overloads)...)
		}
	}
	return cel.Declarations(decls...)
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
