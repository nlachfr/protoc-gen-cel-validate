package validate

import (
	"context"
	"fmt"
	"strings"

	"github.com/Neakxs/protocel/options"
	"github.com/Neakxs/protocel/validate/errors"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type FieldValidateProgram interface {
	Program() ValidateProgram
	IsRequired() bool
}

type fieldValidateProgram struct {
	program  ValidateProgram
	required bool
}

func (p *fieldValidateProgram) Program() ValidateProgram { return p.program }
func (p *fieldValidateProgram) IsRequired() bool         { return p.required }

type MessageValidateProgram interface {
	ValidateWithMask(ctx context.Context, m proto.Message, fm *fieldmaskpb.FieldMask) error
	FieldPrograms() map[string]FieldValidateProgram
}

type messageValidateProgram struct {
	fieldsPrograms map[string]FieldValidateProgram
}

func (p *messageValidateProgram) FieldPrograms() map[string]FieldValidateProgram {
	return p.fieldsPrograms
}

func (p *messageValidateProgram) ValidateWithMask(ctx context.Context, m proto.Message, fm *fieldmaskpb.FieldMask) error {
	if p.fieldsPrograms == nil {
		return fmt.Errorf("validation failed")
	}
	pathsMap := map[string][]string{}
	mdesc := m.ProtoReflect().Descriptor()
	if fm == nil {
		m.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
			pathsMap[fd.TextName()] = []string{""}
			return true
		})
	} else if len(fm.Paths) == 1 && fm.Paths[0] == "*" {
		for i := 0; i < mdesc.Fields().Len(); i++ {
			fdesc := mdesc.Fields().Get(i)
			pathsMap[fdesc.TextName()] = []string{""}
		}
	} else {
		for i := 0; i < len(fm.Paths); i++ {
			parts := strings.SplitN(fm.Paths[i], ".", 2)
			k := parts[0]
			v, ok := pathsMap[k]
			if !ok {
				pathsMap[k] = []string{}
			}
			if len(parts) > 1 {
				pathsMap[k] = append(v, parts[1:]...)
			} else {
				pathsMap[k] = append(v, "")
			}
		}
	}
	vars := map[string]interface{}{}
	for i := 0; i < m.ProtoReflect().Descriptor().Fields().Len(); i++ {
		field := m.ProtoReflect().Descriptor().Fields().Get(i)
		vars[field.TextName()] = m.ProtoReflect().Get(field)
	}
	for i := 0; i < mdesc.Fields().Len(); i++ {
		fdesc := mdesc.Fields().Get(i)
		if paths, ok := pathsMap[fdesc.TextName()]; ok {
			subs := []string{}
			for j := 0; j < len(paths); j++ {
				if paths[j] == "" {
					if pgr, ok := p.fieldsPrograms[fdesc.TextName()]; ok {
						if !isDefaultValue(m, fdesc) {
							if pgr.Program() != nil {
								for _, p := range pgr.Program().CEL() {
									if val, _, err := p.ContextEval(ctx, vars); err != nil {
										return errors.Wrap(err, m, fdesc, nil)
									} else if !types.IsBool(val) || !val.Value().(bool) {
										return errors.New(m, fdesc, nil)
									}
								}
							}
						} else if pgr.IsRequired() {
							for _, behavior := range proto.GetExtension(fdesc.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior) {
								if behavior == annotations.FieldBehavior_REQUIRED {
									return errors.New(m, fdesc, nil)
								}
							}
						}
					}
				} else if paths[j] != "*" {
					subs = append(subs, paths[j])
				}
			}
			if len(subs) > 0 && fdesc.Kind() == protoreflect.MessageKind {
				if v, ok := m.ProtoReflect().Get(fdesc).Message().Interface().(Validater); ok {
					if err := v.ValidateWithMask(ctx, &fieldmaskpb.FieldMask{Paths: subs}); err != nil {
						return errors.Wrap(err, m, fdesc, nil)
					}
				}
			}
		}
	}
	return nil
}

