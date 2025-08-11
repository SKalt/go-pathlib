package pathlib

type Result[T any] struct {
	val T
	err error
}

func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.val
}

func (r Result[T]) Unpack() (T, error) {
	return r.val, r.err
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}
