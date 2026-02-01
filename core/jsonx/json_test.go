package jsonx

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	v := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}
	bs, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"John","age":30}`, string(bs))
}

func TestMarshalToString(t *testing.T) {
	v := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}
	toString, err := MarshalToString(v)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"John","age":30}`, toString)

	_, err = MarshalToString(make(chan int))
	assert.NotNil(t, err)
}

func TestUnmarshal(t *testing.T) {
	const s = `{"name":"John","age":30}`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := Unmarshal([]byte(s), &v)
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalError(t *testing.T) {
	const s = `{"name":"John","age":30`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := Unmarshal([]byte(s), &v)
	assert.NotNil(t, err)
}

func TestUnmarshalFromString(t *testing.T) {
	const s = `{"name":"John","age":30}`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalFromString(s, &v)
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalFromStringError(t *testing.T) {
	const s = `{"name":"John","age":30`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalFromString(s, &v)
	assert.NotNil(t, err)
}

func TestUnmarshalFromRead(t *testing.T) {
	const s = `{"name":"John","age":30}`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalFromReader(strings.NewReader(s), &v)
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalFromReaderError(t *testing.T) {
	const s = `{"name":"John","age":30`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalFromReader(strings.NewReader(s), &v)
	assert.NotNil(t, err)
}

func Test_doMarshalJson(t *testing.T) {
	type args struct {
		v any
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{nil},
			want:    []byte("null"),
			wantErr: assert.NoError,
		},
		{
			name:    "string",
			args:    args{"hello"},
			want:    []byte(`"hello"`),
			wantErr: assert.NoError,
		},
		{
			name:    "int",
			args:    args{42},
			want:    []byte("42"),
			wantErr: assert.NoError,
		},
		{
			name:    "bool",
			args:    args{true},
			want:    []byte("true"),
			wantErr: assert.NoError,
		},
		{
			name: "struct",
			args: args{
				struct {
					Name string `json:"name"`
				}{Name: "test"},
			},
			want:    []byte(`{"name":"test"}`),
			wantErr: assert.NoError,
		},
		{
			name:    "slice",
			args:    args{[]int{1, 2, 3}},
			want:    []byte("[1,2,3]"),
			wantErr: assert.NoError,
		},
		{
			name:    "map",
			args:    args{map[string]int{"a": 1, "b": 2}},
			want:    []byte(`{"a":1,"b":2}`),
			wantErr: assert.NoError,
		},
		{
			name:    "unmarshalable type",
			args:    args{complex(1, 2)},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "channel type",
			args:    args{make(chan int)},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "url with query params",
			args:    args{"https://example.com/api?name=test&age=25"},
			want:    []byte(`"https://example.com/api?name=test&age=25"`),
			wantErr: assert.NoError,
		},
		{
			name:    "url with encoded query params",
			args:    args{"https://example.com/api?data=hello%20world&special=%26%3D"},
			want:    []byte(`"https://example.com/api?data=hello%20world&special=%26%3D"`),
			wantErr: assert.NoError,
		},
		{
			name:    "url with multiple query params",
			args:    args{"http://localhost:8080/users?page=1&limit=10&sort=name&order=asc"},
			want:    []byte(`"http://localhost:8080/users?page=1&limit=10&sort=name&order=asc"`),
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.v)
			if !tt.wantErr(t, err, fmt.Sprintf("Marshal(%v)", tt.args.v)) {
				return
			}

			assert.Equalf(t, string(tt.want), string(got), "Marshal(%v)", tt.args.v)
		})
	}
}

