package connect

import (
	"context"

	"github.com/Neakxs/protocel/authorize"
	"github.com/bufbuild/connect-go"
)

func NewConnectUnaryInterceptor(authzHandler authorize.AuthzInterceptor) connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(uf connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
			if err := authzHandler.Authorize(ctx, ar.Spec().Procedure, ar.Header(), ar.Any()); err != nil {
				return nil, connect.NewError(connect.CodePermissionDenied, err)
			}
			return uf(ctx, ar)
		}
	})
}
