package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type customValidator struct {
	validate *validator.Validate
}

func NewCustomValidator(validate *validator.Validate) Validator {
	validate = validator.New()

	_ = validate.RegisterValidation("uuid", validateUUID)

	return &customValidator{validate: validate}
}

type Validator interface {
	ValidateStruct(interface{}) error
}

func (v *customValidator) ValidateStruct(s interface{}) error {
	err := v.validate.Struct(s)
	if err != nil {
		return errors.New(extractErrorMessage(err.Error()))
	}

	return nil
}

func extractErrorMessage(fullErrorMsg string) string {
	parts := strings.Split(fullErrorMsg, "Error:")
	if len(parts) > 1 {
		return strings.TrimSpace(strings.ToLower(parts[1]))
	}
	return fullErrorMsg
}

func validateUUID(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	_, err := uuid.Parse(str)
	return err == nil
}
