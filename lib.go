package pathlib

import (
	"io/fs"
	"os"
	"time"
)

type kind interface { // TODO: make public?
	PurePath
	~string // TODO: re-restrict to Dir | File | Symlink?
	// see https://blog.chewxy.com/2018/03/18/golang-interfaces/#sealed-interfaces
}

type Readable[T any] interface {
	Read() (T, error)
}
type InfallibleReader[T any] interface {
	MustRead() T
}

// String-only path operations that do not require filesystem access.
type PurePath interface {
	// Navigation
	Join(...string) PathStr
	Parent() Dir
	BaseName() string
	Ext() string
	Parts() []string

	IsAbsolute() bool
	IsLocal() bool
}

// transforms the appearance of a path, but not what it represents.
type Transformer[Self kind] interface { // ~Fallible x3
	Abs() (Self, error)
	Rel(target Dir) (Self, error)
	Localize() (Self, error)
	ExpandUser() (Self, error)
	Clean() Self
	Eq(other Self) bool
}
type InfallibleTransformer[Self kind] interface {
	MustMakeAbs() Self
	MustMakeRel(target Dir) Self
	MustLocalize() Self
	MustExpandUser() Self
}

type Beholder[PathKind kind] interface {
	OnDisk() (OnDisk[PathKind], error)
	Stat() (OnDisk[PathKind], error)
	Lstat() (OnDisk[PathKind], error)
	Exists() bool
}
type InfallibleBeholder[PathKind kind] interface {
	// OnDisk implements Beholder.
	MustBeOnDisk() OnDisk[PathKind]
	MustStat() OnDisk[PathKind]
	MustLstat() OnDisk[PathKind]
}

type OnDisk[PathKind kind] interface {
	fs.FileInfo
	PurePath
	Transformer[PathKind]
	// Readable[any] // refining the type of what gets read would
	// require passing an additional type parameter, which
	// causes weird type-states to become possible, like OnDisk[Dir, struct{...}]
	Observed() time.Time
}

type Maker[T any] interface { // ~Fallible
	// see [os.Create]
	Make(perm ...fs.FileMode) (T, error)
}

type InfallibleMaker[T any] interface {
	MustMake(perm ...fs.FileMode) T
}

type Manipulator[PathKind kind] interface { // ~Fallible x3 + 1
	// see [os.Remove]
	Remove() error
	// see [os.Chmod]
	Chmod(os.FileMode) (PathKind, error)
	// see [os.Chown]
	Chown(uid, gid int) (PathKind, error)
	// see [os.Rename]
	Rename(newPath PathStr) (PathKind, error)
}
type InfallibleManipulator[PathKind kind] interface {
	// see [os.Remove]. Panics if Remove fails.
	MustRemove()
	// see [os.Chmod]. Panics if Chmod fails.
	MustChmod(mode os.FileMode) PathKind
	// see [os.Chown]. Panics if Chown fails.
	MustChown(uid, gid int) PathKind
	// see [os.Rename]. Panics if Rename fails.
	MustRename(newPath PathStr) PathKind
}

// TODO: Dir and Symlink should add a .RemoveAll()

type DirCreator interface { // ~Fallible x3
	// see [os.Mkdir]
	Mkdir(perm fs.FileMode) (Dir, error)
	// see [os.MkdirAll]
	MkdirAll(perm fs.FileMode) (Dir, error)
	// see [os.MkdirTemp]
	MkdirTemp(pattern string) (Dir, error)
}

type Destroyer interface {
	RemoveAll() error
}
type InfallibleDestroyer interface {
	MustRemoveAll()
}

type FileManipulator[P kind] interface {
	Manipulator[P]
	Open(flag int, mode os.FileMode) (*os.File, error) // ~Fallible
}

// TODO: infallible methods/interfaces
