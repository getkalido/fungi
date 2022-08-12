package fungi

import "io"

func DoneChanStream[T, D any](items chan T, done chan D) Stream[T] {
	return &doneChanStream[T, D]{
		source: items,
		done:   done,
	}
}

func ChanStream[T any](items chan T) Stream[T] {
	return &chanStream[T]{
		source: items,
	}
}

type doneChanStream[T, D any] struct {
	source chan T
	done   chan D
}

func (dcs *doneChanStream[T, D]) Next() (item T, err error) {
	select {
	case item = <-dcs.source:
	case <-dcs.done:
		err = io.EOF
	}
	return
}

type chanStream[T any] struct {
	source chan T
}

func (cs *chanStream[T]) Next() (item T, err error) {
	item, more := <-cs.source
	if !more {
		err = io.EOF
	}
	return
}
