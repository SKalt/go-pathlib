package pathlib

import (
	"io"
	"io/fs"
	"os"
)

// TODO: type Fifo
// TODO: type Device      PathStr
// TODO: type RegularFile PathStr

type kind interface {
	PurePath
	PathStr | Dir | Symlink
	// see https://blog.chewxy.com/2018/03/18/golang-interfaces/#sealed-interfaces
}

type ContextManager[T io.Closer] interface {
	With(func(T) error) error
}

type Readable[T any] interface {
	Read() (T, error) // TODO: ReadAll() (T, error)?
}

// String-only path operations that do not require filesystem access.
type PurePath interface {
	// Navigation
	Join(...string) PathStr
	Parent() Dir
	NearestDir() Dir
	BaseName() string
	Ext() string

	IsAbsolute() bool
	IsLocal() bool
}

// transforms the appearance of a path, but not what it represents.
type Transformer[Self kind] interface {
	Abs() (Self, error)
	Rel(target Dir) (Self, error)
	Localize() (Self, error)
	ExpandUser() (Self, error)
	// Stat() (OnDisk[Self], error)
}
type OnDisk[PathKind kind] interface {
	fs.FileInfo
	PurePath
	Transformer[PathKind]
	// Readable[any] // refining the type of what gets read would
	// require passing an additional type parameter, which
	// causes weird type-states to become possible, like OnDisk[Dir, struct{...}]
}

type Manipulator[PathKind kind] interface {
	// see [os.Remove]
	Remove() error
	// see [os.Chmod]
	Chmod(os.FileMode) (PathKind, error)
	// see [os.Chown]
	Chown(uid, gid int) (PathKind, error)
	// see [os.Rename]
	Rename(newPath PathStr) (PathKind, error)
}

type DirCreator interface {
	// see [os.Mkdir]
	Mkdir(perm os.FileMode) (Dir, error)
	// see [os.MkdirAll]
	MkdirAll(perm os.FileMode) (Dir, error)
	// see [os.MkdirTemp]
	MkdirTemp(pattern string) (Dir, error)
}

type FileManipulator[P kind] interface {
	Manipulator[P]
	Open(flag int, mode os.FileMode) (*os.File, error)
}

func Cwd() (Dir, error) {
	dir, err := os.Getwd()
	return Dir(dir), err
}

// TODO: type SymlinkManipulator interface {}
