package metainfo

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

func TestRegisterCustomKeys(t *testing.T) {
	reset()
	RegisterCustomKeys([]string{"a"})
	assert.ElementsMatch(t, []string{"a"}, customKeyStore.keyArr)
	RegisterCustomKeys([]string{"b"})
	assert.ElementsMatch(t, []string{"a", "b"}, customKeyStore.keyArr)
	RegisterCustomKeys([]string{"a", "c"})
	assert.ElementsMatch(t, []string{"a", "b", "c"}, customKeyStore.keyArr)

	t.Run("should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			RegisterCustomKeys([]string{"Aaa"})
		})
	})
}

func TestCustomKeys(t *testing.T) {
	reset()

	testKeys := []string{"a", "b"}
	RegisterCustomKeys(testKeys)

	customMap := map[string]string{
		"a": "a",
		"b": "b",
	}

	header := http.Header{}
	header.Add("a", "a")
	header.Add("b", "b")
	header.Add("notInclude", "notInclude")

	ck := CustomKeysMapPropagator

	checkCtx := func(ctx context.Context) {
		mp := getMap(ctx)
		for _, k := range testKeys {
			assert.Equal(t, k, mp[k])
		}

		_, ok := mp["notInclude"]
		assert.False(t, ok)
	}

	// test http header Extract
	ctx2 := ck.Extract(context.Background(), propagation.HeaderCarrier(header))
	checkCtx(ctx2)

	// test http header Inject
	header2 := http.Header{}
	ck.Inject(ctx2, propagation.HeaderCarrier(header2))
	for _, k := range testKeys {
		assert.Equal(t, header.Get(k), header2.Get(k))
	}

	// test http get map
	assert.Equal(t, customMap, GetMapFromContext(ctx2))
	assert.Equal(t, customMap, GetMapFromPropagator(propagation.HeaderCarrier(header)))

	// test grpc metadata Extract
	md := metadata.New(map[string]string{
		"a":          "a",
		"b":          "b",
		"notInclude": "notInclude",
	})
	ctx3 := ck.Extract(context.Background(), GrpcHeaderCarrier(md))
	checkCtx(ctx3)

	// test grpc metadata Inject
	md2 := metadata.MD{}
	ck.Inject(ctx3, GrpcHeaderCarrier(md2))
	for _, k := range testKeys {
		assert.Equal(t, md.Get(k), md2.Get(k))
	}

	// test grpc get map
	assert.Equal(t, customMap, GetMapFromContext(ctx3))
	assert.Equal(t, customMap, GetMapFromPropagator(GrpcHeaderCarrier(md2)))

	// extract multi times
	header11 := http.Header{}
	header11.Add("a", "a")
	ctx11 := ck.Extract(context.Background(), propagation.HeaderCarrier(header11))
	header12 := http.Header{}
	header12.Add("b", "b")
	ctx12 := ck.Extract(ctx11, propagation.HeaderCarrier(header12))
	checkCtx(ctx12)
}

func TestCustomKeys1(t *testing.T) {
	reset()

	header := http.Header{}
	header.Add("a", "a")
	header.Add("b", "b")
	header.Add("notInclude", "notInclude")

	checkCtx := func(ctx context.Context) {
		mp := getMap(ctx)
		assert.Equal(t, 0, len(mp))
	}

	ck := CustomKeysMapPropagator
	// test http header Extract
	ctx2 := ck.Extract(context.Background(), propagation.HeaderCarrier(header))
	checkCtx(ctx2)
}

func TestRegisterCustomKeys_AutoPrefix(t *testing.T) {
	reset()

	customMap := map[string]string{
		"x-pass-a": "x-pass-a",
		"x-pass-b": "x-pass-b",
	}

	header := http.Header{}
	header.Add("x-pass-a", "x-pass-a")
	header.Add("x-pass-b", "x-pass-b")
	header.Add("notInclude", "notInclude")

	ck := CustomKeysMapPropagator

	checkCtx := func(ctx context.Context) {
		mp := getMap(ctx)
		for k, v := range customMap {
			assert.Equal(t, v, mp[k])
		}

		_, ok := mp["notInclude"]
		assert.False(t, ok)
	}

	// test http header Extract
	ctx2 := ck.Extract(context.Background(), propagation.HeaderCarrier(header))
	checkCtx(ctx2)

	// test http header Inject
	header2 := http.Header{}
	ck.Inject(ctx2, propagation.HeaderCarrier(header2))
	for k := range customMap {
		assert.Equal(t, header.Get(k), header2.Get(k))
	}

	// test http get map
	assert.Equal(t, customMap, GetMapFromContext(ctx2))
	assert.Equal(t, customMap, GetMapFromPropagator(propagation.HeaderCarrier(header)))

	// test grpc metadata Extract
	md := metadata.New(map[string]string{
		"x-pass-a":   "x-pass-a",
		"x-pass-b":   "x-pass-b",
		"notInclude": "notInclude",
	})
	ctx3 := ck.Extract(context.Background(), GrpcHeaderCarrier(md))
	checkCtx(ctx3)

	// test grpc metadata Inject
	md2 := metadata.MD{}
	ck.Inject(ctx3, GrpcHeaderCarrier(md2))
	for k := range customMap {
		assert.Equal(t, md.Get(k), md2.Get(k))
	}

	// test grpc get map
	assert.Equal(t, customMap, GetMapFromContext(ctx3))
	assert.Equal(t, customMap, GetMapFromPropagator(GrpcHeaderCarrier(md2)))

	// extract multi times
	header11 := http.Header{}
	header11.Add("x-pass-a", "x-pass-a")
	ctx11 := ck.Extract(context.Background(), propagation.HeaderCarrier(header11))
	header12 := http.Header{}
	header12.Add("x-pass-b", "x-pass-b")
	ctx12 := ck.Extract(ctx11, propagation.HeaderCarrier(header12))
	checkCtx(ctx12)
}

