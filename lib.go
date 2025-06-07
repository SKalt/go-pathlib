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

type OnDisk[P kind] struct{ fs.FileInfo }

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
	Abs(cwd Dir) Self
	Rel(target Dir) (Self, error)
	Localize() Self
	OnDisk() (*OnDisk[Self], error)
}

type Manipulator[P kind] interface {
	Remove() error
	Rename(newName string) (*OnDisk[P], error)
	Move(newPath PathStr) (*OnDisk[P], error)
	// Info() (os.FileInfo, error)
}
type DirCreator interface {
	// see [os.Mkdir]
	Mkdir(perm os.FileMode) (*OnDisk[Dir], error)
	// see [os.MkdirAll]
	MkdirAll(perm os.FileMode) (*OnDisk[Dir], error)
	// see [os.MkdirTemp]
	MkdirTemp(pattern string) (*OnDisk[Dir], error)
}

type FileManipulator[P kind] interface {
	Manipulator[P]
	Open(flag int, mode os.FileMode) (*os.File, error)
}

// TODO: type SymlinkManipulator interface {}
