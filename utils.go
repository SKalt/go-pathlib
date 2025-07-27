package pathlib

import (
	"io/fs"
	"os"
)

// utility function
func expect[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func isSymLink(m os.FileMode) bool {
	return (m & fs.ModeSymlink) == fs.ModeSymlink
}