func TestRegisterCustomKeys_Mix(t *testing.T) {
	reset()

	customMap := map[string]string{
		"a":        "a",
		"b":        "b",
		"x-pass-a": "x-pass-a",
	}
	RegisterCustomKeys([]string{"a", "b"})

	header := http.Header{}
	for k, v := range customMap {
		header.Add(k, v)
	}

	ck := CustomKeysMapPropagator

	checkCtx := func(ctx context.Context) {
		mp := getMap(ctx)
		for k, v := range customMap {
			assert.Equal(t, v, mp[k])
		}
	}

	// test http header Extract
	ctx2 := ck.Extract(context.Background(), propagation.HeaderCarrier(header))
	checkCtx(ctx2)

	// test http header Inject
	header2 := http.Header{}
	ck.Inject(ctx2, propagation.HeaderCarrier(header2))
	for k := range customMap {
		assert.Equal(t, header.Get(k), header2.Get(k))
	}

	// test http get map
	assert.Equal(t, customMap, GetMapFromContext(ctx2))
	assert.Equal(t, customMap, GetMapFromPropagator(propagation.HeaderCarrier(header)))
}

//goos: darwin
//goarch: arm64
//pkg: code.bydev.io/cht/fiat/backend/lib.git/pkg/transport
//BenchmarkCustomKeys
//BenchmarkCustomKeys-8              	 1414234	       762.8 ns/op
//BenchmarkCustomKeys_10
//BenchmarkCustomKeys_10-8           	  615555	      1849 ns/op
//BenchmarkCustomKeys_50
//BenchmarkCustomKeys_50-8           	  104818	     11497 ns/op
//BenchmarkCustomKeysAutoPass
//BenchmarkCustomKeysAutoPass-8      	  861883	      1333 ns/op
//BenchmarkCustomKeysAutoPass_10
//BenchmarkCustomKeysAutoPass_10-8   	  392179	      3085 ns/op
//BenchmarkCustomKeysAutoPass_50
//BenchmarkCustomKeysAutoPass_50-8   	   75937	     15628 ns/op
//BenchmarkCustomKeysMix
//BenchmarkCustomKeysMix-8           	 1201923	       972.0 ns/op
//BenchmarkCustomKeysMix_10
//BenchmarkCustomKeysMix_10-8        	  450882	      2786 ns/op
//BenchmarkCustomKeysMix_50
//BenchmarkCustomKeysMix_50-8        	   82384	     14543 ns/op

func benchmarkLen(b *testing.B, l int) {
	reset()

	testKeys := make([]string, l)
	h := http.Header{}
	for i := 0; i < len(testKeys); i++ {
		testKeys[i] = fmt.Sprintf("%d", i)
		h.Add(testKeys[i], testKeys[i])
	}

	RegisterCustomKeys(testKeys)
	ck := CustomKeysMapPropagator
	h2 := http.Header{}

	for i := 0; i < b.N; i++ {
		ctx1 := ck.Extract(context.Background(), propagation.HeaderCarrier(h))
		ck.Inject(ctx1, propagation.HeaderCarrier(h2))
	}
}

func benchmarkAutoPassLen(b *testing.B, l int) {
	reset()

	testKeys := make([]string, l)
	h := http.Header{}
	for i := 0; i < len(testKeys); i++ {
		testKeys[i] = fmt.Sprintf("%s%d", PrefixPass, i)
		h.Add(testKeys[i], testKeys[i])
	}

	ck := CustomKeysMapPropagator
	h2 := http.Header{}

	for i := 0; i < b.N; i++ {
		ctx1 := ck.Extract(context.Background(), propagation.HeaderCarrier(h))
		ck.Inject(ctx1, propagation.HeaderCarrier(h2))
	}
}

func benchmarkMixLen(b *testing.B, l int) {
	reset()

	testKeys := make([]string, 0)
	h := http.Header{}
	for i := 0; i < l; i++ {
		var key string
		if i%2 == 0 {
			key = fmt.Sprintf("%d", i)
			testKeys = append(testKeys, key)
		} else {
			key = fmt.Sprintf("%s%d", PrefixPass, i)
		}
		h.Add(key, key)
	}

	RegisterCustomKeys(testKeys)
	ck := CustomKeysMapPropagator
	h2 := http.Header{}

	for i := 0; i < b.N; i++ {
		ctx1 := ck.Extract(context.Background(), propagation.HeaderCarrier(h))
		ck.Inject(ctx1, propagation.HeaderCarrier(h2))
	}
}

func BenchmarkCustomKeys(b *testing.B) {
	benchmarkLen(b, 5)
}

func BenchmarkCustomKeys_10(b *testing.B) {
	benchmarkLen(b, 10)
}

func BenchmarkCustomKeys_50(b *testing.B) {
	benchmarkLen(b, 50)
}

func BenchmarkCustomKeysAutoPass(b *testing.B) {
	benchmarkAutoPassLen(b, 5)
}

func BenchmarkCustomKeysAutoPass_10(b *testing.B) {
	benchmarkAutoPassLen(b, 10)
}

func BenchmarkCustomKeysAutoPass_50(b *testing.B) {
	benchmarkAutoPassLen(b, 50)
}

func BenchmarkCustomKeysMix(b *testing.B) {
	benchmarkMixLen(b, 5)
}

func BenchmarkCustomKeysMix_10(b *testing.B) {
	benchmarkMixLen(b, 10)
}

func BenchmarkCustomKeysMix_50(b *testing.B) {
	benchmarkMixLen(b, 50)
}
