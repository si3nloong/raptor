package middleware

import (
	"fmt"
	"strings"

	"github.com/si3nloong/raptor"
)

// CORSConfig :
type CORSConfig struct {
	Skipper          Skipper
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	MaxAge           int      `yaml:"max_age"`
}

var (
	// DefaultCORSConfig is the default CORS middleware config.
	DefaultCORSConfig = CORSConfig{
		Skipper:       DefaultSkipper,
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{raptor.HEAD, raptor.GET, raptor.POST, raptor.PUT, raptor.PATCH, raptor.DELETE},
		AllowHeaders:  []string{raptor.HeaderOrigin, raptor.HeaderAccept, raptor.HeaderContentType, raptor.HeaderAuthorization, raptor.HeaderContentLength, raptor.HeaderAcceptCharset, "X-Requested-With"},
		ExposeHeaders: []string{raptor.HeaderContentDisposition},
	}
)

// CORS :
func CORS(config ...CORSConfig) raptor.MiddlewareFunc {
	c := DefaultCORSConfig
	if len(config) > 0 {
		c = config[0]
	}
	if c.Skipper == nil {
		c.Skipper = DefaultSkipper
	}
	if c.AllowHeaders == nil {
		c.AllowHeaders = DefaultCORSConfig.AllowHeaders
	}
	if c.AllowMethods == nil {
		c.AllowMethods = DefaultCORSConfig.AllowMethods
	}
	if c.ExposeHeaders == nil {
		c.ExposeHeaders = DefaultCORSConfig.ExposeHeaders
	}
	return corsWithConfig(c)
}

func corsWithConfig(config CORSConfig) raptor.MiddlewareFunc {
	allowOrigins := make(map[string]bool)
	for _, o := range config.AllowOrigins {
		allowOrigins[strings.TrimSpace(o)] = true
	}
	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")

	return func(next raptor.HandlerFunc) raptor.HandlerFunc {
		return func(ctx *raptor.Context) error {
			if config.Skipper != nil && config.Skipper(ctx) {
				return next(ctx)
			}

			origin := string(ctx.Request.Header.Peek(raptor.HeaderOrigin))
			if _, isOk := allowOrigins[origin]; !isOk && !allowOrigins["*"] {
				origin = ""
			}

			ctx.RequestCtx.Response.Header.Add(raptor.HeaderVary, raptor.HeaderOrigin)
			ctx.RequestCtx.Response.Header.Set(raptor.HeaderAccessControlAllowOrigin, origin)
			if config.AllowCredentials {
				ctx.RequestCtx.Response.Header.Set(raptor.HeaderAccessControlAllowCredentials, "true")
			}

			if !ctx.IsMethod(raptor.OPTIONS) {
				ctx.RequestCtx.Response.Header.Set(raptor.HeaderAccessControlExposeHeaders, exposeHeaders)
				return next(ctx)
			}

			ctx.RequestCtx.Response.Header.Add(raptor.HeaderVary, raptor.HeaderAccessControlRequestMethod)
			ctx.RequestCtx.Response.Header.Add(raptor.HeaderVary, raptor.HeaderAccessControlRequestHeaders)
			ctx.RequestCtx.Response.Header.Set(raptor.HeaderAccessControlExposeHeaders, exposeHeaders)
			ctx.RequestCtx.Response.Header.Set(raptor.HeaderAccessControlAllowMethods, allowMethods)

			if allowHeaders != "" {
				ctx.RequestCtx.Response.Header.Set(raptor.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				requestHeaders := string(ctx.Request.Header.Peek(raptor.HeaderAccessControlRequestHeaders))
				if requestHeaders != "" {
					ctx.RequestCtx.Response.Header.Set(raptor.HeaderAccessControlAllowHeaders, requestHeaders)
				}
			}

			if config.MaxAge > 0 {
				ctx.RequestCtx.Response.Header.Set(raptor.HeaderAccessControlMaxAge, fmt.Sprintf("%d", config.MaxAge))
			}

			return ctx.Response().NoContent()
		}
	}
}
