package twirp

import (
	"context"

	"github.com/Neakxs/protocel/authorize"
	"github.com/twitchtv/twirp"
)

func NewTwirpInterceptor(authzHandler authorize.AuthzInterceptor) twirp.Interceptor {
	return func(m twirp.Method) twirp.Method {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			fullMethod, ok := twirp.MethodName(ctx)
			if ok {
				headers, _ := twirp.HTTPRequestHeaders(ctx)
				if err := authzHandler.Authorize(ctx, fullMethod, headers, request); err != nil {
					return nil, twirp.PermissionDenied.Error(err.Error())
				}
			}
			return m(ctx, request)
		}
	}
}
