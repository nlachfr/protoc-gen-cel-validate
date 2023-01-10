package validate

import (
	"context"
	"fmt"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/validate/errors"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ServiceValidateProgram interface {
	Validate(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error
}

type serviceValidateProgram struct {
	methodsDesc     map[string]protoreflect.MethodDescriptor
	methodsPrograms map[string]ValidateProgram
}

func (vp *serviceValidateProgram) Validate(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error {
	if attr == nil || attr.Api == nil {
		return nil
	} else if pgr, ok := vp.methodsPrograms[attr.Api.Operation]; ok {
		req := map[string]interface{}{
			"attribute_context": attr,
			"request":           m,
		}
		for _, p := range pgr.CEL() {
			if val, _, err := p.ContextEval(ctx, req); err != nil {
				return errors.Wrap(err, m, vp.methodsDesc[attr.Api.Operation], attr)
			} else if !types.IsBool(val) || !val.Value().(bool) {
				return errors.New(m, vp.methodsDesc[attr.Api.Operation], attr)
			}
		}
	}
	return nil
}

func BuildServiceValidateProgram(config *ValidateOptions, desc protoreflect.ServiceDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (ServiceValidateProgram, error) {
	if config == nil {
		config = &ValidateOptions{}
	}
	var serviceRule *ValidateRule
	if config.Rules != nil {
		serviceRule = config.Rules[string(desc.FullName())]
	}
	if r := proto.GetExtension(desc.Options(), E_Service).(*ValidateRule); r != nil {
		serviceRule = r
	}
	if serviceRule != nil {
		config.Options = options.Join(config.Options, serviceRule.Options)
	}
	descs := map[string]protoreflect.MethodDescriptor{}
	m := map[string]ValidateProgram{}
	for i := 0; i < desc.Methods().Len(); i++ {
		methodDesc := desc.Methods().Get(i)
		exprs := []string{}
		var methodRule *ValidateRule
		if config.Rules != nil {
			methodRule = config.Rules[string(methodDesc.FullName())]
		}
		if r := proto.GetExtension(methodDesc.Options(), E_Method).(*ValidateRule); r != nil {
			methodRule = r
		}
		if methodRule != nil {
			exprs = methodRule.Exprs
			if methodRule.Expr != "" {
				exprs = append([]string{methodRule.Expr}, exprs...)
			}
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
				m[fmt.Sprintf("/%s/%s", string(desc.FullName()), string(methodDesc.Name()))] = pgr
			}
		}
		descs[fmt.Sprintf("/%s/%s", string(desc.FullName()), string(methodDesc.Name()))] = methodDesc
	}
	return &serviceValidateProgram{methodsDesc: descs, methodsPrograms: m}, nil
}

func BuildMethodValidateProgram(exprs []string, config *ValidateOptions, desc protoreflect.MethodDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (ValidateProgram, error) {
	if config == nil {
		config = &ValidateOptions{}
	}
	lib := &options.Library{EnvOpts: []cel.EnvOption{
		cel.TypeDescs(attribute_context.File_google_rpc_context_attribute_context_proto),
		cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
		cel.TypeDescs(desc.Input().ParentFile()),
		cel.Variable("request", cel.ObjectType(string(desc.Input().FullName()))),
	}}
	if envOpt != nil {
		lib.EnvOpts = append(lib.EnvOpts, envOpt)
	}
	if r := proto.GetExtension(desc.Options(), E_Method).(*ValidateRule); r != nil {
		config.Options = options.Join(config.Options, r.Options)
	}
	lib.EnvOpts = append(lib.EnvOpts, buildValidatersFunctions(desc.Input())...)
	lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(config.Options))
	return BuildValidateProgram(exprs, config, cel.Lib(lib), imports...)
}
