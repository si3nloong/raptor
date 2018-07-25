# Raptor Web Framework

Inspired by Laravel and Iris

Why Raptor web framework? We love fasthttp, we love speed.
Basically Raptor is using [fasthttp](https://github.com/erikdubbelboer/fasthttp), [fasthttprouter](https://github.com/thehowl/fasthttprouter), [ffjson](https://github.com/pquerna/ffjson) packages under the hood.

## Installation

The only requirement is the Go Programming Language

```bash
$ go get -u github.com/si3nloong/raptor
```

## Quick Start

```go
package main

import (
  "github.com/si3nloong/raptor"
)

func main() {
  r := raptor.New()
  r.GET("/", func(c *raptor.Context) error {
    return c.Response().Custom(raptor.Map{"message":"hello world"})
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
    return c.Response().Custom(raptor.Map{"message":"hello world"})
  })

  open := raptor.New()
  open.GET("/", func(c *raptor.Context) error {
    return c.Response().Custom(raptor.Map{"message":"hello world"})
  })

  hosts["api.domain.com"] = api.Handler()
  hosts["open.domain.com"] = open.Handler()

  r := raptor.New()
  r.Start(":8080", hosts.Routing)
}
```

## Variable Binding

```go
  api := raptor.New()
  api.GET("/", func(c *raptor.Context) error {
    var i struct {
	    Name string `json:"name" xml:"name" query:"name"`
	  }

    if err := c.Bind(&i); err != nil {
      return c.Response().BadRequest(c.NewAPIError(err))
    }

	  return c.Response().Custom(raptor.Map{"message":"hello world"})
  })
  api.Start(":8080")
```

## Validation

```go
  api := raptor.New()
  api.GET("/", func(c *raptor.Context) error {
    var i struct {
	    Name string `json:"name" xml:"name" query:"name"`
    }

    if err := c.Bind(&i); err != nil {
      return c.Response().BadRequest(c.NewAPIError(err))
    }

    if message, err := c.Validate(&i); err != nil {
      return c.Response().UnprocessableEntity(c.NewAPIError(err, "", message))
    }

	  return c.Response().Custom(raptor.Map{"message":"hello world"})
  })
  api.Start(":8080")
```

## Error Handling

## Custom Error

[MIT License](https://github.com/si3nloong/raptor/blob/master/LICENSE)