func BuildMessageValidateProgram(config *ValidateOptions, desc protoreflect.MessageDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (MessageValidateProgram, error) {
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
	envOpt = cel.Lib(lib)
	m := map[string]FieldValidateProgram{}
	for i := 0; i < desc.Fields().Len(); i++ {
		fieldDesc := desc.Fields().Get(i)
		if pgr, err := BuildFieldValidateProgram(config, fieldDesc, envOpt, imports...); err != nil {
			return nil, err
		} else {
			m[fieldDesc.TextName()] = pgr
		}
	}
	return &messageValidateProgram{fieldsPrograms: m}, nil
}

func BuildFieldValidateProgram(config *ValidateOptions, desc protoreflect.FieldDescriptor, envOpt cel.EnvOption, imports ...protoreflect.FileDescriptor) (FieldValidateProgram, error) {
	if config == nil {
		config = &ValidateOptions{}
	}
	fieldRule := proto.GetExtension(desc.Options(), E_Field).(*ValidateRule)
	if fieldRule != nil {
		config.Options = options.Join(config.Options, fieldRule.Options)
	}
	lib := &options.Library{}
	if envOpt != nil {
		lib.EnvOpts = append(lib.EnvOpts, envOpt)
	}
	lib.EnvOpts = append(lib.EnvOpts, options.BuildEnvOption(config.Options, desc.ContainingMessage()))
	envOpt = cel.Lib(lib)
	resourceReferenceMap := GenerateResourceTypePatternMapping(imports...)
	exprs := []string{}
	if fieldRule != nil {
		exprs = append(exprs, fieldRule.Exprs...)
		if fieldRule.Expr != "" {
			exprs = append([]string{fieldRule.Expr}, exprs...)
		}
	}
	if resourceReference := proto.GetExtension(desc.Options(), annotations.E_ResourceReference).(*annotations.ResourceReference); resourceReference != nil && !config.ResourceReferenceSupportDisabled {
		var ref string
		if resourceReference.Type != "" {
			if resourceReference.ChildType != "" {
				return nil, fmt.Errorf(`resource reference error: type and child_type are defined`)
			} else if resourceReference.Type != "*" {
				ref = resourceReference.Type
			}
		} else if resourceReference.ChildType != "" {
			ref = resourceReference.ChildType
		}
		if regexp, ok := resourceReferenceMap[ref]; ok {
			if desc.IsList() {
				exprs = append(exprs, fmt.Sprintf(`%s.all(s, s.matches("%s"))`, desc.TextName(), regexp))
			} else if desc.Kind() == protoreflect.StringKind {
				exprs = append(exprs, fmt.Sprintf(`%s.matches("%s")`, desc.TextName(), regexp))
			}
		} else {
			return nil, fmt.Errorf(`cannot find type "%s"`, ref)
		}
	}
	required := false
	if !config.RequiredSupportDisabled {
		for _, behavior := range proto.GetExtension(desc.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior) {
			if behavior == annotations.FieldBehavior_REQUIRED {
				required = true
			}
		}
	}
	var vpgr ValidateProgram
	if len(exprs) > 0 {
		if pgr, err := BuildValidateProgram(exprs, config, envOpt, imports...); err != nil {
			return nil, err
		} else {
			vpgr = pgr
		}
	}
	return &fieldValidateProgram{program: vpgr, required: required}, nil
}

func isDefaultValue(m proto.Message, fdesc protoreflect.FieldDescriptor) bool {
	pf := m.ProtoReflect().Get(fdesc)
	if fdesc.IsList() {
		return pf.List() == nil || pf.List().Len() == 0
	} else if fdesc.IsMap() {
		return pf.Map() == nil || pf.Map().Len() == 0
	} else {
		switch fdesc.Kind() {
		case protoreflect.MessageKind, protoreflect.GroupKind:
			return !pf.Message().IsValid()
		case protoreflect.EnumKind:
			return pf.Enum() == fdesc.Default().Enum()
		case protoreflect.BytesKind:
			return pf.Bytes() == nil || len(pf.Bytes()) == 0
		default:
			return pf.Interface() == fdesc.Default().Interface()
		}
	}
}
