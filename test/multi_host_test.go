package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/si3nloong/raptor"
	"github.com/si3nloong/raptor/middleware"
)

type host map[string]raptor.HandlerFunc

var hosts = make(host)

// HostRouting Implement a CheckHost method on our new type
func (hs host) HostRouting(ctx *raptor.Context) error {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.
	if cb := hosts[string(ctx.Host())]; cb != nil {
		return cb(ctx)
	}

	return ctx.Response().Forbidden(fmt.Errorf("not able to access"))
}

func TestMultiHost(t *testing.T) {

	open := raptor.New()
	open.Use(middleware.CORS())
	open.GET("/test", func(c *raptor.Context) error {
		log.Println("open!!!!")
		return c.Response().Custom(map[string]interface{}{
			"bearer": "Open",
		})
	})

	api := raptor.New()
	api.Use(middleware.CORS())
	api.Use(middleware.Secure())
	api.GET("/test", func(c *raptor.Context) error {
		log.Println("test")
		return c.Response().Custom(map[string]interface{}{
			"Test": "aslkjdlkasjdjlsajlkd",
		})
	})

	hosts["open.wetix.my:9000"] = open.Handler()
	hosts["api.wetix.my:9000"] = api.Handler()

	r := raptor.New()
	r.Start(":9000", hosts.HostRouting)
}
