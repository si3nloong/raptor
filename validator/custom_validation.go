package validator

import (
	"fmt"
	"reflect"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

// CountryCodes : available country codes
var countryCodes = map[string]string{
	"60":  "Malaysia",
	"86":  "China",
	"971": "UAE",
}

// Register Custom Validator
func registerValidation(validate *validator.Validate) {
	validate.RegisterValidation("countrycode", validateCountryCode)
	validate.RegisterValidation("requiredif", validateRequiredIf)
}

func validateCountryCode(fl validator.FieldLevel) bool {
	_, isOK := countryCodes[fl.Field().String()]
	return isOK
}

func validateRequiredIf(fl validator.FieldLevel) bool {
	f := strings.Split(fl.Param(), ":")
	field := f[0]
	value := f[1]
	divefield := strings.Split(field, ".")

	q := reflect.Indirect(fl.Top())
	// switch v.Kind() {
	// case reflect.Slice, reflect.Array:

	// }
	// log.Println("log:", v)
	// log.Println("q: ", v.FieldByIndex([]int{0}).Index(0).Elem())

	// log.Println("log:", v.Index(0))
	// log.Println("debug: ", fl.Top())
	// q := fl.Top().Elem()

	for _, field := range divefield {
		q = q.FieldByName(field)
	}

	if fmt.Sprintf("%v", q.Interface()) == reflect.ValueOf(value).Interface() {
		validate := validator.New()
		if err := validate.Var(fl.Field().Interface(), "required"); err != nil {
			return false
		}
		return true
	}
	return true
}