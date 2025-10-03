package validators

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

var SliceValidation = map[string][]string{
	"roles":           {"admin", "teacher", "student"},
	"grantValid":      {"code", "token", "code token"},
	"responsesValid":  {"code", "token", "code token", "implicit"},
	"nameAppValid":    {"web app", "mobil app", "desktop app"},
	"authMethodValid": {"client_secret_basic", "client_secret_post", "none", "private_key_jwt"},
}

type NotFoundError struct {
	Msg string
}

func (fount NotFoundError) Error() string {
	return fmt.Sprintf("NotFoundError:%s", fount.Msg)
}

func init() {
	Validate = validator.New()

	Validate.RegisterValidation("password", PasswordValidator)
	Validate.RegisterValidation("name", NameValidator)
	Validate.RegisterValidation("rowallowed", InSliceValidator(SliceValidation["roles"]))
	Validate.RegisterValidation("grantallowed", InSliceValidator(SliceValidation["grantValid"]))
	Validate.RegisterValidation("urlallowed", URLArrayValidator)
	Validate.RegisterValidation("responseallowed", ResponseValidator(SliceValidation["responsesValid"]))
	Validate.RegisterValidation("authmethodallowed", ResponseValidator(SliceValidation["authMethodValid"]))
	Validate.RegisterValidation("appallowed", ResponseValidator(SliceValidation["nameAppValid"]))
}

func PasswordValidator(fl validator.FieldLevel) bool {
	p := fl.Field().String()
	if len(p) < 8 {
		return false
	}
	return regexp.MustCompile(`[a-z]`).MatchString(p) &&
		regexp.MustCompile(`[A-Z]`).MatchString(p) &&
		regexp.MustCompile(`[0-9]`).MatchString(p) &&
		regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(p)
}

func NameValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if len(value) < 4 || len(value) > 50 {
		return false
	}

	return !strings.Contains(value, " ")
}

func InSliceValidator(slice []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if len(value) < 4 || len(value) > 50 {
			return false
		}
		return slices.Contains(slice, value)
	}
}

func ResponseValidator(slice []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() != reflect.Slice {
			return false
		}
		for i := 0; i < field.Len(); i++ {
			elem := field.Index(i).String()
			if !slices.Contains(slice, elem) {
				return false
			}
		}

		return true
	}
}

func URLArrayValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i).String()
		if _, err := url.ParseRequestURI(elem); err != nil {
			return false
		}
	}

	return true
}

func ValidateStruct[T any](s T) (err error) {
	return Validate.Struct(s)
}
