package restify_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing/iotest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-chi/chi"
	"github.com/phogolabs/restify"
)

var _ = Describe("Decoder", func() {
	var (
		request *http.Request
		decoder *restify.Decoder
		buffer  *bytes.Buffer
	)

	BeforeEach(func() {
		buffer = &bytes.Buffer{}
		request = httptest.NewRequest("POST", "http://example.com", buffer)
	})

	JustBeforeEach(func() {
		decoder = restify.NewDecoder(request)
	})

	Describe("SetContentType", func() {
		It("sets the content-type successfully", func() {
			decoder.SetContentType(restify.ContentTypeJSON)
			buffer.WriteString("{\"json_name\":\"john\"}")

			input := &Input{}
			Expect(decoder.Decode(input)).To(Succeed())
			Expect(input.Name).To(Equal("john"))
		})
	})

	Describe("Decode", func() {
		Context("when the content-type is application/json", func() {
			BeforeEach(func() {
				request.Header.Set("Content-Type", "application/json")
			})

			It("decodes the input successfully", func() {
				buffer.WriteString("{\"json_name\":\"john\"}")

				input := &Input{}
				Expect(decoder.Decode(input)).To(Succeed())
				Expect(input.Name).To(Equal("john"))
			})

			Context("when the body is empty", func() {
				It("decodes the input successfully", func() {
					input := &Input{}
					Expect(decoder.Decode(input)).To(Succeed())
				})
			})
		})

		Context("when the content-type is application/xml", func() {
			BeforeEach(func() {
				buffer.WriteString(xml.Header)
				buffer.WriteString("<Input><xml_name>peter</xml_name></Input>")
				request.Header.Set("Content-Type", "application/xml")
			})

			It("decodes the input successfully", func() {
				input := &Input{}
				Expect(decoder.Decode(input)).To(Succeed())
				Expect(input.Name).To(Equal("peter"))
			})
		})
	})

	Context("when the content-type is application/x-www-form-urlencoded", func() {
		BeforeEach(func() {
			buffer.WriteString("no=1")
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		})

		It("decodes the input successfully", func() {
			input := &Input{}
			Expect(decoder.Decode(input)).To(Succeed())
			Expect(input.No.Int64).To(Equal(int64(1)))
		})

		Context("when the value cannot be decoded", func() {
			BeforeEach(func() {
				buffer.Reset()
				buffer.WriteString("no=London")
			})

			It("returns an error", func() {
				input := &Input{}
				Expect(decoder.Decode(input)).To(MatchError("cannot convert string 'London' to struct: converting driver.Value type string (\"London\") to a int64: invalid syntax"))
			})
		})

		Context("when the buffer returns an error", func() {
			BeforeEach(func() {
				buffer.Reset()

				var (
					data   = make([]byte, 1)
					reader = iotest.TimeoutReader(buffer)
				)

				reader.Read(data)
				request.Body = ioutil.NopCloser(reader)
			})

			It("returns an error", func() {
				input := &Input{}
				Expect(decoder.Decode(input)).To(MatchError("timeout"))
			})
		})

		Context("when the body cannot be parsed", func() {
			BeforeEach(func() {
				buffer.Reset()
				buffer.WriteString("%no")
			})

			It("returns an error", func() {
				input := &Input{}
				Expect(decoder.Decode(input)).To(MatchError("invalid URL escape \"%no\""))
			})
		})

		Context("when the value is not a pointer", func() {
			It("returns an error", func() {
				input := Input{}
				Expect(decoder.Decode(input)).To(MatchError("the target must be a pointer"))
			})
		})
	})

	Context("when the content-type is a unknown", func() {
		BeforeEach(func() {
			request.Header.Set("Content-Type", "linux/unknown;param=")
		})

		It("returns an error", func() {
			input := &Input{}
			Expect(decoder.Decode(input)).To(MatchError("mime: invalid media parameter"))
		})
	})

	Context("when the header is decoded", func() {
		BeforeEach(func() {
			buffer.WriteString("{}")
			request.Header.Set("X-Version", "1")
			request.Header.Set("X-Ptr", "2")
		})

		It("decodes the input successfully", func() {
			input := &Input{}
			Expect(decoder.Decode(input)).To(Succeed())
			Expect(input.Version).To(Equal(1))
		})

		Context("when the value cannot be decoded", func() {
			BeforeEach(func() {
				request.Header.Set("X-Version", "alpha")
			})

			It("returns an error", func() {
				input := &Input{}
				Expect(decoder.Decode(input)).To(MatchError("cannot convert string 'alpha' to int: strconv.ParseInt: parsing \"alpha\": invalid syntax"))
			})
		})
	})

	Context("when the query is decoded", func() {
		BeforeEach(func() {
			buffer.WriteString("{}")
			request.URL.RawQuery = "filter=1"
		})

		It("decodes the input successfully", func() {
			input := &Input{}
			Expect(decoder.Decode(input)).To(Succeed())
			Expect(input.Filter).To(Equal(1))
		})

		Context("when the value cannot be decoded", func() {
			BeforeEach(func() {
				request.URL.RawQuery = "filter=john"
			})

			It("returns an error", func() {
				input := &Input{}
				Expect(decoder.Decode(input)).To(MatchError("cannot convert string 'john' to int: strconv.ParseInt: parsing \"john\": invalid syntax"))
			})
		})
	})

	Context("when the param is decoded", func() {
		BeforeEach(func() {
			buffer.WriteString("{}")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Keys = []string{"type"}
			rctx.URLParams.Values = []string{"1"}

			ctx := request.Context()
			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

			request = request.WithContext(ctx)
		})

		It("decodes the input successfully", func() {
			input := &Input{}
			Expect(decoder.Decode(input)).To(Succeed())
			Expect(input.Type).To(Equal(1))
		})

		Context("when the value cannot be decoded", func() {
			BeforeEach(func() {
				buffer.WriteString("{}")

				rctx := chi.NewRouteContext()
				rctx.URLParams.Keys = []string{"type"}
				rctx.URLParams.Values = []string{"employee"}

				ctx := request.Context()
				ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

				request = request.WithContext(ctx)
			})

			It("returns an error", func() {
				input := &Input{}
				Expect(decoder.Decode(input)).To(MatchError("cannot convert string 'employee' to int: strconv.ParseInt: parsing \"employee\": invalid syntax"))
			})
		})
	})
})
