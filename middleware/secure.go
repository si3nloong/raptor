package middleware

import (
	"runtime"

	"github.com/si3nloong/raptor"
)

// SecureConfig :
type SecureConfig struct {
	Skipper               Skipper
	XSSProtection         string `yaml:"xss_protection"`
	ContentTypeNosniff    string `yaml:"content_type_nosniff"`
	XFrameOptions         string `yaml:"x_frame_options"`
	HSTSMaxAge            int    `yaml:"hsts_max_age"`
	HSTSExcludeSubdomains bool   `yaml:"hsts_exclude_subdomains"`
	ContentSecurityPolicy string `yaml:"content_security_policy"`
}

// DefaultSecureConfig :
var (
	DefaultSecureConfig = SecureConfig{
		Skipper:            DefaultSkipper,
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}
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
			if config.Skipper != nil && config.Skipper(ctx) {
				return next(ctx)
			}

			if config.XFrameOptions != "" {
				ctx.RequestCtx.Response.Header.Set(raptor.HeaderXFrameOptions, config.XFrameOptions)
			}

			if config.ContentTypeNosniff != "" {
				ctx.RequestCtx.Response.Header.Set(raptor.HeaderXContentTypeOptions, config.ContentTypeNosniff)
			}

			if config.XSSProtection != "" {
				ctx.RequestCtx.Response.Header.Set(raptor.HeaderXXSSProtection, config.XSSProtection)
			}

			if config.ContentSecurityPolicy != "" {
				ctx.RequestCtx.Response.Header.Set(raptor.HeaderContentSecurityPolicy, config.ContentSecurityPolicy)
			}

			ctx.RequestCtx.Response.Header.Set(raptor.HeaderXPoweredBy, runtime.Version())
			return next(ctx)
		}
	}
}
