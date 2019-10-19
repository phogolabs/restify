package restify_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing/iotest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/flaw"
	"github.com/phogolabs/restify"
	"github.com/phogolabs/restify/fake"
)

var _ = Describe("Reactor", func() {
	var (
		reactor  *restify.Reactor
		request  *http.Request
		response *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		buffer := &bytes.Buffer{}

		request = httptest.NewRequest("POST", "http://example.com", buffer)
		response = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		reactor = restify.NewReactor(response, request)
	})

	Describe("Status", func() {
		It("sets the status successfully", func() {
			reactor.Status(http.StatusAccepted)

			Expect(reactor.Render(nil)).To(Succeed())
			Expect(response.Code).To(Equal(http.StatusAccepted))
		})
	})

	Describe("Bind", func() {
		It("binds the input successfully", func() {
			input := &BindableInput{}
			Expect(reactor.Bind(input)).To(Succeed())
			Expect(input.BindCnt).To(Equal(1))
		})

		Context("when the binding fails", func() {
			It("returns an error", func() {
				input := &BindableInput{BindFail: true}
				Expect(reactor.Bind(input)).To(MatchError("oh no"))
				Expect(input.BindCnt).To(Equal(1))
				Expect(input.BindableChild).NotTo(BeNil())
				Expect(input.BindableChild.BindCnt).To(Equal(1))
			})

			Context("when the child fails", func() {
				It("returns an error", func() {
					input := &BindableInput{
						BindableChild: &BindableChild{
							BindFail: true,
						},
					}

					Expect(reactor.Bind(input)).To(MatchError("oh no"))
					Expect(input.BindableChild).NotTo(BeNil())
					Expect(input.BindableChild.BindCnt).To(Equal(1))
				})
			})
		})

		Context("when the request boy returns an error", func() {
			JustBeforeEach(func() {
				var (
					data   = make([]byte, 1)
					reader = iotest.TimeoutReader(&bytes.Buffer{})
				)

				reader.Read(data)
				request.Body = ioutil.NopCloser(reader)
			})

			It("returns an error", func() {
				input := &BindableInput{}
				Expect(reactor.Bind(input)).To(MatchError("timeout"))
			})
		})
	})

	Describe("Render", func() {
		It("renders the output successfully", func() {
			output := &RenderableOutput{}
			Expect(reactor.Render(output)).To(Succeed())
			Expect(output.RenderCnt).To(Equal(1))
			Expect(output.RenderableChild).NotTo(BeNil())
			Expect(output.RenderableChild.RenderCnt).To(Equal(1))
		})

		Context("when the rendering fails", func() {
			It("returns an error", func() {
				output := &RenderableOutput{RenderFail: true}
				Expect(reactor.Render(output)).To(MatchError("oh no"))
				Expect(output.RenderCnt).To(Equal(1))
			})

			Context("when the child fails", func() {
				It("returns an error", func() {
					output := &RenderableOutput{
						RenderableChild: &RenderableChild{
							RenderFail: true,
						},
					}

					Expect(reactor.Render(output)).To(MatchError("oh no"))
					Expect(output.RenderableChild.RenderCnt).To(Equal(1))
				})
			})
		})

		Context("when the output is error", func() {
			It("renders the output successfully", func() {
				output := fmt.Errorf("oh no")
				Expect(reactor.Render(output)).To(Succeed())
				Expect(response.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("when the output is not a pointer", func() {
			It("returns an error", func() {
				output := RenderableOutput{}
				Expect(reactor.Render(output)).To(MatchError("not a struct pointer"))
			})
		})

		Context("when the response writer returns an error", func() {
			JustBeforeEach(func() {
				response := &fake.ResponseWriter{}
				response.HeaderReturns(http.Header{})
				response.WriteReturns(0, fmt.Errorf("oh no"))

				reactor = restify.NewReactor(response, request)
			})

			It("returns an error", func() {
				output := &RenderableOutput{}
				Expect(reactor.Render(output)).To(MatchError("oh no"))
			})
		})
	})

	Describe("RenderWith", func() {
		It("renders the output with the given code successfully", func() {
			output := &RenderableOutput{}
			Expect(reactor.RenderWith(http.StatusCreated, output)).To(Succeed())
			Expect(response.Code).To(Equal(http.StatusCreated))
		})

		Context("when the output is an error", func() {
			It("overrides the error status", func() {
				output := flaw.New("not found")
				Expect(reactor.RenderWith(http.StatusNotFound, output)).To(Succeed())
				Expect(response.Code).To(Equal(http.StatusNotFound))
			})
		})
	})
})
