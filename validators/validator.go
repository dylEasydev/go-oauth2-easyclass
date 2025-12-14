package validators

import (
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

// slice de constante de validation
// on puvait créer un fichier json pour le stocké
var SliceValidation = map[string][]string{
	"roles":           {"admin", "teacher", "student"},
	"tableName":       {"user", "teacher_temp", "student_temps"},
	"grantValid":      {"code", "token", "code token", "client_credentials", "password"},
	"responsesValid":  {"code", "token", "code token", "implicit"},
	"nameAppValid":    {"web app", "mobil app", "desktop app"},
	"authMethodValid": {"client_secret_basic", "client_secret_post", "none", "private_key_jwt"},
}

// initialisation des tags du validateur V10
func init() {
	Validate = validator.New()

	Validate.RegisterValidation("password", PasswordValidator)
	Validate.RegisterValidation("name", NameValidator)
	Validate.RegisterValidation("rowallowed", InSliceValidator(SliceValidation["roles"]))
	Validate.RegisterValidation("grantallowed", ResponseValidator(SliceValidation["grantValid"]))
	Validate.RegisterValidation("urlallowed", URLArrayValidator)
	Validate.RegisterValidation("responseallowed", ResponseValidator(SliceValidation["responsesValid"]))
	Validate.RegisterValidation("authmethodallowed", InSliceValidator(SliceValidation["authMethodValid"]))
	Validate.RegisterValidation("appallowed", InSliceValidator(SliceValidation["nameAppValid"]))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", PasswordValidator)
		v.RegisterValidation("name", NameValidator)
		v.RegisterValidation("rowallowed", InSliceValidator(SliceValidation["roles"]))
		v.RegisterValidation("grantallowed", ResponseValidator(SliceValidation["grantValid"]))
		v.RegisterValidation("urlallowed", URLArrayValidator)
		v.RegisterValidation("tableName", InSliceValidator(SliceValidation["tableName"]))
		v.RegisterValidation("responseallowed", ResponseValidator(SliceValidation["responsesValid"]))
		v.RegisterValidation("authmethodallowed", InSliceValidator(SliceValidation["authMethodValid"]))
		v.RegisterValidation("appallowed", InSliceValidator(SliceValidation["nameAppValid"]))
	}
}

// validation de password ( lettre miniscule , lettre majuscule , carractère numérique et elpha numerique )
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

// validation des champs string ( pas d'espace et une longueur de min 4 et au plus 50 )
func NameValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if len(value) < 4 || len(value) > 50 {
		return false
	}

	return !strings.Contains(value, " ")
}

// valide si le champs est contenus dans la liste des slices fournis en paramètre
func InSliceValidator(slice []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return slices.Contains(slice, value)
	}
}

// valide si toute les donné du chanps []strings sont contenu dans le slice fornis en parmètre
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

// verifie si un tableau URL sont des URl valide
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

// fonction de la validation de structure
func ValidateStruct[T any](s T) (err error) {
	return Validate.Struct(s)
}
