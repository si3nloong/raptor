package test

import (
	"testing"

	"github.com/si3nloong/raptor"
)

func TestServer(t *testing.T) {
	r := raptor.New()
	r.StaticGzip("/gzip", "assets/file.js.gz")
	r.Start(":9001")
}
