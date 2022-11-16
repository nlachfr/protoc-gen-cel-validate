package validate

import (
	"context"
	"testing"

	"github.com/Neakxs/protocel/testdata/validate"
	"github.com/google/cel-go/cel"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestValidateWithMask(t *testing.T) {
	noDepthMap := map[string]cel.Program{}
	noDepthMap["nested"], _ = BuildValidateProgram(`true`, &validate.TestRpcRequest{}, nil)
	fmNoDepthMap := map[string]cel.Program{}
	fmNoDepthMap["nested"], _ = BuildValidateProgram(`nested.validateWithMask(fm)`, &validate.TestRpcRequest{}, nil)
	tests := []struct {
		Name          string
		ValidationMap map[string]cel.Program
		Tests         []struct {
			Name    string
			Message proto.Message
			WantErr bool
		}
	}{
		{
			Name:          "No depth",
			ValidationMap: noDepthMap,
			Tests: []struct {
				Name    string
				Message proto.Message
				WantErr bool
			}{
				{
					Name:    "Empty ref",
					Message: &validate.TestRpcRequest{},
					WantErr: true,
				},
				{
					Name: "OK (ref)",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
					},
					WantErr: false,
				},
			},
		},
		{
			Name:          "Fieldmask validation without depth",
			ValidationMap: fmNoDepthMap,
			Tests: []struct {
				Name    string
				Message proto.Message
				WantErr bool
			}{
				{
					Name: "No fieldmask",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
					},
					WantErr: true,
				},
				{
					Name: "Fieldmask with empty struct fields",
					Message: &validate.TestRpcRequest{
						Ref:    "ref",
						Nested: &validate.Nested{},
						Fm:     &fieldmaskpb.FieldMask{Paths: []string{"name"}},
					},
					WantErr: true,
				},
				{
					Name: "Fieldmask with one invalid field",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
						Nested: &validate.Nested{
							Name: "name",
						},
						Fm: &fieldmaskpb.FieldMask{Paths: []string{"name", "value"}},
					},
					WantErr: true,
				},
				{
					Name: "Fieldmask with invalid nested field",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
						Nested: &validate.Nested{
							Ref: &validate.RefMultiple{
								Name: "noname",
							},
						},
						Fm: &fieldmaskpb.FieldMask{Paths: []string{"ref.name"}},
					},
					WantErr: true,
				},
				{
					Name: "Fieldmask with valid fields",
					Message: &validate.TestRpcRequest{
						Ref: "ref",
						Nested: &validate.Nested{
							Name:  "name",
							Value: "value",
							Ref: &validate.RefMultiple{
								Name: "name",
							},
						},
						Fm: &fieldmaskpb.FieldMask{Paths: []string{"name", "value", "ref.name"}},
					},
					WantErr: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for _, ttt := range tt.Tests {
				t.Run(ttt.Name, func(t *testing.T) {
					err := ValidateWithMask(context.Background(), ttt.Message, &fieldmaskpb.FieldMask{Paths: []string{"*"}}, tt.ValidationMap)
					if (err != nil && !ttt.WantErr) || (err == nil && ttt.WantErr) {
						t.Errorf("wantErr %v, got %v", ttt.WantErr, err)
					}
				})
			}
		})
	}
	if err := ValidateWithMask(context.Background(), &validate.TestRpcRequest{Ref: "ref"}, nil, noDepthMap); err != nil {
		t.Errorf("wantErr false, got %v", err)
	}
	if err := ValidateWithMask(context.Background(), &validate.TestRpcRequest{Ref: "ref"}, &fieldmaskpb.FieldMask{Paths: []string{"ref", "nested"}}, noDepthMap); err != nil {
		t.Errorf("wantErr false, got %v", err)
	}
	if err := ValidateWithMask(context.Background(), nil, nil, nil); err == nil {
		t.Errorf("wantErr true, got <nil>")
	}
}
