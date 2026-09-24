package routes

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// structValidator plugs go-playground/validator into fiber.Config.StructValidator,
// so c.Bind() validates `validate:"..."` tags right after binding.
type structValidator struct {
	validate *validator.Validate
}

func NewStructValidator() fiber.StructValidator {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Report fields by their json/uri name (e.g. "email") instead of the Go field name
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		for _, tag := range []string{"json", "uri"} {
			if name, _, _ := strings.Cut(f.Tag.Get(tag), ","); name != "" && name != "-" {
				return name
			}
		}
		return f.Name
	})
	return &structValidator{validate: v}
}

func (v *structValidator) Validate(out any) error {
	err := v.validate.Struct(out)
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err
	}
	msgs := make([]string, len(verrs))
	for i, fe := range verrs {
		msgs[i] = fmt.Sprintf("%s failed '%s' validation", fe.Field(), fe.Tag())
	}
	return errors.New(strings.Join(msgs, "; "))
}
