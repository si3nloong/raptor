package middleware

import (
	"github.com/si3nloong/raptor"
)

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(*raptor.Context) bool
)

// DefaultSkipper :
func DefaultSkipper(*raptor.Context) bool {
	return false
}
