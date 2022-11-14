package authorize

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

type AuthzInterceptor interface {
	Authorize(ctx context.Context, method string, headers http.Header, request interface{}) error
}

func NewAuthzInterceptor(methodProgramMapping map[string]cel.Program) AuthzInterceptor {
	return &authzInterceptor{
		methodProgramMapping: methodProgramMapping,
	}
}

type authzInterceptor struct {
	methodProgramMapping map[string]cel.Program
}

func (i *authzInterceptor) Authorize(ctx context.Context, method string, headers http.Header, request interface{}) error {
	if pgr, ok := i.methodProgramMapping[method]; ok {
		if val, _, err := pgr.ContextEval(ctx, map[string]interface{}{
			"headers": headers,
			"request": request,
		}); err != nil {
			return err
		} else if !types.IsBool(val) || !val.Value().(bool) {
			return fmt.Errorf(`permission denied on "%s"`, method)
		}
	}
	return nil
}
