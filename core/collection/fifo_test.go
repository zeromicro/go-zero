package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueOrder(t *testing.T) {
	tests := []struct {
		name         string
		size         int
		initial      []int
		takeBefore   int
		additional   []int
		wantCapacity int
	}{
		{
			name:         "within initial capacity",
			size:         8,
			initial:      []int{1, 2, 3},
			wantCapacity: 8,
		},
		{
			name:         "grow from beginning",
			size:         2,
			initial:      []int{1, 2, 3},
			wantCapacity: 4,
		},
		{
			name:         "grow after wrapping",
			size:         4,
			initial:      []int{1, 2, 3, 4},
			takeBefore:   1,
			additional:   []int{5, 6},
			wantCapacity: 8,
		},
		{
			name:         "grow repeatedly",
			size:         1,
			initial:      sequence(20),
			wantCapacity: 32,
		},
		{
			name:         "grow above threshold",
			size:         queueGrowThreshold,
			initial:      sequence(queueGrowThreshold + 1),
			wantCapacity: queueGrowThreshold * 2,
		},
		{
			name:         "grow well above threshold",
			size:         1024,
			initial:      sequence(1025),
			wantCapacity: 1472,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := NewQueue(test.size)
			assert.True(t, queue.Empty())

			for _, value := range test.initial {
				queue.Put(value)
			}
			for i := 0; i < test.takeBefore; i++ {
				value, ok := queue.Take()
				assert.True(t, ok)
				assert.Equal(t, test.initial[i], value)
			}
			for _, value := range test.additional {
				queue.Put(value)
			}

			want := append([]int(nil), test.initial[test.takeBefore:]...)
			want = append(want, test.additional...)
			for _, expected := range want {
				actual, ok := queue.Take()
				assert.True(t, ok)
				assert.Equal(t, expected, actual)
			}

			assert.Equal(t, test.wantCapacity, len(queue.elements))
			assert.True(t, queue.Empty())
			_, ok := queue.Take()
			assert.False(t, ok)
		})
	}
}

func TestQueueTakeClearsElement(t *testing.T) {
	tests := []struct {
		name       string
		size       int
		operations string
	}{
		{
			name:       "take from beginning",
			size:       2,
			operations: "ppt",
		},
		{
			name:       "take after wrapping",
			size:       2,
			operations: "pptptt",
		},
		{
			name:       "take after growing",
			size:       2,
			operations: "pppt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := NewQueue(test.size)
			value := 0

			for _, operation := range test.operations {
				switch operation {
				case 'p':
					value++
					element := new(int)
					*element = value
					queue.Put(element)
				case 't':
					index := queue.head
					_, ok := queue.Take()
					assert.True(t, ok)
					assert.Nil(t, queue.elements[index])
				default:
					t.Fatalf("unknown operation: %q", operation)
				}
			}
		})
	}
}

func TestNewQueueWithInvalidSize(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "zero",
		},
		{
			name: "negative",
			size: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Panics(t, func() {
				NewQueue(test.size)
			})
		})
	}
}

func BenchmarkQueueGrowth(b *testing.B) {
	tests := []struct {
		name  string
		size  int
		count int
	}{
		{
			name:  "initial_1_count_256",
			size:  1,
			count: 256,
		},
		{
			name:  "initial_1_count_4096",
			size:  1,
			count: 4096,
		},
		{
			name:  "initial_1_count_65536",
			size:  1,
			count: 65536,
		},
		{
			name:  "initial_8_count_4096",
			size:  8,
			count: 4096,
		},
		{
			name:  "initial_256_count_4096",
			size:  256,
			count: 4096,
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			elements := make([]any, test.count)
			for i := range elements {
				elements[i] = i
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				queue := NewQueue(test.size)
				for _, element := range elements {
					queue.Put(element)
				}
			}
		})
	}
}

func BenchmarkQueueWrappedGrowth(b *testing.B) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "capacity_8",
			size: 8,
		},
		{
			name: "capacity_256",
			size: 256,
		},
		{
			name: "capacity_4096",
			size: 4096,
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			elements := make([]any, test.size)
			for i := range elements {
				elements[i] = i
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				queue := NewQueue(test.size)
				for _, element := range elements {
					queue.Put(element)
				}
				for j := 0; j < test.size/2; j++ {
					queue.Take()
				}
				for j := 0; j < test.size/2; j++ {
					queue.Put(elements[j])
				}

				// The queue is full and wrapped. One more Put triggers growth
				// and copies both sides of the ring into FIFO order.
				queue.Put(elements[0])
			}
		})
	}
}

func sequence(size int) []int {
	values := make([]int, size)
	for i := range values {
		values[i] = i
	}
	return values
}
