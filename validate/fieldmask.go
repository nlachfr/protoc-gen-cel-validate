package validate

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/cel-go/common/types"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type Validater interface {
	Validate(ctx context.Context) error
	ValidateWithMask(ctx context.Context, fm *fieldmaskpb.FieldMask) error
}

func ValidateWithMask(ctx context.Context, m proto.Message, fm *fieldmaskpb.FieldMask, validationMap map[string]*Program, enforceRequired bool) error {
	if validationMap == nil {
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
	vars := buildProgramVars(m)
	for i := 0; i < mdesc.Fields().Len(); i++ {
		fdesc := mdesc.Fields().Get(i)
		if paths, ok := pathsMap[fdesc.TextName()]; ok {
			subs := []string{}
			for j := 0; j < len(paths); j++ {
				if paths[j] == "" {
					if !isDefaultValue(m, fdesc) {
						if pgr, ok := validationMap[fdesc.TextName()]; ok {
							for _, p := range pgr.rules {
								if val, _, err := p.ContextEval(ctx, vars); err != nil {
									return err
								} else if !types.IsBool(val) || !val.Value().(bool) {
									return fmt.Errorf(`validation failed on %s`, fdesc.FullName())
								}
							}
						}
					} else if enforceRequired {
						for _, behavior := range proto.GetExtension(fdesc.Options(), annotations.E_FieldBehavior).([]annotations.FieldBehavior) {
							if behavior == annotations.FieldBehavior_REQUIRED {
								return fmt.Errorf(`validation failed on %s`, fdesc.FullName())
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
						return err
					}
				}
			}
		}
	}
	return nil
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

func buildProgramVars(m proto.Message) interface{} {
	res := map[string]interface{}{}
	fields := m.ProtoReflect().Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		res[field.TextName()] = m.ProtoReflect().Get(field)
	}
	return res
}
