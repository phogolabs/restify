package restify

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/phogolabs/flaw"
	"gopkg.in/go-playground/validator.v9"
)

// Validator validates an object
type Validator struct {
	r *http.Request
	v *validator.Validate
	// stack trace skip
	skip int
}

// NewValidator creates a new validator
func NewValidator(r *http.Request) *Validator {
	validator := &Validator{r: r, v: validator.New(), skip: 1}
	validator.v.RegisterTagNameFunc(validator.name)
	return validator
}

// Validate validates the object
func (v *Validator) Validate(obj interface{}) (err error) {
	err = v.v.StructCtx(v.r.Context(), obj)

	var verrs validator.ValidationErrors

	if errors.As(err, &verrs) {
		errs := flaw.ErrorCollector{}

		for _, field := range verrs {
			werr := fmt.Errorf("'%s' field does not obey '%s' validation rule", field.Namespace(), field.Tag())
			errs.Wrap(werr)
		}

		stack := flaw.NewStackTraceAt(v.skip)
		err = flaw.Wrap(errs, stack...).WithStatus(http.StatusUnprocessableEntity)
	}

	return err
}

func (v *Validator) name(field reflect.StructField) string {
	content, err := ParseContent(v.r.Header.Get("Content-Type"))
	if err != nil {
		return field.Name
	}

	var (
		keys = []string{"query", "header"}
		tag  = field.Name
	)

	switch content.Type {
	case ContentTypeJSON:
		keys = append(keys, "json")
	case ContentTypeXML:
		keys = append(keys, "xml")
	case ContentTypeForm:
		keys = append(keys, "form")
	}

	for index := len(keys) - 1; index >= 0; index-- {
		value := field.Tag.Get(keys[index])

		if value == "" || value == "-" {
			continue
		}

		tag = value
	}

	if idx := strings.Index(tag, ","); idx != -1 {
		tag = tag[:idx]
	}

	return tag
}
