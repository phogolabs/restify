package restify_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw"
	"github.com/phogolabs/restify"
)

var _ = Describe("Validator", func() {
	var (
		request   *http.Request
		validator *restify.Validator
	)

	BeforeEach(func() {
		request = httptest.NewRequest("POST", "http://example.com", nil)
		validator = restify.NewValidator(request)
	})

	It("validates an object successfully", func() {
		input := &Input{}
		Expect(validator.Validate(input)).To(MatchError("message: 'Input' validation failed cause: ['Name' field does not obey 'required' validation rule]"))
	})

	Context("when the content-type is unknown", func() {
		BeforeEach(func() {
			request.Header.Set("Content-Type", "application/unknown")
		})

		It("validates an object successfully", func() {
			input := &Input{}
			Expect(validator.Validate(input)).To(MatchError("message: 'Input' validation failed cause: ['Name' field does not obey 'required' validation rule]"))
		})
	})

	Context("when the content-type is application/json", func() {
		BeforeEach(func() {
			request.Header.Set("Content-Type", "application/json")
		})

		It("validates an object successfully", func() {
			var (
				input = &Input{}
				err   = validator.Validate(input)
			)

			Expect(flaw.Status(err)).To(Equal(http.StatusUnprocessableEntity))
			Expect(err).To(MatchError("message: 'Input' validation failed cause: ['json_name' field does not obey 'required' validation rule]"))
		})
	})

	Context("when the content-type is application/xml", func() {
		BeforeEach(func() {
			request.Header.Set("Content-Type", "application/xml")
		})

		It("validates an object successfully", func() {
			var (
				input = &Input{}
				err   = validator.Validate(input)
			)

			Expect(flaw.Status(err)).To(Equal(http.StatusUnprocessableEntity))
			Expect(err).To(MatchError("message: 'Input' validation failed cause: ['xml_name' field does not obey 'required' validation rule]"))
		})
	})

	Context("when the content-type is application/x-www-form-urlencoded", func() {
		BeforeEach(func() {
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		})

		It("validates an object successfully", func() {
			var (
				input = &Input{}
				err   = validator.Validate(input)
			)

			Expect(flaw.Status(err)).To(Equal(http.StatusUnprocessableEntity))
			Expect(err).To(MatchError("message: 'Input' validation failed cause: ['form_name' field does not obey 'required' validation rule]"))
		})
	})

	Context("when the content-type has wrong parameter", func() {
		BeforeEach(func() {
			request.Header.Set("Content-Type", "application/json;charset=")
		})

		It("validates an object successfully", func() {
			var (
				input = &Input{}
				err   = validator.Validate(input)
			)

			Expect(flaw.Status(err)).To(Equal(http.StatusUnprocessableEntity))
			Expect(err).To(MatchError("message: 'Input' validation failed cause: ['Name' field does not obey 'required' validation rule]"))
		})
	})
})
