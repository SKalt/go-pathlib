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

// Convenience method to cast get the untyped string representation of the path.
//
// String implements [Transformer].
func (p onDisk[P]) String() string {
	return string(p.Path())
}

// Returns true if the two paths represent the same path.
//
// Eq implements [Transformer].
func (p onDisk[P]) Eq(q P) bool {
	return PathStr(p.Path()).Eq(PathStr(q))
}

// Remove ".", "..", and repeated slashes from a path.
//
// See [path/filepath.Clean].
//
// Clean implements [Transformer]
func (p onDisk[P]) Clean() P {
	return clean(p.Path())
}

// Returns an absolute path, or an error if the path cannot be made absolute. Note that there may be more than one
// absolute path for a given input path.
//
// See [path/filepath.Abs].
//
// Abs implements [Transformer].
func (p onDisk[P]) Abs() (P, error) {
	return abs(p.Path())
}

// Localize implements [Transformer].
func (p onDisk[P]) Localize() (P, error) {
	return localize(p.Path())
}

// Returns a relative path to the target directory, or an error if the path cannot be made relative.
//
// See [path/filepath.Rel].
//
// Rel implements [Transformer].
func (p onDisk[P]) Rel(base Dir) (P, error) {
	return rel(base, p.Path())
}

func (p onDisk[P]) ExpandUser() (P, error) {
	return expandUser(p.Path())
}

// Mover --------------------------------------------------------------------
var _ Remover[PathStr] = onDisk[PathStr]{}

// See [os.Remove].
//
// Remove implements [Remover].
func (p onDisk[P]) Remove() error {
	return remove(p.Path())
}

// Rename implements [Remover].
func (p onDisk[P]) Rename(destination PathStr) (P, error) {
	return rename(p.Path(), destination)
}

// Changer ----------------------------------------------------------------------
var _ Changer = onDisk[PathStr]{}

// See [os.Chmod].
//
// Chmod implements [Changer].
func (p onDisk[P]) Chmod(mode fs.FileMode) error {
	return chmod(p.Path(), mode)
}

// See [os.Chown].
//
// Chown implements [Changer].
func (p onDisk[P]) Chown(uid, gid int) error {
	return chown(p.Path(), uid, gid)
}
