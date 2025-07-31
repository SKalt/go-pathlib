package pathlib

import (
	"fmt"
	"io/fs"
)

type WrongTypeOnDisk[P Kind] struct{ fs.FileInfo }

func (w WrongTypeOnDisk[P]) Error() string {
	return fmt.Sprintf(
		"%T(%s) unexpectedly found a file with mode `%s` on disk",
		(*P)(nil), // reveal the type name as "*T"
		w.Name(),
		w.Mode(),
	)[1:] // chop the leading "*" off "*T"
}
