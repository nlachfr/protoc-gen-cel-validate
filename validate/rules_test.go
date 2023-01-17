package validate

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		Name       string
		BuildOpts  []BuildOption
		MethodDesc protoreflect.MethodDescriptor
		Request    proto.Message
		WantErr    bool
	}{}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

		})
	}
}
