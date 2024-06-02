package chain

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// A constructor for middleware
// that writes its own "tag" into the RW and does nothing else.
// Useful in checking if a chain is behaving in the right order.
func tagMiddleware(tag string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tag))
			h.ServeHTTP(w, r)
		})
	}
}

// Not recommended (https://golang.org/pkg/reflect/#Value.Pointer),
// but the best we can do.
func funcsEqual(f1, f2 any) bool {
	val1 := reflect.ValueOf(f1)
	val2 := reflect.ValueOf(f2)
	return val1.Pointer() == val2.Pointer()
}

var testApp = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("app\n"))
})

func TestNew(t *testing.T) {
	c1 := func(h http.Handler) http.Handler {
		return nil
	}

	c2 := func(h http.Handler) http.Handler {
		return http.StripPrefix("potato", nil)
	}

	slice := []Middleware{c1, c2}
	c := New(slice...)
	for k := range slice {
		assert.True(t, funcsEqual(c.(chain).middlewares[k], slice[k]),
			"New does not add constructors correctly")
	}
}

func TestThenWorksWithNoMiddleware(t *testing.T) {
	assert.True(t, funcsEqual(New().Then(testApp), testApp),
		"Then does not work with no middleware")
}

func TestThenTreatsNilAsDefaultServeMux(t *testing.T) {
	assert.Equal(t, http.DefaultServeMux, New().Then(nil),
		"Then does not treat nil as DefaultServeMux")
}

func TestThenFuncTreatsNilAsDefaultServeMux(t *testing.T) {
	assert.Equal(t, http.DefaultServeMux, New().ThenFunc(nil),
		"ThenFunc does not treat nil as DefaultServeMux")
}

func TestThenFuncConstructsHandlerFunc(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	chained := New().ThenFunc(fn)
	rec := httptest.NewRecorder()

	chained.ServeHTTP(rec, (*http.Request)(nil))

	assert.Equal(t, reflect.TypeOf((http.HandlerFunc)(nil)), reflect.TypeOf(chained),
		"ThenFunc does not construct HandlerFunc")
}

func TestThenOrdersHandlersCorrectly(t *testing.T) {
	t1 := tagMiddleware("t1\n")
	t2 := tagMiddleware("t2\n")
	t3 := tagMiddleware("t3\n")

	chained := New(t1, t2, t3).Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTP(w, r)

	assert.Equal(t, "t1\nt2\nt3\napp\n", w.Body.String(),
		"Then does not order handlers correctly")
}

func TestAppendAddsHandlersCorrectly(t *testing.T) {
	c := New(tagMiddleware("t1\n"), tagMiddleware("t2\n"))
	c = c.Append(tagMiddleware("t3\n"), tagMiddleware("t4\n"))
	h := c.Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", http.NoBody)
	assert.Nil(t, err)

	h.ServeHTTP(w, r)
	assert.Equal(t, "t1\nt2\nt3\nt4\napp\n", w.Body.String(),
		"Append does not add handlers correctly")
}

func TestExtendAddsHandlersCorrectly(t *testing.T) {
	c := New(tagMiddleware("t3\n"), tagMiddleware("t4\n"))
	c = c.Prepend(tagMiddleware("t1\n"), tagMiddleware("t2\n"))
	h := c.Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	h.ServeHTTP(w, r)
	assert.Equal(t, "t1\nt2\nt3\nt4\napp\n", w.Body.String(),
		"Extend does not add handlers in correctly")
}
