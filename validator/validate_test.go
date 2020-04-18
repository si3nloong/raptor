package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {
	var i struct {
		Type   string `validate:"required"`
		Email  string `validate:"required_if=Type eq TYPE"`
		Int    int    `validate:"required"`
		Int64  int64  `validate:"required"`
		Nested struct {
			Num         int `validate:"required_if=Int eq 1"`
			GreaterThan int `validate:"required_if=Int64 gt 100"`
		}
	}

	{
		i.Type = "TYPE"
		i.Email = "test@gmail.com"
		i.Nested.Num = 100
		i.Nested.GreaterThan = 100
		i.Int = 1
		i.Int64 = 1123123123

		err := Validate("form", i)
		require.NoError(t, err)
	}

	// {
	// 	i.Int = 0
	// 	err := Validate("form", i)
	// 	require.Error(t, err)
	// }
}
