package validate

import (
	"fmt"
	reflect "reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ValidationError struct {
	Descriptor protoreflect.Descriptor
}

func (e *ValidationError) Error() string {
	if e.Descriptor != nil {
		return fmt.Sprintf(`validation failed on "%s"`, e.Descriptor.FullName())
	}
	return "validation failed"
}

func (e *ValidationError) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	return nil, e
}

func (e *ValidationError) ConvertToType(typeVal ref.Type) ref.Val {
	return e
}

func (e *ValidationError) Equal(other ref.Val) ref.Val {
	return e
}

func (e *ValidationError) String() string {
	return e.Error()
}

func (e *ValidationError) Type() ref.Type {
	return types.ErrType
}

func (e *ValidationError) Value() interface{} {
	return e
}

type MethodValidationError struct {
	AttributeContext *attribute_context.AttributeContext
}

func (e *MethodValidationError) Error() string {
	if e.AttributeContext != nil && e.AttributeContext.Api != nil {
		return fmt.Sprintf(`validation failed on "%s"`, e.AttributeContext.Api.Operation)
	}
	return "validation failed"
}

func (e *MethodValidationError) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	return nil, e
}

func (e *MethodValidationError) ConvertToType(typeVal ref.Type) ref.Val {
	return e
}

func (e *MethodValidationError) Equal(other ref.Val) ref.Val {
	return e
}

func (e *MethodValidationError) String() string {
	return e.Error()
}

func (e *MethodValidationError) Type() ref.Type {
	return types.ErrType
}

func (e *MethodValidationError) Value() interface{} {
	return e
}
