package pathlib

import (
	"io/fs"
	"os"
)


func (p *PathOnDisk[P]) Parent() Dir {
	return PathStr(p.original).Parent()
}

func (p *PathOnDisk[P]) Path() P {
	return p.original
}

func (p *PathOnDisk[P]) Chmod(mode os.FileMode) PathOnDisk[P] {
	result := p.Map(func(fi fs.FileInfo) (*fs.FileInfo, error) {
		err := os.Chmod(fi.Name(), mode)
		if err != nil {
			return nil, err
		}
		r := lstat(fi.Name())
		return r.value, r.err
	}).(result[fs.FileInfo])

	return PathOnDisk[P]{p.original, result}
}
