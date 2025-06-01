package pathlib

import (
	"os"
)


func (p OnDisk[P]) Parent() Dir {
	return p.Path().Parent()
}

func (p OnDisk[P]) Path() P {
	return P(p.Name())
}

func (p *OnDisk[P]) Chmod(mode os.FileMode) error {
	return os.Chmod(p.Name(), mode)
}

func (p OnDisk[P]) IsRegular() (isRegular bool) {
	return p.Mode().IsRegular()
}
