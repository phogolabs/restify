package restify

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/phogolabs/inflate"
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
		err = inflate.NewHeaderDecoder(d.request.Header).Decode(obj)
	}

	if err == nil {
		// encode cookies
		err = inflate.NewCookieDecoder(d.request.Cookies()).Decode(obj)
	}

	if err == nil {
		// encode the url query params
		if d.request.URL != nil {
			err = inflate.NewQueryDecoder(d.request.URL.Query()).Decode(obj)
		}
	}

	if err == nil {
		// encode the url path params
		if ctx, ok := d.request.Context().Value(chi.RouteCtxKey).(*chi.Context); ok {
			err = inflate.NewPathDecoder(&ctx.URLParams).Decode(obj)
		}
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
	//	10 MB limit
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

	decoder := inflate.NewFormDecoder(values)
	return decoder.Decode(v)
}
