package raptor

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBind(t *testing.T) {
	var (
		err error
	)

	var i struct {
		Name           []byte  `query:"name"`
		Signature      string  `query:"signature"`
		Amount         float64 `query:"amount"`
		DiscountAmount int64   `query:"discountAmount"`
		Total          uint    `query:"total"`
	}

	v := reflect.ValueOf(&i)
	err = bindQuery("query", v, map[string][]string{
		"name":           {"Hello World"},
		"signature":      {"D86029B21A0DF5E4AFBA48D0FB9861393B278AAC"},
		"amount":         {"19.85"},
		"discountAmount": {"-105"},
		"total":          {"100"},
	})

	require.NoError(t, err)
	require.Equal(t, "Hello World", string(i.Name))
	require.Equal(t, "D86029B21A0DF5E4AFBA48D0FB9861393B278AAC", i.Signature)
	require.Equal(t, float64(19.85), i.Amount)
	require.Equal(t, int64(-105), i.DiscountAmount)
	require.Equal(t, uint(100), i.Total)
}
