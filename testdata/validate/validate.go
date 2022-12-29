package validate

import (
	"context"
	"fmt"
	"regexp"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (r *TestRpcRequest) Validate(ctx context.Context) error {
	return r.ValidateWithMask(ctx, &fieldmaskpb.FieldMask{Paths: []string{"*"}})
}

func (r *TestRpcRequest) ValidateWithMask(ctx context.Context, fm *fieldmaskpb.FieldMask) error {
	if r == nil {
		return fmt.Errorf("nil request")
	}
	if r.Ref == "" {
		return fmt.Errorf("required failed")
	} else if ok, _ := regexp.MatchString(`^refs/[a-zA-Z0-9]+$`, r.Ref); !ok {
		return fmt.Errorf("resource reference failed")
	}
	return nil
}

func (r *Nested) Validate(ctx context.Context) error {
	return r.ValidateWithMask(ctx, &fieldmaskpb.FieldMask{Paths: []string{"*"}})
}

func (r *Nested) ValidateWithMask(ctx context.Context, fm *fieldmaskpb.FieldMask) error {
	if fm == nil {
		return fmt.Errorf("nil fieldmask")
	}
	if r == nil {
		return fmt.Errorf("nil nested")
	}
	for _, path := range fm.Paths {
		switch path {
		case "name":
			if len(r.Name) == 0 {
				return fmt.Errorf("name is empty")
			}
		case "value":
			if r.Value != "value" {
				return fmt.Errorf("invalid value")
			}
		case "ref.name":
			if r.Ref == nil || r.Ref.Name != "name" {
				return fmt.Errorf("invalid ref.name")
			}
		}
	}
	return nil
}
