package restify

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/go-chi/chi"
	"github.com/mitchellh/mapstructure"
	"github.com/phogolabs/flaw"
)

// Binder interface for managing request payloads.
type Binder interface {
	Bind(r *http.Request) error
}

// Decoder decodes a request
type Decoder struct {
	request *http.Request
	content *Content
}

// NewDecoder creates a new decoder
func NewDecoder(r *http.Request) *Decoder {
	return &Decoder{request: r}
}

// SetContentType overrides the Content-Type attribute
func (d *Decoder) SetContentType(kind ContentType) {
	d.content = &Content{
		Type:   kind,
		Params: make(map[string]string),
	}
}

// Decode decodes the request
func (d *Decoder) Decode(obj interface{}) (err error) {
	content := d.content

	if content == nil {
		content, err = ParseContent(d.request.Header.Get("Content-Type"))
		if err != nil {
			return err
		}
	}

	switch content.Type {
	case ContentTypeXML:
		err = DecodeXML(d.request.Body, obj)
	case ContentTypeForm:
		err = DecodeForm(d.request.Body, obj)
	default:
		err = DecodeJSON(d.request.Body, obj)
	}

	if errors.Is(err, io.EOF) {
		err = nil
	}

	if err == nil {
		// encode header
		err = NewHeaderDecoder(d.request.Header).Decode(obj)
	}

	if err == nil {
		// encode the url query params
		if d.request.URL != nil {
			err = NewQueryDecoder(d.request.URL.Query()).Decode(obj)
		}
	}

	if err == nil {
		// encode the url params
		if ctx, ok := d.request.Context().Value(chi.RouteCtxKey).(*chi.Context); ok {
			err = NewParamDecoder(ctx.URLParams).Decode(obj)
		}
	}

	return err
}

// MapDecoder decodes a header
type MapDecoder struct {
	data map[string]interface{}
	tag  string
}

// NewHeaderDecoder creates a new header decoder
func NewHeaderDecoder(header http.Header) *MapDecoder {
	decoder := &MapDecoder{
		data: make(map[string]interface{}),
		tag:  "header",
	}

	for key := range header {
		decoder.data[key] = header.Get(key)
	}

	return decoder
}

// NewParamDecoder creates a new param decoder
func NewParamDecoder(params chi.RouteParams) *MapDecoder {
	decoder := &MapDecoder{
		data: make(map[string]interface{}),
		tag:  "param",
	}

	for index, key := range params.Keys {
		decoder.data[key] = params.Values[index]
	}

	return decoder
}

// NewQueryDecoder creates a new query decoder
func NewQueryDecoder(values url.Values) *MapDecoder {
	decoder := &MapDecoder{
		data: make(map[string]interface{}),
		tag:  "query",
	}

	for key := range values {
		decoder.data[key] = values.Get(key)
	}

	return decoder
}

// NewFormDecoder creates a new query decoder
func NewFormDecoder(values url.Values) *MapDecoder {
	decoder := &MapDecoder{
		data: make(map[string]interface{}),
		tag:  "form",
	}

	for key := range values {
		decoder.data[key] = values.Get(key)
	}

	return decoder
}

// Decode decodes the object from a map
func (d *MapDecoder) Decode(obj interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           obj,
		DecodeHook:       d.convert,
		TagName:          d.tag,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return d.error(err)
	}

	err = decoder.Decode(d.data)
	return d.error(err)
}

func (d *MapDecoder) convert(source reflect.Type, target reflect.Type, data interface{}) (interface{}, error) {
	element := target

	if target.Kind() == reflect.Ptr {
		element = target.Elem()
	}

	value := reflect.New(element)

	if value.CanInterface() {
		if scanner, ok := value.Interface().(sql.Scanner); ok {
			if err := scanner.Scan(data); err != nil {
				return nil, err
			}

			return value.Interface(), nil
		}
	}

	return data, nil
}

func (d *MapDecoder) error(err error) error {
	type ErrorList interface {
		WrappedErrors() []error
	}

	if list, ok := err.(ErrorList); ok {
		errs := flaw.ErrorCollector{}

		for _, errx := range list.WrappedErrors() {
			errs.Wrap(errx)
		}

		return errs
	}

	return err
}

// DecodeJSON decodes a object from JSON reader
func DecodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return json.NewDecoder(r).Decode(v)
}

// DecodeXML decodes a object from XML reader
func DecodeXML(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return xml.NewDecoder(r).Decode(v)
}

// DecodeForm decodes a object from Form reader
func DecodeForm(r io.ReadCloser, v interface{}) error {
	// 10 MB limit
	max := int64(10 << 20)

	r = http.MaxBytesReader(nil, r, max)

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	values, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	decoder := NewFormDecoder(values)
	return decoder.Decode(v)
}
