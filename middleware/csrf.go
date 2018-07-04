package middleware

import (
	"github.com/si3nloong/raptor"
	"github.com/valyala/fasthttp"
)

// CSRFConfig :
type CSRFConfig struct {
	Skipper     Skipper
	TokenLength uint8  `yaml:"token_length"`
	TokenLookup string `yaml:"token_lookup"`
	ContextKey  string `yaml:"context_key"`
	Cookie      *fasthttp.Cookie
}

var (
	// DefaultCSRFConfig :
	DefaultCSRFConfig = CSRFConfig{
		Skipper: DefaultSkipper,
	}
)

// CSRF :
func CSRF(config ...CSRFConfig) raptor.MiddlewareFunc {
	c := DefaultCSRFConfig
	if len(config) > 0 {
		c = config[0]
	}
	if c.Skipper == nil {
		c.Skipper = DefaultSkipper
	}
	return csrfWithConfig(c)
}

func csrfWithConfig(config CSRFConfig) raptor.MiddlewareFunc {
	return func(next raptor.HandlerFunc) raptor.HandlerFunc {
		return func(ctx *raptor.Context) error {
			if config.Skipper != nil && config.Skipper(ctx) {
				return next(ctx)
			}

			return next(ctx)
		}
	}
}
