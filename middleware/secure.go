package middleware

import (
	"runtime"

	"github.com/si3nloong/raptor"
)

// SecureConfig :
type SecureConfig struct {
}

// DefaultSecureConfig :
var (
	DefaultSecureConfig = SecureConfig{}
)

// Secure :
func Secure(config ...SecureConfig) raptor.MiddlewareFunc {
	c := DefaultSecureConfig
	if len(config) > 0 {
		c = config[0]
	}
	return secureWithConfig(c)
}

// Secure :
func secureWithConfig(config SecureConfig) raptor.MiddlewareFunc {
	return func(next raptor.HandlerFunc) raptor.HandlerFunc {
		return func(ctx *raptor.Context) error {
			ctx.RequestCtx.Response.Header.Set(raptor.HeaderXFrameOptions, "deny")
			ctx.RequestCtx.Response.Header.Set(raptor.HeaderXContentTypeOptions, "nosniff")
			ctx.RequestCtx.Response.Header.Set(raptor.HeaderXXSSProtection, "1; mode=block")
			ctx.RequestCtx.Response.Header.Set(raptor.HeaderXPoweredBy, runtime.Version())
			return next(ctx)
		}
	}
}
