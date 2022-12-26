package options

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/interpreter/functions"
	v1alpha1 "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

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
				for i := 0; i < desc.Fields().Len(); i++ {
					reservedNames[desc.Fields().Get(i).TextName()] = true
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
