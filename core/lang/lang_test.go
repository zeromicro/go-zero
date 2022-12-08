package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepr(t *testing.T) {
	var (
		f32 float32 = 1.1
		f64         = 2.2
		i8  int8    = 1
		i16 int16   = 2
		i32 int32   = 3
		i64 int64   = 4
		u8  uint8   = 5
		u16 uint16  = 6
		u32 uint32  = 7
		u64 uint64  = 8
	)
	tests := []struct {
		v      interface{}
		expect string
	}{
		{
			nil,
			"",
		},
		{
			mockStringable{},
			"mocked",
		},
		{
			new(mockStringable),
			"mocked",
		},
		{
			newMockPtr(),
			"mockptr",
		},
		{
			&mockOpacity{
				val: 1,
			},
			"{1}",
		},
		{
			true,
			"true",
		},
		{
			false,
			"false",
		},
		{
			f32,
			"1.1",
		},
		{
			f64,
			"2.2",
		},
		{
			i8,
			"1",
		},
		{
			i16,
			"2",
		},
		{
			i32,
			"3",
		},
		{
			i64,
			"4",
		},
		{
			u8,
			"5",
		},
		{
			u16,
			"6",
		},
		{
			u32,
			"7",
		},
		{
			u64,
			"8",
		},
		{
			[]byte(`abcd`),
			"abcd",
		},
		{
			mockOpacity{val: 1},
			"{1}",
		},
	}

	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			assert.Equal(t, test.expect, Repr(test.v))
		})
	}
}

type mockStringable struct{}

func (m mockStringable) String() string {
	return "mocked"
}

type mockPtr struct{}

func newMockPtr() *mockPtr {
	return new(mockPtr)
}

func (m *mockPtr) String() string {
	return "mockptr"
}

type mockOpacity struct {
	val int
}
