package utils

import (
	"github.com/go-playground/validator/v10"
)

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return IsSupportedCurrency(currency)
	}
	return false
}
