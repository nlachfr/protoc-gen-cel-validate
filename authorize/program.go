package authorize

import (
	"fmt"
	"net/http"

	"github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func BuildAuthzProgramFromDesc(expr string, imports []protoreflect.FileDescriptor, msgDesc protoreflect.MessageDescriptor, config *AuthorizeOptions, libs ...cel.Library) (cel.Program, error) {
	envOpts := []cel.EnvOption{
		cel.TypeDescs(msgDesc.Parent()),
	}
	for i := 0; i < len(imports); i++ {
		envOpts = append(envOpts, cel.TypeDescs(imports[i]))
	}
	for i := 0; i < len(libs); i++ {
		envOpts = append(envOpts, cel.Lib(libs[i]))
	}
	return buildAuthzProgram(expr, msgDesc, config, envOpts...)
}

func BuildAuthzProgram(expr string, msg proto.Message, config *AuthorizeOptions, libs ...cel.Library) (cel.Program, error) {
	envOpts := []cel.EnvOption{
		cel.Types(msg),
	}
	for i := 0; i < len(libs); i++ {
		envOpts = append(envOpts, cel.Lib(libs[i]))
	}
	return buildAuthzProgram(expr, msg.ProtoReflect().Descriptor(), config, envOpts...)
}

func buildAuthzProgram(expr string, desc protoreflect.MessageDescriptor, config *AuthorizeOptions, envOpts ...cel.EnvOption) (cel.Program, error) {
	envOpts = append(envOpts,
		cel.Declarations(
			decls.NewVar(
				"headers",
				decls.NewMapType(
					decls.String,
					decls.NewListType(decls.String),
				),
			),
			decls.NewVar(
				"request",
				decls.NewObjectType(string(desc.FullName())),
			),
		),
		cel.Function("get",
			cel.MemberOverload(
				"get",
				[]*cel.Type{cel.MapType(cel.StringType, cel.ListType(cel.StringType)), cel.StringType},
				cel.StringType,
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					var h http.Header
					switch m := lhs.Value().(type) {
					case map[string][]string:
						h = http.Header(m)
					case http.Header:
						h = m
					default:
						return types.String("")
					}
					if s, ok := rhs.Value().(string); ok {
						return types.String(h.Get(s))
					}
					return types.String("")
				}),
			),
		),
		cel.Function("values",
			cel.MemberOverload(
				"values",
				[]*cel.Type{cel.MapType(cel.StringType, cel.ListType(cel.StringType)), cel.StringType},
				cel.ListType(cel.StringType),
				cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					var h http.Header
					switch m := lhs.Value().(type) {
					case map[string][]string:
						h = http.Header(m)
					case http.Header:
						h = m
					default:
						return types.NewStringList(nil, []string{})
					}
					if s, ok := rhs.Value().(string); ok {
						return types.NewStringList(TypeAdapterFunc(func(value interface{}) ref.Val { return types.String(value.(string)) }), h.Values(s))
					}
					return types.String("")
				}),
			),
		),
	)
	if config != nil {
		if config.Options != nil {
			envOpts = append(envOpts, options.BuildEnvOption(config.Options))
			if macros, err := options.BuildMacros(config.Options, expr, envOpts); err != nil {
				return nil, fmt.Errorf("build macros error: %v", err)
			} else {
				envOpts = append(envOpts, cel.Macros(macros...))
			}
		}
		envOpts = append(envOpts, options.BuildStdLib(config.Options))
	} else {
		envOpts = append(envOpts, options.BuildStdLib(nil))
	}
	env, err := cel.NewCustomEnv(envOpts...)
	if err != nil {
		return nil, fmt.Errorf("new env error: %w", err)
	}
	ast, issues := env.Compile(expr)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("compile error: %w", issues.Err())
	}
	if !ast.OutputType().IsAssignableType(cel.BoolType) {
		return nil, fmt.Errorf("output type not bool")
	}
	pgr, err := env.Program(ast, cel.OptimizeRegex(interpreter.MatchesRegexOptimization))
	if err != nil {
		return nil, fmt.Errorf("program error: %w", err)
	}
	return pgr, nil
}
