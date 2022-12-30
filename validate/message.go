package validate

import (
	"fmt"

	"github.com/Neakxs/protocel/options"
	"github.com/google/cel-go/cel"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func BuildMessageValidateProgram(config *ValidateOptions, desc protoreflect.MessageDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (map[string]*Program, error) {
	if config == nil {
		config = &ValidateOptions{}
	}
	messageRule := proto.GetExtension(desc.Options(), E_Message).(*ValidateRule)
	if messageRule != nil {
		config.Options = options.Join(config.Options, messageRule.Options)
	}
	lib := &options.Library{}
	if envOpt != nil {
		lib.EnvOpts = append(lib.EnvOpts, envOpt)
	}
	lib.EnvOpts = append(lib.EnvOpts, cel.DeclareContextProto(desc))
	lib.EnvOpts = append(lib.EnvOpts, buildValidatersFunctions(desc)...)
	lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(config.Options, desc))
	envOpt = cel.Lib(lib)
	m := map[string]*Program{}
	resourceReferenceMap, err := GenerateResourceTypePatternMapping(desc.ParentFile(), imports...)
	if err != nil {
		return nil, err
	}
	for i := 0; i < desc.Fields().Len(); i++ {
		fieldDesc := desc.Fields().Get(i)
		exprs := []string{}
		if resourceReference := proto.GetExtension(fieldDesc.Options(), annotations.E_ResourceReference).(*annotations.ResourceReference); resourceReference != nil && !config.ResourceReferenceSupportDisabled {
			var regexp string
			if resourceReference.Type != "" {
				if resourceReference.Type != "*" {
					regexp = fmt.Sprintf("^%s$", resourceReferenceMap[resourceReference.Type])
				}
			} else if resourceReference.ChildType != "" {
				regexp = fmt.Sprintf("^%s", resourceReferenceMap[resourceReference.ChildType])
			}
			if regexp != "" {
				if fieldDesc.IsList() {
					exprs = append(exprs, fmt.Sprintf(`%s.all(s, s.matches("%s"))`, fieldDesc.TextName(), regexp))
				} else if fieldDesc.Kind() == protoreflect.StringKind {
					exprs = append(exprs, fmt.Sprintf(`%s.matches("%s")`, fieldDesc.TextName(), regexp))
				}
			}
		}
		if fieldRule := proto.GetExtension(fieldDesc.Options(), E_Field).(*ValidateRule); fieldRule != nil {
			exprs = append(exprs, fieldRule.Exprs...)
			if fieldRule.Expr != "" {
				exprs = append([]string{fieldRule.Expr}, exprs...)
			}
			if len(exprs) == 0 && messageRule != nil {
				defaultExprs := messageRule.Exprs
				if messageRule.Expr != "" {
					defaultExprs = append([]string{messageRule.Expr}, exprs...)
				}
				exprs = defaultExprs
			}
		}
		if len(exprs) > 0 {
			if pgr, err := BuildValidateProgram(exprs, config, envOpt, imports...); err != nil {
				return nil, err
			} else {
				m[string(fieldDesc.FullName())] = pgr
			}
		}
	}
	return m, nil
}