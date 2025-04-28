package fungi

import "io"

// Concat joins together multiple streams of type T. For each input stream,
// it will stream all elements of the input stream, before it begins
// returning elements from the next input stream.
//
// Example
//
//	err := fungi.Loop(func(in int64) {
//		fmt.Println(in)
//	})(
//		fungi.Concat(
//			fungi.SliceStream([]int64{1, 2, 3}),
//			fungi.SliceStream([]int64{4, 5}),
//		),
//	)
func Concat[T any](streams ...Stream[T]) Stream[T] {
	return &concat[T]{sources: streams}
}

type concat[T any] struct {
	sources []Stream[T]
}

func (c *concat[T]) Next() (item T, err error) {
	for len(c.sources) > 0 {
		firstSource := c.sources[0]
		item, err = firstSource.Next()
		if err == nil {
			return item, err
		}
		if err != io.EOF {
			return
		}
		c.sources = c.sources[1:]
	}

	err = io.EOF
	return
}
