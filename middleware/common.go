package middleware

import "github.com/go-chi/chi/v5/middleware"

var (
	// RequestID is a middleware that injects a request ID into the context of each
	// request. A request ID is a string of the form "host.example.com/random-0001",
	// where "random" is a base62 random string that uniquely identifies this go
	// process, and where the last number is an atomically incremented request
	// counter.
	RequestID = middleware.RequestID

	// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
	// of parsing either the X-Forwarded-For header or the X-Real-IP header (in that
	// order).
	//
	// This middleware should be inserted fairly early in the middleware stack to
	// ensure that subsequent layers (e.g., request loggers) which examine the
	// RemoteAddr will see the intended value.
	//
	// You should only use this middleware if you can trust the headers passed to
	// you (in particular, the two headers this middleware uses), for example
	// because you have placed a reverse proxy like HAProxy or nginx in front of
	// chi. If your reverse proxies are configured to pass along arbitrary header
	// values from the client, or if you use this middleware without a reverse
	// proxy, malicious clients will be able to make you very sad (or, depending on
	// how you're using RemoteAddr, vulnerable to an attack of some sort).
	RealIP = middleware.RealIP

	// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
	// a router (or subrouter) from being cached by an upstream proxy and/or client.
	//
	// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
	//      Expires: Thu, 01 Jan 1970 00:00:00 UTC
	//      Cache-Control: no-cache, private, max-age=0
	//      X-Accel-Expires: 0
	//      Pragma: no-cache (for HTTP/1.0 proxies/clients)
	NoCache = middleware.NoCache
)
