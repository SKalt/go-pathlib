package pathlib

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Dir PathStr

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

func (d Dir) Chdir() error {
	return os.Chdir(string(d))
}

// Readable --------------------------------------------------------------------
var _ Readable[[]os.DirEntry] = Dir(".")

// a wrapper around [os.ReadDir]. Read() returns all the entries of the directory sorted
// by filename. If an error occurred reading the directory, Read returns the entries it was
// able to read before the error, along with the error.
func (d Dir) Read() ([]os.DirEntry, error) {
	return os.ReadDir(string(d))
}

// PurePath --------------------------------------------------------------------
var _ PurePath = Dir(".")

// BaseName implements PurePath.
func (d Dir) BaseName() string {
	return PathStr(d).BaseName()
}

// IsAbsolute implements PurePath.
func (d Dir) IsAbsolute() bool {
	return PathStr(d).IsAbsolute()
}

// IsLocal implements PurePath.
func (d Dir) IsLocal() bool {
	return PathStr(d).IsLocal()
}

// Join implements PurePath.
func (d Dir) Join(parts ...string) PathStr {
	return PathStr(d).Join(parts...)
}

// NearestDir implements PurePath.
func (d Dir) NearestDir() Dir {
	return d
}

// Parent implements PurePath.
func (d Dir) Parent() Dir {
	return PathStr(d).Parent()
}

// Ext implements PurePath.
func (d Dir) Ext() string {
	return PathStr(d).Ext()
}

// Transformer -----------------------------------------------------------------
var _ Transformer[Dir] = Dir(".")

// Abs implements Transformer.
func (d Dir) Abs() (Dir, error) {
	abs, err := PathStr(d).Abs()
	return Dir(abs), err
}

// Localize implements Transformer.
func (d Dir) Localize() (Dir, error) {
	q, err := PathStr(d).Localize()
	return Dir(q), err // this will panic if d is not absolute
}

// Rel implements Transformer.
func (d Dir) Rel(target Dir) (Dir, error) {
	relative, err := PathStr(d).Rel(target)
	return Dir(relative), err
}

// ExpandUser implements TildeTransformer.
func (d Dir) ExpandUser() (Dir, error) {
	p, err := PathStr(d).ExpandUser()
	return Dir(p), err
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Dir] = Dir(".")

func (d Dir) OnDisk() (OnDisk[Dir], error) {
	actual, err := PathStr(d).OnDisk()
	if err != nil {
		return nil, err
	}
	if !actual.IsDir() {
		return nil, WrongTypeOnDisk[Dir]{actual}
	}
	return onDisk[Dir]{actual}, nil
}

// Exists implements Beholder.
func (root Dir) Exists() bool {
	return PathStr(root).Exists()
}

// Lstat implements Beholder.
func (root Dir) Lstat() (OnDisk[Dir], error) {
	return root.OnDisk()
}

// Stat implements Beholder.
func (root Dir) Stat() (result OnDisk[Dir], err error) {
	var info fs.FileInfo
	info, err = os.Stat(string(root))
	if err != nil {
		return
	}
	if !info.IsDir() {
		err = WrongTypeOnDisk[Dir]{info}
		return
	}
	result = onDisk[Dir]{info}
	return
}
