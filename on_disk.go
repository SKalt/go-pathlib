package pathlib

import (
	"io/fs"
	"os"
	"time"
)

type onDisk[P kind] struct {
	fs.FileInfo
	observed time.Time
}

var _ OnDisk[PathStr] = onDisk[PathStr]{}

func (p onDisk[P]) Path() P {
	return P(p.Name())
}

func (p onDisk[P]) Observed() time.Time {
	return p.observed
}

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

func (p onDisk[P]) Parts() []string {
	return p.Path().Parts()
}

// func typeIs[A, B kind]() (typesAreEqual bool) {
// 	_, typesAreEqual = any((*A)(nil)).(*B)
// 	return
// }

// Transformer -----------------------------------------------------------------
var _ Transformer[PathStr] = onDisk[PathStr]{}
var _ InfallibleTransformer[PathStr] = onDisk[PathStr]{}

func (p onDisk[P]) Eq(q P) bool {
	return PathStr(p.Path()).Eq(PathStr(q))
}

// Clean implements Transformer
func (p onDisk[P]) Clean() P {
	return P(PathStr(p.Path()).Clean())
}

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

// MustExpandUser implements InfallibleTransformer.
func (p onDisk[P]) MustExpandUser() P {
	return expect(p.ExpandUser())
}

// MustLocalize implements InfallibleTransformer.
func (p onDisk[P]) MustLocalize() P {
	return expect(p.Localize())
}

// MustMakeAbs implements InfallibleTransformer.
func (p onDisk[P]) MustMakeAbs() P {
	return expect(p.Abs())
}

// MustMakeRel implements InfallibleTransformer.
func (p onDisk[P]) MustMakeRel(target Dir) P {
	return expect(p.Rel(target))
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[PathStr] = onDisk[PathStr]{}
var _ InfallibleManipulator[PathStr] = onDisk[PathStr]{}

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
	err = os.Chmod(string(result), mode)
	return
}

func (p onDisk[P]) Chown(uid, gid int) (result P, err error) {
	result = p.Path()
	err = os.Lchown(string(result), uid, gid)
	return
}

// MustChmod implements InfallibleManipulator.
func (p onDisk[P]) MustChmod(mode os.FileMode) P {
	return expect(p.Chmod(mode))
}

// MustChown implements InfallibleManipulator.
func (p onDisk[P]) MustChown(uid int, gid int) P {
	return expect(p.Chown(uid, gid))
}

// MustRemove implements InfallibleManipulator.
func (p onDisk[P]) MustRemove() {
	expect[any](nil, p.Remove())
}

// MustRename implements InfallibleManipulator.
func (p onDisk[P]) MustRename(newPath PathStr) P {
	return expect(p.Rename(newPath))
}
