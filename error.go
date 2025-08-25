package pathlib

import (
	"fmt"
)

type WrongTypeOnDisk[P Kind] struct{ Observed Info[P] }

func (w WrongTypeOnDisk[P]) Error() string {
	path := w.Observed.Path()
	return fmt.Sprintf(
		"%T(%q) unexpectedly has mode %s on-disk",
		path, path,
		w.Observed.Mode(),
	)
}
