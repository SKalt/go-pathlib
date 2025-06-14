package pathlib

import (
	"io/fs"
	"os"
)

type onDisk[P kind] struct{ fs.FileInfo }

var _ OnDisk[PathStr] = onDisk[PathStr]{}

func (p onDisk[P]) Path() P {
	return P(p.Name())
}

// func (p onDisk[P]) IsRegular() (isRegular bool) {
// 	return p.Mode().IsRegular()
// }

var _ fs.FileInfo = onDisk[PathStr]{}

// PurePath --------------------------------------------------------------------
var _ PurePath = onDisk[PathStr]{}

// Implements PurePath
func (p onDisk[P]) Parent() Dir {
	return p.Path().Parent()
}

// BaseName implements PurePath.
func (p onDisk[P]) BaseName() string {
	return p.Path().BaseName()
}

// Ext implements PurePath.
func (p onDisk[P]) Ext() string {
	return p.Path().Ext()
}

// IsAbsolute implements PurePath.
func (p onDisk[P]) IsAbsolute() bool {
	return p.Path().IsAbsolute()
}

// IsLocal implements PurePath.
func (p onDisk[P]) IsLocal() bool {
	return p.Path().IsLocal()
}

// Join implements PurePath.
func (p onDisk[P]) Join(parts ...string) PathStr {
	return p.Path().Join(parts...)
}
func typeIs[A, B kind]() (typesAreEqual bool) {
	_, typesAreEqual = any((*A)(nil)).(*B)
	return
}

// NearestDir implements PurePath.
func (p onDisk[P]) NearestDir() Dir {
	if typeIs[P, Dir]() {
		return Dir(p.Path())
	} else {
		return p.Parent()
	}
}

// Transformer -----------------------------------------------------------------
var _ Transformer[PathStr] = onDisk[PathStr]{}

// Abs implements Transformer.
func (p onDisk[P]) Abs() (P, error) {
	abs, err := PathStr(p.Path()).Abs()
	return P(abs), err
}

// Localize implements Transformer.
func (p onDisk[P]) Localize() (P, error) {
	q, err := PathStr(p.Path()).Localize()
	return P(q), err
}

// Rel implements Transformer.
func (p onDisk[P]) Rel(target Dir) (P, error) {
	q, err := PathStr(p.Path()).Rel(target)
	return P(q), err
}

func (p onDisk[P]) ExpandUser() (P, error) {
	q, err := PathStr(p.Path()).ExpandUser()
	return P(q), err
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[PathStr] = onDisk[PathStr]{}

// Remove implements Manipulator.
func (p onDisk[P]) Remove() error {
	return os.Remove(p.Name())
}

// Rename implements Manipulator.
func (p onDisk[P]) Rename(destination PathStr) (result P, err error) {
	result = P(destination)
	err = os.Rename(p.Name(), string(destination))
	return
}

func (p onDisk[P]) Chmod(mode fs.FileMode) (result P, err error) {
	result = p.Path()
	// os.Lchown()
	err = os.Chmod(string(result), mode)
	return
}

func (p onDisk[P]) Chown(uid, gid int) (result P, err error) {
	result = p.Path()
	err = os.Lchown(string(result), uid, gid)
	return
}
