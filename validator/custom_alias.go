package validator

import (
	validator "github.com/go-playground/validator/v10"
)

// Register Custom Alias
func registerAlias(validate *validator.Validate) {
	validate.RegisterAlias("phonenumber", "numeric,min=6,max=20")
}
