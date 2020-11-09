package auth

import (
	"net/http"
)

/* REQUEST */

// An HTTP request with the User. It compose with standard http.Request, so
// you can access with all method and field simply.
type Request struct {
	http.Request
	// The User can be nil if the handler UserLevel is set to LevelNo, else it
	// is never nil.
	User *User
}

// The HTTP request to use when you need standard http.Request.
func (r *Request) R() *http.Request {
	return &r.Request
}

/* HANDLERS */

// Like http.handler but with a auth.Request.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *Request)
}

type HandlerFunc func(http.ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *Request) {
	f(w, r)
}

/* ERRORRESPONSE */

// Variant of http.Error function.
type ErrorResponse func(http.ResponseWriter, *http.Request, error, int)

// A simple default ErrorResponse, it's a binding of http.Error.
func ErrorResponseDefault(w http.ResponseWriter, _ *http.Request, err error, code int) {
	http.Error(w, err.Error(), code)
}
