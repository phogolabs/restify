# restify

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][action-img]][action-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

An Opinionated package for building HTTP RESTFul API

## Installation

Make sure you have a working Go environment. Go version 1.13.x is supported.

[See the install instructions for Go](http://golang.org/doc/install.html).

To install `restify`, simply run:

```
$ go get github.com/phogolabs/restify
```

## Getting Started

```golang
type CreateUserInput struct {
	FirstName string `json:"first_name" header:"-"            validate:"required"`
	LastName  string `json:"last_name"  header:"-"            validate:"required"`
	CreatedBy string `json:"-"          header:"X-Created-By" validate:"-"`
}

type CreateUserOutput struct {
	UserID string `json:"user_id"`
}

func (output *CreateUserOutput) Status() int {
	return http.StatusCreated
}

func create(w http.ResponseWriter, r *http.Request) {
	reactor := restify.NewReactor(w, r)

	var (
		input  = &CreateUserInput{}
		output = &CreateUserOutput{}
	)

	if err := reactor.Bind(input); err != nil {
		reactor.Render(err)
		return
	}

	// TODO: implement your logic here

	if err := reactor.Render(output); err != nil {
		reactor.Render(err)
		return
	}
}
```

## Contributing

We are open for any contributions. Just fork the
[project](https://github.com/phogolabs/restify).

[report-img]: https://goreportcard.com/badge/github.com/phogolabs/restify
[report-url]: https://goreportcard.com/report/github.com/phogolabs/restify
[codecov-url]: https://codecov.io/gh/phogolabs/restify
[codecov-img]: https://codecov.io/gh/phogolabs/restify/branch/master/graph/badge.svg
[action-img]: https://github.com/phogolabs/restify/workflows/pipeline/badge.svg
[action-url]: https://github.com/phogolabs/restify/actions
[godoc-url]: https://godoc.org/github.com/phogolabs/restify
[godoc-img]: https://godoc.org/github.com/phogolabs/restify?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
