package pathlib

import (
	"fmt"
)

type WrongTypeOnDisk[P Kind] struct{ OnDisk[P] }

func (w WrongTypeOnDisk[P]) Error() string {
	path := w.Path()
	return fmt.Sprintf(
		"%T(%q) unexpectedly has mode %s on-disk",
		path, path,
		w.Mode(),
	)
}

type Error[P Kind] struct {
	Path P
	Op   string
	Err  error
}

func (err Error[P]) Error() string {
	return fmt.Sprintf("%T(%q): %s: %s", err.Path, err.Path, err.Op, err.Err.Error())
}

func (err Error[P]) Unwrap() error {
	return err.Err
}

// TODO: figure out errors.Is() semantics
// func (err Error[P]) Is(target error) bool {
// 	if t, ok := target.(Error[P]); ok {
// 		result := t.Op == err.Op
// 		if inner, ok := err.Err.(interface {Is(error) bool}); ok {
// 			result = result && inner.Is(t.Err)
// 		}
// 		return result
// 	}
// 	return false

// }