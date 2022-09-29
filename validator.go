package webutils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	qerrors "github.com/cyberhorsey/errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/gommon/log"
)

// Use a single instance of the validator structs to cache the structs
var (
	Validator   *govalidator.Validate
	utranslator *ut.UniversalTranslator
	translator  ut.Translator
)

func init() {
	en := en.New()
	utranslator = ut.New(en, en)
	translator, _ = utranslator.GetTranslator("en")
	Validator = govalidator.New()
	_ = en_translations.RegisterDefaultTranslations(Validator, translator)
	Validator.RegisterTagNameFunc(jsonTagName)
}

func jsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	if name == "-" {
		return ""
	}

	return name
}

// Validate data using the provided validate struct tags per github.com/go-playground/validator
func Validate(data interface{}) []error {
	errs := make([]error, 0)

	err := Validator.Struct(data)
	if err == nil {
		return nil
	}

	for _, fieldErr := range err.(govalidator.ValidationErrors) {
		log.Errorf("validation error: %v", err)

		errs = append(errs, qerrors.Validation.NewWithDetail(ValidationErrorDetail(fieldErr)))
	}

	return errs
}

// ValidationErrorDetail returns the validation error description for fieldErr
func ValidationErrorDetail(fieldErr govalidator.FieldError) string {
	tag := fieldErr.Tag()
	msg := ""

	switch tag {
	case "required":
		msg = fmt.Sprintf("%v is required", fieldErr.Field())
	case "required_without":
		msg = fmt.Sprintf(
			"%v is required unless %v is provided",
			underscore(fieldErr.Field()),
			underscore(fieldErr.Param()),
		)
	case "required_if":
		var requiredIfCondition string

		// required_if=SomeField value SomeOtherField value2
		if params := strings.Split(fieldErr.Param(), " "); len(params) > 1 {
			var requiredIfConditions []string
			for i := 0; i < len(params); i += 2 {
				requiredIfConditions = append(
					requiredIfConditions,
					fmt.Sprintf("%v is %v", underscore(params[i]), params[i+1]),
				)
			}

			requiredIfCondition = strings.Join(requiredIfConditions, ", ")
		} else {
			requiredIfCondition = fieldErr.Param()
		}

		msg = fmt.Sprintf(
			"%v is required if %v",
			underscore(fieldErr.Field()),
			requiredIfCondition,
		)
	default:
		msg = fieldErr.Translate(translator)
	}

	return msg
}

// Underscore a camel case string
// Inspired by: https://play.golang.org/p/z3_83edIpG
func underscore(s string) string {
	// remove all spaces/tabs followed by uppercase character
	spaceRegexp := regexp.MustCompile("[[:space:][:blank:]]([A-Z]+)")
	spaceless := spaceRegexp.ReplaceAll([]byte(s), []byte("$1"))
	s = string(spaceless)

	// convert all spaces/tabs not followed by capital letter to underscore
	spaceRegexp = regexp.MustCompile("[[:space:][:blank:]]([^A-Z]+)")
	spaceless = spaceRegexp.ReplaceAll([]byte(s), []byte("_$1"))
	s = string(spaceless)

	camelRegexp := regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

	var humps []string

	for _, sub := range camelRegexp.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			humps = append(humps, sub[1])
		}

		if sub[2] != "" {
			humps = append(humps, sub[2])
		}
	}

	return strings.ToLower(strings.Join(humps, "_"))
}
