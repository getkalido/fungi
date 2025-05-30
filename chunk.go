package fungi

import "io"

// Chunk converts from a fungi.Stream[T] to a fungi.Stream[[]T]. It takes
// a stream of T and returns chunked slices of T in order.
//
// Example:
//
//	someStream := fungi.SliceStream([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//	err := fungi.Loop(func(in []int64) {
//		fmt.Println(in)
//	})(fungi.Chunk[int64](3)(someStream))
func Chunk[T any](chunkSize int) StreamTransformer[T, []T] {
	if chunkSize <= 0 {
		panic("chunk size must be greater than 0")
	}
	return func(source Stream[T]) Stream[[]T] {
		return &chunk[T]{
			source:       source,
			chunkSize:    chunkSize,
			currentSlice: make([]T, 0, chunkSize),
		}
	}
}

type chunk[T any] struct {
	source       Stream[T]
	chunkSize    int
	currentSlice []T
	err          error
}

func (c *chunk[T]) Next() (items []T, err error) {
	if c.source == nil {
		panic("nil source for chunking")
	}
	if c.err != nil {
		return nil, c.err
	}

	for len(c.currentSlice) < c.chunkSize {
		item, err := c.source.Next()
		if err == nil {
			c.currentSlice = append(c.currentSlice, item)
			continue
		}
		c.err = err
		if err != io.EOF {
			return nil, err
		}
		if len(c.currentSlice) == 0 {
			return nil, io.EOF
		}
		break
	}

	var chunkLen int
	if len(c.currentSlice) < c.chunkSize {
		chunkLen = len(c.currentSlice)
	} else {
		chunkLen = c.chunkSize
	}
	items = c.currentSlice[:chunkLen]
	c.currentSlice = c.currentSlice[chunkLen:]
	return items, nil
}
