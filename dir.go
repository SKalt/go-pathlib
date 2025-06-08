package pathlib

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type Dir PathStr

var (
	_ PurePath                = Dir(".")
	_ Transformer[Dir]        = Dir(".")
	_ Readable[[]os.DirEntry] = Dir(".")
)

// a wrapper around [os.ReadDir]. Read() returns all the entries of the directory sorted
// by filename. If an error occurred reading the directory, Read returns the entries it was
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
func (d Dir) Abs() (Dir, error) {
	abs, err := PathStr(d).Abs()
	return Dir(abs), err
}

// Localize implements PurePath.
func (d Dir) Localize() (Dir, error) {
	q, err := PathStr(d).Localize()
	return Dir(q), err // this will panic if d is not absolute
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

func (d Dir) OnDisk() (*onDisk[Dir], error) {
	actual, err := PathStr(d).OnDisk()
	if err != nil {
		return nil, err
	}
	if !actual.IsDir() {
		return nil, errors.New("not a directory: " + string(d))
	}
	return &onDisk[Dir]{*actual}, nil
}

// ExpandUser implements TildeTransformer.
func (d Dir) ExpandUser() (Dir, error) {
	p, err := PathStr(d).ExpandUser()
	return Dir(p), err
}
