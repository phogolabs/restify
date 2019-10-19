package restify_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/restify"
)

var _ = Describe("Content", func() {
	Context("when the content-type is text/plain", func() {
		It("parses the content successfully", func() {
			content, err := restify.ParseContent("text/plain")
			Expect(err).To(BeNil())
			Expect(content.Type).To(Equal(restify.ContentTypePlainText))
		})
	})

	Context("when the content-type is text/html", func() {
		It("parses the content successfully", func() {
			content, err := restify.ParseContent("text/html")
			Expect(err).To(BeNil())
			Expect(content.Type).To(Equal(restify.ContentTypeHTML))
		})
	})

	Context("when the content-type is application/json", func() {
		It("parses the content successfully", func() {
			content, err := restify.ParseContent("application/json")
			Expect(err).To(BeNil())
			Expect(content.Type).To(Equal(restify.ContentTypeJSON))
		})
	})

	Context("when the content-type is application/xml", func() {
		It("parses the content successfully", func() {
			content, err := restify.ParseContent("application/xml")
			Expect(err).To(BeNil())
			Expect(content.Type).To(Equal(restify.ContentTypeXML))
		})
	})

	Context("when the content-type is application/x-www-form-urlencoded", func() {
		It("parses the content successfully", func() {
			content, err := restify.ParseContent("application/x-www-form-urlencoded")
			Expect(err).To(BeNil())
			Expect(content.Type).To(Equal(restify.ContentTypeForm))
		})
	})

	Context("when the content-type is text/event-stream", func() {
		It("parses the content successfully", func() {
			content, err := restify.ParseContent("text/event-stream")
			Expect(err).To(BeNil())
			Expect(content.Type).To(Equal(restify.ContentTypeEventStream))
		})
	})

	Context("when the content-type is empty", func() {
		It("parses the content successfully", func() {
			content, err := restify.ParseContent("")
			Expect(err).To(BeNil())
			Expect(content.Type).To(Equal(restify.ContentTypeUnknown))
		})
	})

	Context("when the content-type has invalid parameter", func() {
		It("parses the content successfully", func() {
			content, err := restify.ParseContent("application/json;charset=")
			Expect(err).To(MatchError("mime: invalid media parameter"))
			Expect(content).To(BeNil())
		})
	})

})
