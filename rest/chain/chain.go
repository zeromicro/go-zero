package chain

// This is a modified version of https://github.com/justinas/alice
// The original code is licensed under the MIT license.
// It's modified for couple reasons:
// - Added the Chain interface
// - Added support for the Chain.Prepend(...) method

import "net/http"

type (
	// Chain defines a chain of middleware.
	Chain interface {
		Append(middlewares ...Middleware) Chain
		Prepend(middlewares ...Middleware) Chain
		Then(h http.Handler) http.Handler
		ThenFunc(fn http.HandlerFunc) http.Handler
	}

	// Middleware is an HTTP middleware.
	Middleware func(http.Handler) http.Handler

	// chain acts as a list of http.Handler middlewares.
	// chain is effectively immutable:
	// once created, it will always hold
	// the same set of middlewares in the same order.
	chain struct {
		middlewares []Middleware
	}
)

// New creates a new Chain, memorizing the given list of middleware middlewares.
// New serves no other function, middlewares are only called upon a call to Then() or ThenFunc().
func New(middlewares ...Middleware) Chain {
	return chain{middlewares: append(([]Middleware)(nil), middlewares...)}
}

// Append extends a chain, adding the specified middlewares as the last ones in the request flow.
//
//	c := chain.New(m1, m2)
//	c.Append(m3, m4)
//	// requests in c go m1 -> m2 -> m3 -> m4
func (c chain) Append(middlewares ...Middleware) Chain {
	return chain{middlewares: join(c.middlewares, middlewares)}
}

// Prepend extends a chain by adding the specified chain as the first one in the request flow.
//
//	c := chain.New(m3, m4)
//	c1 := chain.New(m1, m2)
//	c.Prepend(c1)
//	// requests in c go m1 -> m2 -> m3 -> m4
func (c chain) Prepend(middlewares ...Middleware) Chain {
	return chain{middlewares: join(middlewares, c.middlewares)}
}

// Then chains the middleware and returns the final http.Handler.
//
//	New(m1, m2, m3).Then(h)
//
// is equivalent to:
//
//	m1(m2(m3(h)))
//
// When the request comes in, it will be passed to m1, then m2, then m3
// and finally, the given handler
// (assuming every middleware calls the following one).
//
// A chain can be safely reused by calling Then() several times.
//
//	stdStack := chain.New(ratelimitHandler, csrfHandler)
//	indexPipe = stdStack.Then(indexHandler)
//	authPipe = stdStack.Then(authHandler)
//
// Note that middlewares are called on every call to Then() or ThenFunc()
// and thus several instances of the same middleware will be created
// when a chain is reused in this way.
// For proper middleware, this should cause no problems.
//
// Then() treats nil as http.DefaultServeMux.
func (c chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}

	return h
}

// ThenFunc works identically to Then, but takes
// a HandlerFunc instead of a Handler.
//
// The following two statements are equivalent:
//
//	c.Then(http.HandlerFunc(fn))
//	c.ThenFunc(fn)
//
// ThenFunc provides all the guarantees of Then.
func (c chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	// This nil check cannot be removed due to the "nil is not nil" common mistake in Go.
	// Required due to: https://stackoverflow.com/questions/33426977/how-to-golang-check-a-variable-is-nil
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fn)
}

func join(a, b []Middleware) []Middleware {
	mids := make([]Middleware, 0, len(a)+len(b))
	mids = append(mids, a...)
	mids = append(mids, b...)
	return mids
}
