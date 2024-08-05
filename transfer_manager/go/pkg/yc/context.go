package yc

import (
	"context"

	"github.com/doublecloud/tross/transfer_manager/go/pkg/contextutil"
)

var (
	withUserAuthCtxKey = contextutil.NewContextKey()
)

func WithUserAuth(ctx context.Context) context.Context {
	return context.WithValue(ctx, withUserAuthCtxKey, true)
}

func IsWithUserAuth(ctx context.Context) bool {
	value, ok := ctx.Value(withUserAuthCtxKey).(bool)
	return ok && value
}