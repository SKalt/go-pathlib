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

// type Error[P Kind] struct {
// 	Path P
// 	Op   string
// 	Err  error
// }

// func (err Error[P]) Error() string {
// 	return fmt.Sprintf("%T(%q): %s: %s", err.Path, err.Path, err.Op, err.Err.Error())
// }

// func (err Error[P]) Unwrap() error {
// 	return err.Err
// }

// TODO: figure out errors.Is() semantics
