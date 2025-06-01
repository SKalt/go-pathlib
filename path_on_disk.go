package pathlib

import (
	"errors"
	"io/fs"
)

func (p PathOnDisk[P]) Exists() bool {
	return  !errors.Is(p.err, fs.ErrExist)
}

func (p PathOnDisk[P]) IsRegular() (isRegular bool) {
	return p.value != nil  && (*p.value).Mode().IsRegular()
}
