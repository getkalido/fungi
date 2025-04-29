package fungi

// Loop calls an infallible function on each element of the input stream.
//
// See [ForEach] for a version which can return errors.
func Loop[T any](handle func(T)) StreamHandler[T] {
	return ForEach(func(item T) error { handle(item); return nil })
}
