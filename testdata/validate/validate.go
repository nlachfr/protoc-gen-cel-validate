package validate

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

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
