package gin

import (
	"github.com/go-playground/validator/v10"

	"github.com/croixxant/go-sandbox/entity"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if c, ok := fieldLevel.Field().Interface().(string); ok {
		return entity.IsSupportedCurrency(c)
	}
	return false
}
