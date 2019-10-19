package restify

import (
	"errors"
	"reflect"

	"github.com/phogolabs/flaw"
	"github.com/phogolabs/log"
)

var (
	rendererType = reflect.TypeOf(new(Renderer)).Elem()
	binderType   = reflect.TypeOf(new(Binder)).Elem()
)

func (r *Reactor) binder(obj interface{}) error {
	value := reflect.ValueOf(obj)
	value = reflect.Indirect(value)

	if value.Kind() == reflect.Struct {
		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)

			if !field.Type().Implements(binderType) {
				continue
			}

			if r.none(field) {
				if !field.CanSet() {
					continue
				}

				r.set(field)
			}

			if child, ok := field.Interface().(Binder); ok {
				if err := r.binder(child); err != nil {
					return err
				}
			}
		}
	}

	if binder, ok := obj.(Binder); ok {
		if err := binder.Bind(r.decoder.request); err != nil {
			return err
		}
	}

	return nil
}

func (r *Reactor) renderer(obj interface{}) error {
	value := reflect.ValueOf(obj)
	value = reflect.Indirect(value)

	if value.Kind() == reflect.Struct {
		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)

			if !field.Type().Implements(rendererType) {
				continue
			}

			if r.none(field) {
				if !field.CanSet() {
					continue
				}

				r.set(field)
			}

			if child, ok := field.Interface().(Renderer); ok {
				if err := r.renderer(child); err != nil {
					return err
				}
			}
		}
	}

	if renderer, ok := obj.(Renderer); ok {
		if err := renderer.Render(r.encoder.response, r.decoder.request); err != nil {
			return err
		}
	}

	return nil
}

func (r *Reactor) none(f reflect.Value) bool {
	switch f.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return f.IsNil()
	default:
		return false
	}
}

func (r *Reactor) set(field reflect.Value) {
	var (
		origin = field.Type()
		target = field.Type()
	)

	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	value := reflect.New(target)

	if origin.Kind() != reflect.Ptr {
		value = value.Elem()
	}

	field.Set(value)
}

func (r *Reactor) error(err error) *flaw.Error {
	var errx *flaw.Error

	if !errors.As(err, &errx) {
		stack := flaw.NewStackTraceAt(2)
		errx = flaw.Wrap(err, stack...)
	}

	return errx
}

func (r *Reactor) log(err error) error {
	logger := log.GetContext(r.encoder.request.Context())
	logger.WithError(err).Errorf("http request %v caused an error",
		r.encoder.request.RequestURI)

	return err
}
