package validatorx

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var roomCodeRegex = regexp.MustCompile(`^[A-Z0-9]{5}$`)

// RegisterCustomValidators 注册项目级自定义参数校验器。
func RegisterCustomValidators() error {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil
	}

	return engine.RegisterValidation("room_code", func(fl validator.FieldLevel) bool {
		return roomCodeRegex.MatchString(fl.Field().String())
	})
}