func TestMarshalWithBuffer(t *testing.T) {
	v := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}

	bs, err := MarshalWithBuffer(v)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"John","age":30}`, string(bs))

	// Test consistency with Marshal
	bs_marshal, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, string(bs_marshal), string(bs))
}

func TestMarshalWithBufferError(t *testing.T) {
	_, err := MarshalWithBuffer(make(chan int))
	assert.NotNil(t, err)
}

func TestMarshalWithBufferReuse(t *testing.T) {
	// First marshal
	v1 := map[string]string{"key": "value1"}
	bs1, err := MarshalWithBuffer(v1)
	assert.Nil(t, err)
	assert.Equal(t, `{"key":"value1"}`, string(bs1))

	// Reset buffer and reuse
	v2 := map[string]string{"key": "value2"}
	bs2, err := MarshalWithBuffer(v2)
	assert.Nil(t, err)
	assert.Equal(t, `{"key":"value2"}`, string(bs2))
}

func TestMarshalWithBufferMultipleTypes(t *testing.T) {
	tests := []struct {
		name string
		args any
		want string
	}{
		{
			name: "nil",
			args: nil,
			want: "null",
		},
		{
			name: "string",
			args: "hello",
			want: `"hello"`,
		},
		{
			name: "int",
			args: 42,
			want: "42",
		},
		{
			name: "bool",
			args: true,
			want: "true",
		},
		{
			name: "struct",
			args: struct {
				Name string `json:"name"`
			}{Name: "test"},
			want: `{"name":"test"}`,
		},
		{
			name: "slice",
			args: []int{1, 2, 3},
			want: "[1,2,3]",
		},
		{
			name: "map",
			args: map[string]int{"a": 1, "b": 2},
			want: `{"a":1,"b":2}`,
		},
		{
			name: "url with special characters",
			args: "https://example.com/api?name=test&age=25",
			want: `"https://example.com/api?name=test&age=25"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalWithBuffer(tt.args)
			assert.Nil(t, err, "MarshalWithBuffer should not return error for %v", tt.args)
			assert.Equal(t, tt.want, string(got), "MarshalWithBuffer(%v)", tt.args)
		})
	}
}

// 基准测试数据结构
type Person struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Age         int               `json:"age"`
	Address     Address           `json:"address"`
	PhoneNumber string            `json:"phone_number"`
	IsActive    bool              `json:"is_active"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

func generateTestData(size int) []Person {
	people := make([]Person, size)
	for i := 0; i < size; i++ {
		people[i] = Person{
			ID:    i + 1,
			Name:  fmt.Sprintf("Person %d", i+1),
			Email: fmt.Sprintf("person%d@example.com", i+1),
			Age:   20 + (i % 50),
			Address: Address{
				Street:  fmt.Sprintf("%d Main St", (i+1)*100),
				City:    "New York",
				State:   "NY",
				ZipCode: fmt.Sprintf("100%02d", i%100),
				Country: "USA",
			},
			PhoneNumber: fmt.Sprintf("+1-555-%04d", i+1),
			IsActive:    i%2 == 0,
			Tags:        []string{"tag1", "tag2", "tag3"},
			Metadata: map[string]string{
				"department": "Engineering",
				"level":      "Senior",
				"location":   "Remote",
			},
		}
	}
	return people
}

// 基准测试：对比 Marshal 和 MarshalWithBuffer 的性能
func BenchmarkMarshal_vs_MarshalWithBuffer(b *testing.B) {
	person := generateTestData(1)[0]

	b.Run("Marshal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := Marshal(person)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MarshalWithBuffer", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := MarshalWithBuffer(person)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// 测试不同大小数据的性能
func BenchmarkMarshal_DifferentSizes(b *testing.B) {
	sizes := []int{1, 10, 100, 1000}

	for _, size := range sizes {
		data := generateTestData(size)

		b.Run(fmt.Sprintf("Marshal_%d_items", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Marshal(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(fmt.Sprintf("MarshalWithBuffer_%d_items", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := MarshalWithBuffer(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// 内存分配对比测试
func BenchmarkMarshal_MemoryAllocs(b *testing.B) {
	person := generateTestData(1)[0]

	b.Run("Marshal_Allocs", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := Marshal(person)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MarshalWithBuffer_Allocs", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := MarshalWithBuffer(person)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// 并发性能测试
func BenchmarkMarshal_Concurrent(b *testing.B) {
	person := generateTestData(1)[0]

	b.Run("Marshal_Concurrent", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := Marshal(person)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("MarshalWithBuffer_Concurrent", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := MarshalWithBuffer(person)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// 测试简单数据类型的性能
func BenchmarkMarshal_SimpleTypes(b *testing.B) {
	simpleData := map[string]any{
		"string": "hello world",
		"int":    42,
		"bool":   true,
		"slice":  []int{1, 2, 3, 4, 5},
		"map":    map[string]int{"a": 1, "b": 2, "c": 3},
	}

	b.Run("Marshal_SimpleTypes", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := Marshal(simpleData)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MarshalWithBuffer_SimpleTypes", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := MarshalWithBuffer(simpleData)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
