package pathlib

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func (p PathStr) OnDisk() (onDisk *OnDisk[PathStr], err error) {
	var info os.FileInfo
	info, err = lstat(p)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	onDisk = &OnDisk[PathStr]{info}
	return
}

// Note: single-field structs have the same size as their field

func (p PathStr) Exists() (exists bool) {
	_, err := p.OnDisk()
	return !errors.Is(err, fs.ErrNotExist)
}

// Returns true if the path is absolute, false otherwise.
func (p PathStr) IsAbsolute() bool {
	return isAbsolute(p)
}

// returns true if the path is local/relative, false otherwise.
func (p PathStr) IsLocal() bool {
	return isLocal(p)
}

func (p PathStr) Read() (result any, err error) {
	var onDisk *OnDisk[PathStr]
	onDisk, err = p.OnDisk()
	if err != nil {
		return
	}
	mode := (*onDisk).Mode()

	if mode.IsRegular() {
		result, err = os.ReadFile(string(p))
	} else if mode.IsDir() {
		result, err = Dir(p).Read()
	} else if isSymLink(mode) {
		result, err = Symlink(p).Read()
	} else if isCharDevice(mode) {
		// TODO
	} else if isDevice(mode) {
		// TODO
	} else if isFifo(mode) {
		// TODO
	} else if isSocket(mode) {
		// TODO
	} else if isTemporary(mode) {
		// TODO
	}
	return
}

func (p PathStr) Open() (*os.File, error) {
	return os.Open(string(p))
}

func (p PathStr) WithOpen(cb func(*os.File) error) error { // FIXME: name
	f, err := p.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	return cb(f)
}

// Abs implements PurePath.
func (p PathStr) Abs(cwd Dir) PathStr {
	if p.IsAbsolute() {
		return p
	}
	// assume that cwd is absolute
	if !cwd.IsAbsolute() {
		panic("cwd must be absolute")
	}
	return cwd.Join(string(p)) // join the path with the current working directory
}

// Localize implements PurePath.
func (p PathStr) Localize() PathStr {
	filepath.Localize(string(p))
	panic("unimplemented") // TODO: implement this
}

// Rel implements PurePath.
func (p PathStr) Rel(target Dir) (PathStr, error) {
	result, err := filepath.Rel(string(p), string(target))
	if err != nil {
		return "", errors.Join(err, errors.New("unable to make relative path"))
	}
	return PathStr(result), nil
}
