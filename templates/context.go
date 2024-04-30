package templates

import "context"

type ContextKey string

const HxRequestContextKey ContextKey = "hx-request"

func IsHxRequest(ctx context.Context) bool {
	if isHxRequest, err := ctx.Value(HxRequestContextKey).(bool); err {
		return isHxRequest
	}
	return false
}
