package pathlib

import (
	"errors"
	"io/fs"
)

func (p PathOnDisk[any]) Exists() bool {
	return  !errors.Is(p.err, fs.ErrExist)
}

func (p PathOnDisk[any]) IsRegular() (isRegular bool) {
	return p.value != nil  && (*p.value).Mode().IsRegular()
}
