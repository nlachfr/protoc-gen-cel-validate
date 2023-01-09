package errors

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type err interface {
	error
	ConvertToNative(typeDesc reflect.Type) (interface{}, error)
	ConvertToType(typeVal ref.Type) ref.Val
	Equal(other ref.Val) ref.Val
	String() string
	Type() ref.Type
	Value() interface{}
	Unwrap() error
}

type ValidateError interface {
	err
	GetAttributeContext() *attribute_context.AttributeContext
	GetMessage() proto.Message
	GetDescriptor() protoreflect.Descriptor
}

func New(message proto.Message, desc protoreflect.Descriptor, ctx *attribute_context.AttributeContext) ValidateError {
	return &validateError{Message: message, Descriptor: desc, AttributeContext: ctx}
}

func Wrap(err error, message proto.Message, desc protoreflect.Descriptor, ctx *attribute_context.AttributeContext) ValidateError {
	return &validateError{Err: err, Message: message, Descriptor: desc, AttributeContext: ctx}
}

type validateError struct {
	Err              error
	AttributeContext *attribute_context.AttributeContext
	Message          proto.Message
	Descriptor       protoreflect.Descriptor
}

func (e *validateError) GetAttributeContext() *attribute_context.AttributeContext {
	return e.AttributeContext
}
func (e *validateError) GetMessage() proto.Message {
	return e.Message
}
func (e *validateError) GetDescriptor() protoreflect.Descriptor {
	return e.Descriptor
}

func (e *validateError) Error() string {
	if e.Descriptor != nil {
		if e.Err != nil {
			return fmt.Sprintf(`validation failed on "%s": %s`, e.Descriptor.FullName(), e.Err.Error())
		}
		return fmt.Sprintf(`validation failed on "%s"`, e.Descriptor.FullName())
	}
	return "validation failed"
}

func (e *validateError) Unwrap() error {
	if e.Err != nil {
		return e.Err
	}
	return nil
}

func (e *validateError) ConvertToNative(typeDesc reflect.Type) (interface{}, error) { return nil, e }
func (e *validateError) ConvertToType(typeVal ref.Type) ref.Val                     { return e }
func (e *validateError) Equal(other ref.Val) ref.Val                                { return e }
func (e *validateError) String() string                                             { return e.Error() }
func (e *validateError) Type() ref.Type                                             { return types.ErrType }
func (e *validateError) Value() interface{}                                         { return e }
