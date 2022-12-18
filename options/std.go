package options

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/interpreter/functions"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

func BuildStdLib(options *Options, descs ...protoreflect.MessageDescriptor) cel.EnvOption {
	if options != nil && options.StdlibOverridingEnabled {
		reservedNames := map[string]bool{}
		if options.Globals != nil {
			for k, _ := range options.Globals.Constants {
				reservedNames[k] = true
			}
			for k, _ := range options.Globals.Functions {
				reservedNames[k] = true
			}
		}
		if options.Overloads != nil {
			for k, _ := range options.Overloads.Functions {
				reservedNames[k] = true
			}
			for k, _ := range options.Overloads.Variables {
				reservedNames[k] = true
			}
		}
		for _, desc := range descs {
			for i := 0; i < desc.ReservedNames().Len(); i++ {
				reservedNames[string(desc.ReservedNames().Get(i))] = true
			}
		}
		decls := []*expr.Decl{}
		for _, decl := range checker.StandardDeclarations() {
			if _, ok := reservedNames[decl.Name]; !ok {
				decls = append(decls, decl)
			}
		}
		return cel.Lib(&library{
			EnvOpts: []cel.EnvOption{cel.Declarations(decls...)},
			PgrOpts: []cel.ProgramOption{cel.Functions(functions.StandardOverloads()...)},
		})
	} else {
		return cel.StdLib()
	}
}
