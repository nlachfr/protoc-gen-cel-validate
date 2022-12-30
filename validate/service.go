package validate

import (
	"github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func BuildServiceValidateProgram(config *ValidateOptions, desc protoreflect.ServiceDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (map[string]*Program, error) {
	if config == nil {
		config = &ValidateOptions{}
	}
	serviceRule := proto.GetExtension(desc.Options(), E_Service).(*ValidateRule)
	if serviceRule != nil {
		config.Options = options.Join(config.Options, serviceRule.Options)
	}
	m := map[string]*Program{}
	for i := 0; i < desc.Methods().Len(); i++ {
		methodDesc := desc.Methods().Get(i)
		if methodRule := proto.GetExtension(methodDesc.Options(), E_Method).(*ValidateRule); methodRule != nil {
			exprs := methodRule.Exprs
			if methodRule.Expr != "" {
				exprs = append([]string{methodRule.Expr}, exprs...)
			}
			if len(exprs) == 0 && serviceRule != nil {
				defaultExprs := serviceRule.Exprs
				if serviceRule.Expr != "" {
					defaultExprs = append([]string{serviceRule.Expr}, exprs...)
				}
				exprs = defaultExprs
			}
			if len(exprs) > 0 {
				if pgr, err := BuildMethodValidateProgram(exprs, config, methodDesc, envOpt, imports...); err != nil {
					return nil, err
				} else {
					m[string(methodDesc.FullName())] = pgr
				}
			}
		}
	}
	return m, nil
}
