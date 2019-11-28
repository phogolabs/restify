package restify_test

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/json"
	"github.com/phogolabs/schema"
)

func TestRestify(t *testing.T) {
	log.SetHandler(json.New(GinkgoWriter))

	RegisterFailHandler(Fail)
	RunSpecs(t, "Restify Suite")
}

type Input struct {
	Version int              `json:"-" xml:"-" header:"X-Version"`
	Filter  int              `json:"-" xml:"-" query:"filter"`
	Type    int              `json:"-" xml:"-" path:"type"`
	No      schema.NullInt64 `json:"-" xml:"-" form:"no"`
	Ptr     *int             `json:"-" xml:"-" header:"X-Ptr"`
	Name    string           `json:"json_name" xml:"xml_name" form:"form_name,option" default:"john" validate:"required"`
}

type BindableInput struct {
	Name          string         `json:"json_name" xml:"xml_name" form:"form_name,option" default:"john" validate:"required"`
	BindCnt       int            `json:"-" xml:"-"`
	BindFail      bool           `json:"-" xml:"-"`
	BindableChild *BindableChild `json:"-" xml:"-"`
	SentinelChild SentinelChild  `json:"-" xml:"-"`
}

func (i *BindableInput) Bind(r *http.Request) error {
	i.BindCnt++

	if i.BindFail {
		return fmt.Errorf("oh no")
	}

	return nil
}

type BindableChild struct {
	BindCnt  int  `json:"-" xml:"-"`
	BindFail bool `json:"-" xml:"-"`
}

func (i *BindableChild) Bind(r *http.Request) error {
	i.BindCnt++

	if i.BindFail {
		return fmt.Errorf("oh no")
	}

	return nil
}

type Output struct {
	Name string `json:"json_name" xml:"xml_name" form:"Name"`
}

type RenderableOutput struct {
	Name            string           `json:"json_name" xml:"xml_name" form:"Name"`
	RenderCnt       int              `json:"-" xml:"-"`
	RenderFail      bool             `json:"-" xml:"-"`
	RenderableChild *RenderableChild `json:"-" xml:"-"`
	SentinelChild   SentinelChild    `json:"-" xml:"-"`
}

func (o *RenderableOutput) Render(w http.ResponseWriter, r *http.Request) error {
	o.RenderCnt++

	if o.RenderFail {
		return fmt.Errorf("oh no")
	}

	return nil
}

type RenderableChild struct {
	RenderCnt  int  `json:"-" xml:"-"`
	RenderFail bool `json:"-" xml:"-"`
}

func (o *RenderableChild) Render(w http.ResponseWriter, r *http.Request) error {
	o.RenderCnt++

	if o.RenderFail {
		return fmt.Errorf("oh no")
	}

	return nil
}

type OutputError struct{}

func (e *OutputError) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("oh no")
}

func (e *OutputError) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return fmt.Errorf("oh no")
}

type SentinelChild int

func (s SentinelChild) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s SentinelChild) Bind(r *http.Request) error {
	return nil
}
