package validate

import (
	"context"
	"fmt"

	"github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func NewValidateInterceptor(methodProgramMapping map[string]cel.Program) ValidateInterceptor {
	return &validateInterceptor{
		methodProgramMapping: methodProgramMapping,
	}
}

type ValidateInterceptor interface {
	Validate(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error
}

type validateInterceptor struct {
	methodProgramMapping map[string]cel.Program
}

func (i *validateInterceptor) Validate(ctx context.Context, attr *attribute_context.AttributeContext, m proto.Message) error {
	if attr == nil || attr.Request == nil {
		return nil
	} else if pgr, ok := i.methodProgramMapping[attr.Request.Method]; ok {
		if m == nil {
			return fmt.Errorf(`nil message`)
		}
		req := map[string]interface{}{
			"attribute_context": attr,
		}
		fields := m.ProtoReflect().Descriptor().Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.Get(i)
			req[f.TextName()] = m.ProtoReflect().Get(f)
		}
		if val, _, err := pgr.ContextEval(ctx, req); err != nil {
			return err
		} else if !types.IsBool(val) || !val.Value().(bool) {
			return fmt.Errorf(`validation failed on "%s`, attr.Request.Path)
		}
	}
	return nil
}

func BuildMethodValidateProgram(expr string, config *ValidateOptions, desc protoreflect.MessageDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (cel.Program, error) {
	lib := &options.Library{EnvOpts: []cel.EnvOption{
		cel.TypeDescs(attribute_context.File_google_rpc_context_attribute_context_proto),
		cel.Variable("attribute_context", cel.ObjectType(string((&attribute_context.AttributeContext{}).ProtoReflect().Descriptor().FullName()))),
	}}
	if envOpt != nil {
		lib.EnvOpts = append(lib.EnvOpts, envOpt)
	}
	return BuildValidateProgram(expr, config, desc, cel.Lib(lib), imports...)
}
