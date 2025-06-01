package pathlib

import (
	"io/fs"
	"os"
	"path/filepath"
)

var (
	_ PurePath[Dir]               = Dir("/tmp")
)

// a wrapper around [os.ReadDir]. Read() returns all the entries of the directory sorted
// by filename. If an error occured reading the directory, Read returns the entries it was
// able to read before the error, along with the error.
func (d Dir) Read() ([]os.DirEntry, error) {
	return os.ReadDir(string(d))
}

// A wrapper around [path/filepath.WalkDir], which has the following properties:
//
// > The files are walked in lexical order, which makes the output deterministic but [reading the] entire directory into memory before proceeding to walk that directory.
//
// > WalkDir does not follow symbolic links.
//
// > WalkDir calls [callback] with paths that use the separator character appropriate for the operating system.
func (root Dir) Walk(
	callback func(path string, d fs.DirEntry, err error) error,
) error {
	return filepath.WalkDir(string(root), callback)
}

func (d Dir) Glob(pattern string) ([]string, error) {
	return filepath.Glob(filepath.Join(string(d), pattern))
}

// implmenting IPurePath

// BaseName implements IPurePath.
func (d Dir) BaseName() string {
	return PathStr(d).BaseName()
}

// IsAbsolute implements IPurePath.
func (d Dir) IsAbsolute() bool {
	return PathStr(d).IsAbsolute()
}

// IsLocal implements IPurePath.
func (d Dir) IsLocal() bool {
	return PathStr(d).IsLocal()
}

// Join implements IPurePath.
func (d Dir) Join(parts ...string) PathStr {
	return PathStr(d).Join(parts...)
}

// NearestDir implements IPurePath.
func (d Dir) NearestDir() Dir {
	return d
}

// Parent implements IPurePath.
func (d Dir) Parent() Dir {
	return PathStr(d).Parent()
}


// Abs implements PurePath.
func (d Dir) Abs(cwd Dir) Dir {
	return Dir(PathStr(d).Abs(cwd)) // this will panic if cwd is not absolute
}

// Localize implements PurePath.
func (d Dir) Localize() Dir {
	return Dir(PathStr(d).Localize()) // this will panic if d is not absolute
}

// Rel implements PurePath.
func (d Dir) Rel(target Dir) (Dir, error) {
	relative, err := PathStr(d).Rel(target)
	return Dir(relative), err
}

func (d Dir) Ext() string {
	return PathStr(d).Ext()
}

func (d Dir) Chdir() error {
	return os.Chdir(string(d))
}
