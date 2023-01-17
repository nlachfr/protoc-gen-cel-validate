package validate

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// GetExtension wraps proto.GetExtension function with dynamicpb.Message support
func GetExtension(m protoreflect.ProtoMessage, xt protoreflect.ExtensionType) interface{} {
	if m == nil {
		return proto.GetExtension(m, xt)
	}
	switch t := m.ProtoReflect().Get(xt.TypeDescriptor()).Interface().(type) {
	case *dynamicpb.Message:
		mm := xt.InterfaceOf(xt.New()).(proto.Message)
		if raw, err := proto.Marshal(t); err != nil {
			panic(err)
		} else if err = proto.Unmarshal(raw, mm); err != nil {
			panic(err)
		}
		return mm
	default:
		if t == nil {
			return xt.InterfaceOf(xt.Zero())
		}
		return proto.GetExtension(m, xt)
	}
}

func IsDefaultValue(m proto.Message, fdesc protoreflect.FieldDescriptor) bool {
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
