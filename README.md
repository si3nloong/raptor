## Quick Start

```go
package main

import (
	"github.com/si3nloong/raptor"
)

func main() {
	r := raptor.New()
	r.GET("/", func(c *raptor.Context) error {
		return c.SuccessString("application/json", `{"message":"hello world"}`)
	})
	r.Start(":8080")
}
```

## Multiple Hosts

```go
import (
    "github.com/si3nloong/raptor"
)

type host map[string]raptor.HandlerFunc

// Routing is to route to specific domain
func (hs host) Routing(ctx *raptor.Context) error {
	if cb := hosts[string(ctx.Host())]; cb != nil {
		cb = middleware.CORS(corsConfig)(cb)
		return cb(ctx)
	}
	return ctx.Response().NotFound(fmt.Errorf("page not found"))
}

func main() {
    api := raptor.New()
    api.GET("/", func(c *raptor.Context) error {
        return c.SuccessString("application/json", `{"message":"hello world"}`)
    })

    open := raptor.New()
    open.GET("/", func(c *raptor.Context) error {
        return c.SuccessString("application/json", `{"message":"hello world"}`)
    })

    hosts["api.domain.com"] = api.Handler()
    hosts["open.domain.com"] = open.Handler()

	r := raptor.New()
	r.Start(":8080", hosts.Routing)
}
```

## Error Handling

packages we use :

- github.com/valyala/fasthttp
- github.com/buaazp/fasthttprouter
