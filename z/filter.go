package z

// Generator is the type that generates the elements of T.
type Generator[T any] chan T

// Next returns the next element in the collection.
func (f Generator[T]) Next() *T {
	c, ok := <-f
	if !ok {
		return nil
	}
	return &c
}

// Slice returns a slice of the collection.
func (f Generator[T]) Slice() []T {
	result := []T{}
	for x := range f {
		result = append(result, x)
	}
	return result
}

// FilterMap returns a generator that generates the results of applying the function to each element in the collection.
func FilterMap[T any, R any](collection []T, fn func(T, int) (R, bool)) Generator[R] {
	ch := make(chan R)
	go func() {
		for i, item := range collection {
			if r, ok := fn(item, i); ok {
				ch <- r
			}
		}
		close(ch)
	}()
	return ch
}
