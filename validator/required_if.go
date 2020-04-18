package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

func (v *Validator) validateRequiredIf(fl validator.FieldLevel) bool {
	param := fl.Param()
	paths := strings.SplitN(param, " ", 3)
	if len(paths) != 3 {
		panic("invalid parameter for required_if")
	}

	for idx := range paths {
		paths[idx] = strings.TrimSpace(paths[idx])
	}

	// v, k, nullable, found := fl.GetStructFieldOKAdvanced2(fl.Top(), paths[0])
	val, k, _, found := fl.GetStructFieldOKAdvanced2(fl.Top(), paths[0])
	if !found {
		panic(fmt.Sprintf("invalid parent field name %s for field %s", paths[0], fl.FieldName()))
	}

	var valid bool
	switch paths[1] {
	case "eq":
		valid = validateEqual(val, k, paths)

	case "ne":
		valid = validateNotEqual(val, k, paths)

	case "gt":
		valid = validateGreaterThan(val, k, paths)

	case "gte":
		valid = validateGreaterThanOrEqual(val, k, paths)

	case "lt":
		valid = validateLesserThan(val, k, paths)

	case "lte":
		valid = validateLesserThanOrEqual(val, k, paths)

	default:
		panic(fmt.Sprintf("unsupported operator %s", paths[1]))
	}

	if !valid {
		return true
	}

	if err := v.Var(fl.Field().Interface(), "required"); err != nil {
		return false
	}
	return true
}

func validateEqual(v reflect.Value, k reflect.Kind, paths []string) bool {
	switch k {
	case reflect.String:
		return v.String() == paths[2]
	case reflect.Bool:
		flag, err := strconv.ParseBool(paths[2])
		if err != nil {
			panic(err)
		}
		return v.Bool() == flag
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Int() == num
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Uint() == num
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(paths[2], 64)
		if err != nil {
			panic(err)
		}
		return v.Float() == f
	default:
		panic(fmt.Sprintf("unsupported data type %s", k))
	}
}

func validateNotEqual(v reflect.Value, k reflect.Kind, paths []string) bool {
	switch k {
	case reflect.String:
		return v.String() != paths[2]
	case reflect.Bool:
		flag, err := strconv.ParseBool(paths[2])
		if err != nil {
			panic(err)
		}
		return v.Bool() != flag
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Int() != num
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Uint() != num
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(paths[2], 64)
		if err != nil {
			panic(err)
		}
		return v.Float() != f
	default:
		panic(fmt.Sprintf("unsupported data type %s", k))
	}
}

func validateGreaterThan(v reflect.Value, k reflect.Kind, paths []string) bool {
	switch k {
	case reflect.String:
		return v.String() > paths[2]
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Int() > num
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Uint() > num
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(paths[2], 64)
		if err != nil {
			panic(err)
		}
		return v.Float() > f
	default:
		panic(fmt.Sprintf("unsupported data type %s", k))
	}
}

func validateGreaterThanOrEqual(v reflect.Value, k reflect.Kind, paths []string) bool {
	switch k {
	case reflect.String:
		return v.String() >= paths[2]
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Int() >= num
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Uint() >= num
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(paths[2], 64)
		if err != nil {
			panic(err)
		}
		return v.Float() >= f
	default:
		panic(fmt.Sprintf("unsupported data type %s", k))
	}
}

func validateLesserThan(v reflect.Value, k reflect.Kind, paths []string) bool {
	switch k {
	case reflect.String:
		return v.String() < paths[2]
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Int() < num
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Uint() < num
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(paths[2], 64)
		if err != nil {
			panic(err)
		}
		return v.Float() < f
	default:
		panic(fmt.Sprintf("unsupported data type %s", k))
	}
}

func validateLesserThanOrEqual(v reflect.Value, k reflect.Kind, paths []string) bool {
	switch k {
	case reflect.String:
		return v.String() <= paths[2]
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.ParseInt(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Int() <= num
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		num, err := strconv.ParseUint(paths[2], 10, 64)
		if err != nil {
			panic(err)
		}
		return v.Uint() <= num
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(paths[2], 64)
		if err != nil {
			panic(err)
		}
		return v.Float() <= f
	default:
		panic(fmt.Sprintf("unsupported data type %s", k))
	}
}
