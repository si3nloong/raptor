package validator

import (
	validator "gopkg.in/go-playground/validator.v9"
)

// Register Custom Alias
func registerAlias(validate *validator.Validate) {
	validate.RegisterAlias("phonenumber", "numeric,min=6,max=20")
}
