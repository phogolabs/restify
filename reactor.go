package restify

import (
	"net/http"

	"github.com/phogolabs/inflate"
)

// Reactor works with http response and request
type Reactor struct {
	decoder   *Decoder
	encoder   *Encoder
	validator *Validator
}

// NewReactor creates a new reactor
func NewReactor(w http.ResponseWriter, r *http.Request) *Reactor {
	validator := NewValidator(r)
	validator.skip = 2

	reactor := &Reactor{
		decoder:   NewDecoder(r),
		encoder:   NewEncoder(w, r),
		validator: validator,
	}

	return reactor
}

// Status sets the response status code
func (r *Reactor) Status(code int) {
	r.encoder.Status(code)
}

// Bind binds the object
func (r *Reactor) Bind(obj interface{}) (err error) {
	// decode the body
	if err = r.decoder.Decode(obj); err != nil {
		return err
	}

	// set defaults
	if err == nil {
		err = inflate.SetDefault(obj)
	}

	// validate
	if err == nil {
		err = r.validator.Validate(obj)
	}

	if err == nil {
		// bind the input
		err = r.binder(obj)
	}

	return err
}

// Render renders the response
func (r *Reactor) Render(obj interface{}) (err error) {
	// handle error
	if errx, ok := obj.(error); ok {
		obj = r.log(r.error(errx))
	}

	if obj != nil {
		// set defaults
		if err = inflate.SetDefault(obj); err != nil {
			return err
		}
	}

	// encode the body
	if err = r.encoder.Encode(obj); err != nil {
		return err
	}

	if err == nil {
		// render the output
		if renderer, ok := obj.(Renderer); ok {
			err = r.renderer(renderer)
		}
	}

	return err
}

// RenderWith renders the response with code
func (r *Reactor) RenderWith(code int, obj interface{}) (err error) {
	r.Status(code)

	// handle error
	if errx, ok := obj.(error); ok {
		obj = r.error(errx).WithStatus(code)
	}

	return r.Render(obj)
}
