package middleware

import (
	"fmt"
	"log"

	"github.com/si3nloong/raptor"
)

// RecoverConfig :
type RecoverConfig struct {
	Skipper Skipper
}

// Recover :
func Recover(config ...SecureConfig) raptor.MiddlewareFunc {
	c := DefaultSecureConfig
	if len(config) > 0 {
		c = config[0]
	}
	return recoverWithConfig(c)
}

func recoverWithConfig(config SecureConfig) raptor.MiddlewareFunc {
	return func(next raptor.HandlerFunc) raptor.HandlerFunc {
		return func(ctx *raptor.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, isOk := r.(error)
					if !isOk {
						err = fmt.Errorf("%v", r)
					}
					// stack := make([]byte, config.StackSize)
					// length := runtime.Stack(stack, !config.DisableStackAll)
					// if !config.DisablePrintStack {
					// 	c.Logger().Printf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					// }
					log.Println(err)
					// c.Error(err)
				}
			}()
			return next(ctx)
		}
	}
}
