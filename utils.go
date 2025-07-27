package pathlib

import (
	"io/fs"
	"os"
)

// Panics if err != nil. This probably-inlined function is defined to cut down on boilerplate
// when transforming fallible `Method() (T, error)`s into infallible `MustMethod() T`s.,
func expect[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func isSymLink(m os.FileMode) bool {
	return (m & fs.ModeSymlink) == fs.ModeSymlink
}
