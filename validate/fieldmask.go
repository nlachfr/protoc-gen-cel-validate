package validate

import (
	"context"
	"fmt"
	reflect "reflect"
	"strings"

	"github.com/google/cel-go/cel"
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

func ValidateWithMask(ctx context.Context, m proto.Message, fm *fieldmaskpb.FieldMask, validationMap map[string]cel.Program) error {
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
					if !reflect.ValueOf(m.ProtoReflect().Get(fdesc).Interface()).IsZero() {
						if pgr, ok := validationMap[fdesc.TextName()]; ok {
							if val, _, err := pgr.ContextEval(ctx, vars); err != nil {
								return err
							} else if !types.IsBool(val) || !val.Value().(bool) {
								return fmt.Errorf(`validation failed on %s`, fdesc.FullName())
							}
						}
					} else {
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

func buildProgramVars(m proto.Message) interface{} {
	res := map[string]interface{}{}
	fields := m.ProtoReflect().Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		res[field.TextName()] = m.ProtoReflect().Get(field)
	}
	return res
}
