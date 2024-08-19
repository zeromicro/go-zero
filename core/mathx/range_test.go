package mathx

import "testing"

func TestAtLeast(t *testing.T) {
	t.Run("test int", func(t *testing.T) {
		if got := AtLeast(10, 5); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(3, 5); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(5, 5); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test int8", func(t *testing.T) {
		if got := AtLeast(int8(10), int8(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(int8(3), int8(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(int8(5), int8(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test int16", func(t *testing.T) {
		if got := AtLeast(int16(10), int16(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(int16(3), int16(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(int16(5), int16(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test int32", func(t *testing.T) {
		if got := AtLeast(int32(10), int32(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(int32(3), int32(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(int32(5), int32(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test int64", func(t *testing.T) {
		if got := AtLeast(int64(10), int64(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(int64(3), int64(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(int64(5), int64(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test uint", func(t *testing.T) {
		if got := AtLeast(uint(10), uint(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(uint(3), uint(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(uint(5), uint(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test uint8", func(t *testing.T) {
		if got := AtLeast(uint8(10), uint8(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(uint8(3), uint8(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(uint8(5), uint8(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test uint16", func(t *testing.T) {
		if got := AtLeast(uint16(10), uint16(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(uint16(3), uint16(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(uint16(5), uint16(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test uint32", func(t *testing.T) {
		if got := AtLeast(uint32(10), uint32(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(uint32(3), uint32(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(uint32(5), uint32(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test uint64", func(t *testing.T) {
		if got := AtLeast(uint64(10), uint64(5)); got != 10 {
			t.Errorf("AtLeast() = %v, want 10", got)
		}
		if got := AtLeast(uint64(3), uint64(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
		if got := AtLeast(uint64(5), uint64(5)); got != 5 {
			t.Errorf("AtLeast() = %v, want 5", got)
		}
	})

	t.Run("test float32", func(t *testing.T) {
		if got := AtLeast(float32(10.0), float32(5.0)); got != 10.0 {
			t.Errorf("AtLeast() = %v, want 10.0", got)
		}
		if got := AtLeast(float32(3.0), float32(5.0)); got != 5.0 {
			t.Errorf("AtLeast() = %v, want 5.0", got)
		}
		if got := AtLeast(float32(5.0), float32(5.0)); got != 5.0 {
			t.Errorf("AtLeast() = %v, want 5.0", got)
		}
	})

	t.Run("test float64", func(t *testing.T) {
		if got := AtLeast(10.0, 5.0); got != 10.0 {
			t.Errorf("AtLeast() = %v, want 10.0", got)
		}
		if got := AtLeast(3.0, 5.0); got != 5.0 {
			t.Errorf("AtLeast() = %v, want 5.0", got)
		}
		if got := AtLeast(5.0, 5.0); got != 5.0 {
			t.Errorf("AtLeast() = %v, want 5.0", got)
		}
	})
}

func TestAtMost(t *testing.T) {
	t.Run("test int", func(t *testing.T) {
		if got := AtMost(10, 5); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(3, 5); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(5, 5); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test int8", func(t *testing.T) {
		if got := AtMost(int8(10), int8(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(int8(3), int8(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(int8(5), int8(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test int16", func(t *testing.T) {
		if got := AtMost(int16(10), int16(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(int16(3), int16(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(int16(5), int16(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test int32", func(t *testing.T) {
		if got := AtMost(int32(10), int32(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(int32(3), int32(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(int32(5), int32(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test int64", func(t *testing.T) {
		if got := AtMost(int64(10), int64(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(int64(3), int64(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(int64(5), int64(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test uint", func(t *testing.T) {
		if got := AtMost(uint(10), uint(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(uint(3), uint(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(uint(5), uint(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test uint8", func(t *testing.T) {
		if got := AtMost(uint8(10), uint8(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(uint8(3), uint8(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(uint8(5), uint8(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test uint16", func(t *testing.T) {
		if got := AtMost(uint16(10), uint16(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(uint16(3), uint16(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(uint16(5), uint16(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test uint32", func(t *testing.T) {
		if got := AtMost(uint32(10), uint32(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(uint32(3), uint32(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(uint32(5), uint32(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test uint64", func(t *testing.T) {
		if got := AtMost(uint64(10), uint64(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
		if got := AtMost(uint64(3), uint64(5)); got != 3 {
			t.Errorf("AtMost() = %v, want 3", got)
		}
		if got := AtMost(uint64(5), uint64(5)); got != 5 {
			t.Errorf("AtMost() = %v, want 5", got)
		}
	})

	t.Run("test float32", func(t *testing.T) {
		if got := AtMost(float32(10.0), float32(5.0)); got != 5.0 {
			t.Errorf("AtMost() = %v, want 5.0", got)
		}
		if got := AtMost(float32(3.0), float32(5.0)); got != 3.0 {
			t.Errorf("AtMost() = %v, want 3.0", got)
		}
		if got := AtMost(float32(5.0), float32(5.0)); got != 5.0 {
			t.Errorf("AtMost() = %v, want 5.0", got)
		}
	})

	t.Run("test float64", func(t *testing.T) {
		if got := AtMost(10.0, 5.0); got != 5.0 {
			t.Errorf("AtMost() = %v, want 5.0", got)
		}
		if got := AtMost(3.0, 5.0); got != 3.0 {
			t.Errorf("AtMost() = %v, want 3.0", got)
		}
		if got := AtMost(5.0, 5.0); got != 5.0 {
			t.Errorf("AtMost() = %v, want 5.0", got)
		}
	})
}

func TestBetween(t *testing.T) {
	t.Run("test int", func(t *testing.T) {
		if got := Between(10, 5, 15); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(3, 5, 15); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(20, 5, 15); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(5, 5, 15); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(15, 5, 15); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test int8", func(t *testing.T) {
		if got := Between(int8(10), int8(5), int8(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(int8(3), int8(5), int8(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(int8(20), int8(5), int8(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(int8(5), int8(5), int8(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(int8(15), int8(5), int8(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test int16", func(t *testing.T) {
		if got := Between(int16(10), int16(5), int16(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(int16(3), int16(5), int16(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(int16(20), int16(5), int16(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(int16(5), int16(5), int16(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(int16(15), int16(5), int16(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test int32", func(t *testing.T) {
		if got := Between(int32(10), int32(5), int32(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(int32(3), int32(5), int32(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(int32(20), int32(5), int32(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(int32(5), int32(5), int32(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(int32(15), int32(5), int32(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test int64", func(t *testing.T) {
		if got := Between(int64(10), int64(5), int64(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(int64(3), int64(5), int64(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(int64(20), int64(5), int64(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(int64(5), int64(5), int64(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(int64(15), int64(5), int64(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test uint", func(t *testing.T) {
		if got := Between(uint(10), uint(5), uint(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(uint(3), uint(5), uint(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint(20), uint(5), uint(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(uint(5), uint(5), uint(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint(15), uint(5), uint(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test uint8", func(t *testing.T) {
		if got := Between(uint8(10), uint8(5), uint8(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(uint8(3), uint8(5), uint8(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint8(20), uint8(5), uint8(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(uint8(5), uint8(5), uint8(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint8(15), uint8(5), uint8(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test uint16", func(t *testing.T) {
		if got := Between(uint16(10), uint16(5), uint16(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(uint16(3), uint16(5), uint16(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint16(20), uint16(5), uint16(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(uint16(5), uint16(5), uint16(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint16(15), uint16(5), uint16(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test uint32", func(t *testing.T) {
		if got := Between(uint32(10), uint32(5), uint32(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(uint32(3), uint32(5), uint32(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint32(20), uint32(5), uint32(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(uint32(5), uint32(5), uint32(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint32(15), uint32(5), uint32(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test uint64", func(t *testing.T) {
		if got := Between(uint64(10), uint64(5), uint64(15)); got != 10 {
			t.Errorf("Between() = %v, want 10", got)
		}
		if got := Between(uint64(3), uint64(5), uint64(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint64(20), uint64(5), uint64(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
		if got := Between(uint64(5), uint64(5), uint64(15)); got != 5 {
			t.Errorf("Between() = %v, want 5", got)
		}
		if got := Between(uint64(15), uint64(5), uint64(15)); got != 15 {
			t.Errorf("Between() = %v, want 15", got)
		}
	})

	t.Run("test float32", func(t *testing.T) {
		if got := Between(float32(10.0), float32(5.0), float32(15.0)); got != 10.0 {
			t.Errorf("Between() = %v, want 10.0", got)
		}
		if got := Between(float32(3.0), float32(5.0), float32(15.0)); got != 5.0 {
			t.Errorf("Between() = %v, want 5.0", got)
		}
		if got := Between(float32(20.0), float32(5.0), float32(15.0)); got != 15.0 {
			t.Errorf("Between() = %v, want 15.0", got)
		}
		if got := Between(float32(5.0), float32(5.0), float32(15.0)); got != 5.0 {
			t.Errorf("Between() = %v, want 5.0", got)
		}
		if got := Between(float32(15.0), float32(5.0), float32(15.0)); got != 15.0 {
			t.Errorf("Between() = %v, want 15.0", got)
		}
	})

	t.Run("test float64", func(t *testing.T) {
		if got := Between(10.0, 5.0, 15.0); got != 10.0 {
			t.Errorf("Between() = %v, want 10.0", got)
		}
		if got := Between(3.0, 5.0, 15.0); got != 5.0 {
			t.Errorf("Between() = %v, want 5.0", got)
		}
		if got := Between(20.0, 5.0, 15.0); got != 15.0 {
			t.Errorf("Between() = %v, want 15.0", got)
		}
		if got := Between(5.0, 5.0, 15.0); got != 5.0 {
			t.Errorf("Between() = %v, want 5.0", got)
		}
		if got := Between(15.0, 5.0, 15.0); got != 15.0 {
			t.Errorf("Between() = %v, want 15.0", got)
		}
	})
}
