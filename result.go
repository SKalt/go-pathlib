package pathlib

type Result[T any] struct {
	Val T
	Err error
}

func (r Result[T]) Unwrap() T {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.Val
}

func (r Result[T]) Unpack() (T, error) {
	return r.Val, r.Err
}

func (r Result[T]) IsOk() bool {
	return r.Err == nil
}
