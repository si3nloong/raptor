# Raptor Web Framework

Inspired by Laravel and Iris

## Installation

The only requirement is the Go Programming Language

```go
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

	return c.SuccessString("application/json", `{"message":"hello world"}`)
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

	return c.SuccessString("application/json", `{"message":"hello world"}`)
  })
  api.Start(":8080")
```

## Error Handling

## Custom Error

[MIT License](https://github.com/si3nloong/raptor/blob/master/LICENSE)
