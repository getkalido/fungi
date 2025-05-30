package fungi

// Flatten converts from a fungi.Stream[[]T] to a fungi.Stream[T]. It takes
// a stream of slices of T and returns individual T elements in order.
//
// This does not handle deduplication of T, which can be performed with
// [fungi.UniqueBy].
//
// Example:
//
//	someStream := fungi.SliceStream([][]int64{{1, 2, 3}, {4, 5}})
//	slice, err := fungi.Loop(func (in int64) {
//		fmt.Println(in)
//	})(fungi.Flatten[int64](someStream))
func Flatten[T any](source Stream[[]T]) Stream[T] {
	return &flatten[T]{source: source}
}

type flatten[T any] struct {
	source       Stream[[]T]
	currentSlice []T
}

func (c *flatten[T]) Next() (item T, err error) {
	if c.source == nil {
		panic("nil source for flatmap")
	}

	for len(c.currentSlice) == 0 {
		c.currentSlice, err = c.source.Next()
		if err != nil {
			return
		}
	}

	item = c.currentSlice[0]
	c.currentSlice = c.currentSlice[1:]
	return
}
