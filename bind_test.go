package raptor

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestBind(t *testing.T) {
	var i struct {
		Name           []byte  `query:"name"`
		Signature      string  `query:"signature"`
		Amount         float64 `query:"amount"`
		DiscountAmount int64   `query:"discountAmount"`
		Total          uint    `query:"total"`
	}

	v := reflect.ValueOf(&i)
	if err := bindQuery("query", v, map[string][]string{
		"name":           []string{"Hello World"},
		"signature":      []string{"D86029B21A0DF5E4AFBA48D0FB9861393B278AAC"},
		"amount":         []string{"19.85"},
		"discountAmount": []string{"-105"},
		"total":          []string{"100"},
	}); err != nil {
		log.Fatal(err)
	}

	if string(i.Name) != "Hello World" {
		log.Fatal(fmt.Errorf("Unexpected bind value %s", i.Name))
	}

	if i.Signature != "D86029B21A0DF5E4AFBA48D0FB9861393B278AAC" {
		log.Fatal(fmt.Errorf("Unexpected bind value %v", i.Signature))
	}

	if i.Amount != 19.85 {
		log.Fatal(fmt.Errorf("Unexpected bind value %v", i.Amount))
	}

	if i.DiscountAmount != -105 {
		log.Fatal(fmt.Errorf("Unexpected bind value %v", i.DiscountAmount))
	}

	if i.Total != 100 {
		log.Fatal(fmt.Errorf("Unexpected bind value %v", i.Total))
	}
}
