package homework

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")
var ErrValidationError = errors.New("validation error")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	stringArray := make([]string, 0)
	for i := 0; i < len(v); i++ {
		stringArray = append(stringArray, v[i].Err.Error())
	}
	return strings.Join(stringArray, ", ")
}

func makeValidationErrors(err []error) ValidationErrors {
	var v ValidationErrors
	for i := 0; i < len(err); i++ {
		v = append(v, ValidationError{err[i]})
	}
	return v
}

func Validate(v any) error {
	if reflect.ValueOf(v).Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var errorStack []error = make([]error, 0)
	dt := reflect.TypeOf(v)
	dv := reflect.ValueOf(v)

	for i := 0; i < dt.NumField(); i++ {

		validation := dt.Field(i).Tag.Get("validate")
		if validation == "" {
			continue
		}
		requirement := strings.Split(validation, ":")
		if len(requirement) < 2 {
			errorStack = append(errorStack, ErrInvalidValidatorSyntax)
			continue
		}

		validator := requirement[0]
		value := requirement[1]
		if value == "" {
			errorStack = append(errorStack, ErrInvalidValidatorSyntax)
			continue
		}
		if !dt.Field(i).IsExported() {
			errorStack = append(errorStack, ErrValidateForUnexportedFields)
			continue
		}

		typeLine := dt.Field(i).Type.Kind()
		var err error = nil
		switch typeLine {
		case reflect.Int:
			err = validateInteger(validator, value, dv.Field(i).Interface().(int))
		case reflect.String:
			err = validateString(validator, value, dv.Field(i).Interface().(string))
		case reflect.Slice:
			if dt.Field(i).Type.String() == "[]string" {
				str := dv.Field(i).Interface().([]string)
				for _, el := range str {
					errorStack = writeDownError(validateString(validator, value, el), errorStack, dt.Field(i).Name+"."+strconv.Itoa(i))
				}
			} else if dt.Field(i).Type.String() == "[]int" {
				str := dv.Field(i).Interface().([]int)
				for _, el := range str {
					errorStack = writeDownError(validateInteger(validator, value, el), errorStack, dt.Field(i).Name+"."+strconv.Itoa(i))
				}
			}
			continue
		}
		errorStack = writeDownError(err, errorStack, dt.Field(i).Name)
	}

	if len(errorStack) == 0 {
		return nil
	} else {
		var e error = makeValidationErrors(errorStack)
		return e
	}
}

func writeDownError(newError error, errorArray []error, name string) []error {
	if errors.Is(newError, ErrInvalidValidatorSyntax) {
		errorArray = append(errorArray, newError)
	} else if errors.Is(newError, ErrValidationError) {
		errorArray = append(errorArray, errors.New("Field "+name+" isn't valid"))
	}
	return errorArray
}

func contains[T comparable](t []T, needle T) bool {
	for _, v := range t {
		if v == needle {
			return true
		}
	}
	return false
}

func checkValue(value int, validator string, f func(int, int) bool) error {
	check, e := strconv.Atoi(validator)
	if e != nil {
		return ErrInvalidValidatorSyntax
	} else if !f(value, check) {
		return ErrValidationError
	}
	return nil
}

func validateInteger(validation string, value string, val int) error {
	switch validation {
	case "min":
		return checkValue(val, value, func(a, b int) bool { return a >= b })
	case "max":
		return checkValue(val, value, func(a, b int) bool { return a <= b })
	case "in":
		arr := strings.Split(value, ",")
		var intArr []int = make([]int, 0)
		for _, el := range arr {
			n, e := strconv.Atoi(el)
			if e != nil {
				return ErrInvalidValidatorSyntax
			} else {
				intArr = append(intArr, n)
			}
		}
		if !contains(intArr, val) {
			return ErrValidationError
		}
	default:
		return ErrInvalidValidatorSyntax
	}
	return nil
}

func validateString(validation string, value string, val string) error {
	switch validation {
	case "len":
		return checkValue(len(val), value, func(a, b int) bool { return a == b })
	case "in":
		if !contains(strings.Split(value, ","), val) {
			return ErrValidationError
		}
	case "max":
		return checkValue(len(val), value, func(a, b int) bool { return a < b })
	case "min":
		return checkValue(len(val), value, func(a, b int) bool { return a > b })
	default:
		return ErrInvalidValidatorSyntax
	}
	return nil
}
