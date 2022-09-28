package restify_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw"
	"github.com/phogolabs/restify"
	"github.com/phogolabs/restify/fake"
)

var _ = Describe("Encoder", func() {
	var (
		request  *http.Request
		response *httptest.ResponseRecorder
		encoder  *restify.Encoder
	)

	BeforeEach(func() {
		request = httptest.NewRequest("POST", "http://example.com", nil)
		response = httptest.NewRecorder()
		encoder = restify.NewEncoder(response, request)
	})

	Describe("Status", func() {
		It("sets the status successfully", func() {
			encoder.Status(http.StatusNoContent)

			Expect(encoder.Encode(nil)).To(Succeed())
			Expect(response.Code).To(Equal(http.StatusNoContent))
			Expect(response.Body.String()).To(Equal("null\n"))
		})
	})

	Describe("Encode", func() {
		Context("when the accepted content-type is application/json", func() {
			BeforeEach(func() {
				request.Header.Set("Accepted", "application/json")
			})

			It("encodes the output successfully", func() {
				output := &Output{Name: "John"}
				Expect(encoder.Encode(output)).To(Succeed())
				Expect(response.Code).To(Equal(http.StatusOK))
				Expect(response.Body.String()).To(Equal("{\"json_name\":\"John\"}\n"))
			})

			Context("when the output is an error", func() {
				It("encodes the output successfully", func() {
					Expect(encoder.Encode(flaw.Errorf("oh no"))).To(Succeed())
					Expect(response.Code).To(Equal(http.StatusInternalServerError))
					Expect(response.Body.String()).To(Equal("{\"error_message\":\"oh no\"}\n"))
				})
			})

			Context("when the response returns an error", func() {
				BeforeEach(func() {
					response := &fake.ResponseWriter{}
					response.HeaderReturns(http.Header{})
					response.WriteReturns(0, fmt.Errorf("oh no"))

					encoder = restify.NewEncoder(response, request)
				})

				It("returns an error", func() {
					output := &Output{Name: "John"}
					Expect(encoder.Encode(output)).To(MatchError("oh no"))
				})
			})

			Context("when the marshalling fails", func() {
				It("returns an error", func() {
					output := &OutputError{}
					Expect(encoder.Encode(output)).To(MatchError("json: error calling MarshalJSON for type *restify_test.OutputError: oh no"))
				})
			})
		})

		Context("when the accepted content-type is application/xml", func() {
			BeforeEach(func() {
				request.Header.Set("Accepted", "application/xml")
			})

			It("encodes the output successfully", func() {
				output := &Output{Name: "John"}
				Expect(encoder.Encode(output)).To(Succeed())
				Expect(response.Code).To(Equal(http.StatusOK))
				Expect(response.Body.String()).To(Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Output><xml_name>John</xml_name></Output>"))
			})

			Context("when the output is an error", func() {
				It("encodes the output successfully", func() {
					Expect(encoder.Encode(flaw.Errorf("oh no"))).To(Succeed())
					Expect(response.Code).To(Equal(http.StatusInternalServerError))
					Expect(response.Body.String()).To(Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Error><ErrorMessage>oh no</ErrorMessage></Error>"))
				})
			})

			Context("when the response returns an error", func() {
				BeforeEach(func() {
					response := &fake.ResponseWriter{}
					response.HeaderReturns(http.Header{})
					response.WriteReturns(0, fmt.Errorf("oh no"))

					encoder = restify.NewEncoder(response, request)
				})

				It("returns an error", func() {
					output := &Output{Name: "John"}
					Expect(encoder.Encode(output)).To(MatchError("oh no"))
				})
			})

			Context("when the marshalling fails", func() {
				It("returns an error", func() {
					output := &OutputError{}
					Expect(encoder.Encode(output)).To(MatchError("oh no"))
				})
			})
		})

		Context("when the accepted content-type is application/x-www-form-urlencoded", func() {
			BeforeEach(func() {
				request.Header.Set("Accepted", "application/x-www-form-urlencoded")
			})

			It("encodes the output successfully", func() {
				output := &Output{Name: "John"}
				Expect(encoder.Encode(output)).To(Succeed())
				Expect(response.Code).To(Equal(http.StatusOK))
				Expect(response.Body.String()).To(Equal("Name=John"))
			})

			Context("when the value is a nil pointer", func() {
				It("returns an error", func() {
					var output *Output
					Expect(encoder.Encode(output)).To(MatchError("form: Encode(nil *restify_test.Output)"))
				})
			})
		})

		Context("when the accepted content-type has wrong parameter", func() {
			BeforeEach(func() {
				request.Header.Set("Accepted", "linux/unknown;param=")
			})

			It("returns an error", func() {
				output := &Output{Name: "John"}
				Expect(encoder.Encode(output)).To(MatchError("mime: invalid media parameter"))
			})
		})

		Context("when the accepted content-type is not provided", func() {
			It("encodes the output successfully", func() {
				output := &Output{Name: "John"}
				Expect(encoder.Encode(output)).To(Succeed())
				Expect(response.Code).To(Equal(http.StatusOK))
				Expect(response.Body.String()).To(Equal("{\"json_name\":\"John\"}\n"))
			})
		})
	})

	Describe("SetContentType", func() {
		It("encodes the output successfully", func() {
			encoder.SetContentType(restify.ContentTypeXML)

			Expect(encoder.Encode(flaw.Errorf("oh no"))).To(Succeed())
			Expect(response.Code).To(Equal(http.StatusInternalServerError))
			Expect(response.Body.String()).To(Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<Error><ErrorMessage>oh no</ErrorMessage></Error>"))
		})
	})
})
