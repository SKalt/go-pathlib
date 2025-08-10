package pathlib

import (
	"io/fs"
	"os"
)

// Note: single-field structs have the same size as their field
type onDisk[P Kind] struct {
	fs.FileInfo
}

var _ OnDisk[PathStr] = onDisk[PathStr]{}

// Path implements [OnDisk].
func (p onDisk[P]) Path() P {
	return P(p.Name())
}

var _ fs.FileInfo = onDisk[PathStr]{}

// PurePath --------------------------------------------------------------------
var _ PurePath = onDisk[PathStr]{}

// Parent implements [PurePath]
func (p onDisk[P]) Parent() Dir {
	return p.Path().Parent()
}

// BaseName implements [PurePath].
func (p onDisk[P]) BaseName() string {
	return p.Path().BaseName()
}

// Ext implements [PurePath].
func (p onDisk[P]) Ext() string {
	return p.Path().Ext()
}

// IsAbsolute implements [PurePath].
func (p onDisk[P]) IsAbsolute() bool {
	return p.Path().IsAbsolute()
}

// IsLocal implements [PurePath].
func (p onDisk[P]) IsLocal() bool {
	return p.Path().IsLocal()
}

// Join implements [PurePath].
func (p onDisk[P]) Join(parts ...string) PathStr {
	return p.Path().Join(parts...)
}

// Parts implements [PurePath].
func (p onDisk[P]) Parts() []string {
	return p.Path().Parts()
}

// Transformer -----------------------------------------------------------------
var _ Transformer[PathStr] = onDisk[PathStr]{}

func (p onDisk[P]) String() string {
	return p.Name()
}

func (p onDisk[P]) Eq(q P) bool {
	return PathStr(p.Path()).Eq(PathStr(q))
}

// Clean implements [Transformer]
func (p onDisk[P]) Clean() P {
	return clean(p.Path())
}

// Abs implements [Transformer].
func (p onDisk[P]) Abs() Result[P] {
	return abs(p.Path())
}

// Localize implements [Transformer].
func (p onDisk[P]) Localize() Result[P] {
	return localize(p.Path())
}

// Rel implements [Transformer].
func (p onDisk[P]) Rel(target Dir) Result[P] {
	return rel(p.Path(), target)
}

func (p onDisk[P]) ExpandUser() Result[P] {
	return expandUser(p.Path())
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[PathStr] = onDisk[PathStr]{}

// Remove implements [Manipulator].
func (p onDisk[P]) Remove() error {
	return os.Remove(p.Name())
}

// Rename implements [Manipulator].
func (p onDisk[P]) Rename(destination PathStr) Result[P] {
	return rename(p.Path(), destination)
}

func (p onDisk[P]) Chmod(mode fs.FileMode) Result[P] {
	return chmod(p.Path(), mode)
}

func (p onDisk[P]) Chown(uid, gid int) Result[P] {
	return chown(p.Path(), uid, gid)
}
