# restify

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][travis-img]][travis-url]
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
	spew.Dump(input)

	if err := reactor.RenderWith(http.StatusCreated, output); err != nil {
		reactor.Render(err)
		return
	}
}
```

## Contributing

We are welcome to any contributions. Just fork the
[project](https://github.com/phogolabs/restify).

[travis-img]: https://travis-ci.org/phogolabs/restify.svg?branch=master
[travis-url]: https://travis-ci.org/phogolabs/restify
[report-img]: https://goreportcard.com/badge/github.com/phogolabs/restify
[report-url]: https://goreportcard.com/report/github.com/phogolabs/restify
[codecov-url]: https://codecov.io/gh/phogolabs/restify
[codecov-img]: https://codecov.io/gh/phogolabs/restify/branch/master/graph/badge.svg
[godoc-url]: https://godoc.org/github.com/phogolabs/restify
[godoc-img]: https://godoc.org/github.com/phogolabs/restify?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
