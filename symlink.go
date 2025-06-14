package pathlib

import (
	"os"
)

type Symlink PathStr

// -----------------------------------------------------------------------------
var _ Readable[PathStr] = Symlink("")

// Read implements Readable.
func (s Symlink) Read() (PathStr, error) {
	link, err := os.Readlink(string(s))
	return PathStr(link), err
}

// -----------------------------------------------------------------------------
var _ PurePath = Symlink("")

// BaseName implements PurePath.
func (s Symlink) BaseName() string {
	return PathStr(s).BaseName()
}

// IsAbsolute implements PurePath.
func (s Symlink) IsAbsolute() bool {
	return PathStr(s).IsAbsolute()
}

// IsLocal implements PurePath.
func (s Symlink) IsLocal() bool {
	return PathStr(s).IsLocal()
}

// Join implements [PurePath].
func (s Symlink) Join(parts ...string) PathStr {
	return PathStr(s).Join(parts...)
}

// Parent implements PurePath.
func (s Symlink) Parent() Dir {
	return PathStr(s).Parent()
}

// -----------------------------------------------------------------------------
var _ Transformer[Symlink] = Symlink("")

// Abs implements Transformer.
func (s Symlink) Abs() (Symlink, error) {
	q, err := PathStr(s).Abs()
	return Symlink(q), err
}

// Localize implements Transformer.
func (s Symlink) Localize() (Symlink, error) {
	q, err := PathStr(s).Localize()
	return Symlink(q), err
}

// Rel implements Transformer.
func (s Symlink) Rel(target Dir) (Symlink, error) {
	result, err := PathStr(s).Rel(target)
	return Symlink(result), err
}

func (s Symlink) ExpandUser() (Symlink, error) {
	q, err := PathStr(s).ExpandUser()
	return Symlink(q), err
}

func (s Symlink) Ext() string {
	return PathStr(s).Ext()
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Symlink] = Symlink("")

// OnDisk implements Beholder.
func (s Symlink) OnDisk() (OnDisk[Symlink], error) {
	actual, err := PathStr(s).OnDisk()
	if err != nil {
		return nil, err
	}
	if !isSymLink(actual.Mode()) {
		return nil, WrongTypeOnDisk[Symlink]{actual}
	}
	return onDisk[Symlink]{actual}, nil
}

// Exists implements Beholder.
func (s Symlink) Exists() bool {
	panic("unimplemented")
}

// Lstat implements Beholder.
func (s Symlink) Lstat() (OnDisk[Symlink], error) {
	panic("unimplemented")
}

// Stat implements Beholder.
func (s Symlink) Stat() (OnDisk[Symlink], error) {
	panic("unimplemented")
}

// // https://go.dev/play/p/mWNvcZLrjog
// // https://godbolt.org/z/1caPfvzfh

// func temp[T kind]() {
// 	switch any((*T)(nil)).(type) {
// 	case *PathStr:
// 	}
// }
