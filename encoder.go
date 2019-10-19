package restify

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/go-playground/form"
)

// Renderer interface for managing response payloads.
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

// Encoder encodes the request
type Encoder struct {
	code     int
	response http.ResponseWriter
	request  *http.Request
	content  *Content
}

// NewEncoder creates a new encoder
func NewEncoder(w http.ResponseWriter, r *http.Request) *Encoder {
	return &Encoder{code: 200, response: w, request: r}
}

// Status sets the response status code
func (e *Encoder) Status(code int) {
	e.code = code
}

// SetContentType overrides the Content-Type attribute
func (e *Encoder) SetContentType(kind ContentType) {
	e.content = &Content{
		Type:   kind,
		Params: make(map[string]string),
	}
}

// Encode encodes the request
func (e *Encoder) Encode(obj interface{}) (err error) {
	content := e.content

	if content == nil {
		content, err = ParseContent(e.header())
		if err != nil {
			return err
		}
	}

	type Response interface {
		Status() int
	}

	if response, ok := obj.(Response); ok {
		if status := response.Status(); status > 0 {
			e.code = status
		}
	}

	// write status code
	e.response.WriteHeader(e.code)

	// encode the body
	switch content.Type {
	case ContentTypeXML:
		err = EncodeXML(e.response, obj)
	case ContentTypeForm:
		err = EncodeForm(e.response, obj)
	default:
		err = EncodeJSON(e.response, obj)
	}

	return err
}

func (e *Encoder) header() string {
	headers := []string{"Accepted", "Content-Type"}

	for _, header := range headers {
		if value := e.request.Header.Get(header); value != "" {
			return value
		}
	}

	return ""
}

// EncodeJSON marshals 'obj' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func EncodeJSON(w http.ResponseWriter, obj interface{}) error {
	var (
		buffer  = &bytes.Buffer{}
		encoder = json.NewEncoder(buffer)
	)

	encoder.SetEscapeHTML(true)

	if err := encoder.Encode(obj); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := buffer.WriteTo(w)
	return err
}

// EncodeXML marshals 'obj' to JSON, setting the Content-Type as application/xml. It
// will automatically prepend a generic XML header (see encoding/xml.Header) if
// one is not found in the first 100 bytes of 'v'.
func EncodeXML(w http.ResponseWriter, obj interface{}) error {
	var (
		buffer  = &bytes.Buffer{}
		encoder = xml.NewEncoder(buffer)
	)

	if err := encoder.Encode(obj); err != nil {
		return err
	}

	if !bytes.HasPrefix(buffer.Bytes(), []byte("<?xml")) {
		fmt.Fprint(w, xml.Header)
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	_, err := buffer.WriteTo(w)
	return err
}

// EncodeForm marshals 'obj' to URL encoded data,  and setting the
// Content-Type as application/x-www-form-urlencoded.
func EncodeForm(w http.ResponseWriter, obj interface{}) error {
	encoder := form.NewEncoder()
	encoder.SetMode(form.ModeExplicit)

	values, err := encoder.Encode(obj)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	_, err = fmt.Fprint(w, values.Encode())
	return err
}
