package restify

import (
	"mime"
)

// ContentType is an enumeration of common HTTP content types.
type ContentType int

// ContentTypes handled by this package.
const (
	ContentTypeUnknown ContentType = iota
	ContentTypePlainText
	ContentTypeHTML
	ContentTypeJSON
	ContentTypeXML
	ContentTypeForm
	ContentTypeEventStream
)

// Content represents the content type
type Content struct {
	Type   ContentType
	Params map[string]string
}

// ParseContent parse the content type
func ParseContent(header string) (*Content, error) {
	content := &Content{
		Type:   ContentTypeUnknown,
		Params: make(map[string]string),
	}

	if header == "" {
		return content, nil
	}

	media, params, err := mime.ParseMediaType(header)

	if err != nil {
		return nil, err
	}

	content.Params = params

	switch media {
	case "text/plain":
		content.Type = ContentTypePlainText
	case "text/html", "application/xhtml+xml":
		content.Type = ContentTypeHTML
	case "application/json", "text/javascript":
		content.Type = ContentTypeJSON
	case "text/xml", "application/xml":
		content.Type = ContentTypeXML
	case "application/x-www-form-urlencoded":
		content.Type = ContentTypeForm
	case "text/event-stream":
		content.Type = ContentTypeEventStream
	}

	return content, nil
}
