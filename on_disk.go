package pathlib

import (
	"io/fs"
)

type onDisk[P Kind] struct {
	// type path associated with the file info
	p P
	fs.FileInfo
}

var _ Info[PathStr] = onDisk[PathStr]{}

// Path implements [Info].
func (p onDisk[P]) Path() P {
	return p.p
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
	return string(p.Path())
}

func (p onDisk[P]) Eq(q P) bool {
	return PathStr(p.Path()).Eq(PathStr(q))
}

// Clean implements [Transformer]
func (p onDisk[P]) Clean() P {
	return clean(p.Path())
}

// Abs implements [Transformer].
func (p onDisk[P]) Abs() (P, error) {
	return abs(p.Path())
}

// Localize implements [Transformer].
func (p onDisk[P]) Localize() (P, error) {
	return localize(p.Path())
}

// Rel implements [Transformer].
func (p onDisk[P]) Rel(base Dir) (P, error) {
	return rel(base, p.Path())
}

func (p onDisk[P]) ExpandUser() (P, error) {
	return expandUser(p.Path())
}

// Mover --------------------------------------------------------------------
var _ Mover[PathStr] = onDisk[PathStr]{}

// Remove implements [Mover].
func (p onDisk[P]) Remove() error {
	return remove(p.Path())
}

// Rename implements [Mover].
func (p onDisk[P]) Rename(destination PathStr) (P, error) {
	return rename(p.Path(), destination)
}

// Changer ----------------------------------------------------------------------
var _ Changer = onDisk[PathStr]{}

// Chmod implements [Changer].
func (p onDisk[P]) Chmod(mode fs.FileMode) error {
	return chmod(p.Path(), mode)
}

// Chown implements [Changer].
func (p onDisk[P]) Chown(uid, gid int) error {
	return chown(p.Path(), uid, gid)
}
