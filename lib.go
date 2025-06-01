package pathlib

import (
	"io"
	"io/fs"
	"os"
)

type PathStr string

type Dir PathStr

type Symlink PathStr

// TODO: type Fifo
// TODO: type Device      PathStr
// TODO: type Regularfile PathStr

type kind interface {
	PurePath
	PathStr | Dir | Symlink
	// see https://blog.chewxy.com/2018/03/18/golang-interfaces/#sealed-interfaces
}

type ContextManager[T io.Closer] interface {
	With(func(T) error) error
}

type OnDisk[P PurePath] struct{ fs.FileInfo }

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

type Transmogrifier[Self kind] interface {
	Abs(cwd Dir) Self
	Rel(target Dir) (Self, error)
	Localize() Self
}

type ExistingManipulator[P kind] interface {
	Remove() error
	Rename(newName string) OnDisk[P]
	Move(newPath PathStr) error
	Info() (os.FileInfo, error)
}
type DirManipulator interface {
	ExistingManipulator[Dir]
	Mkdir(perm os.FileMode) OnDisk[Dir]
	MkdirAll(perm os.FileMode) OnDisk[Dir]
	MkdirTemp(pattern string) OnDisk[Dir]
}

type FileManipulator[P kind] interface {
	ExistingManipulator[P]
	Open(flag int, mode os.FileMode) (*os.File, error)
}

var (
	_ PurePath = PathStr(".")
	_ PurePath = Dir(".")
	_ PurePath = Symlink("link")

	_ Readable[[]os.DirEntry] = Dir(".")
	_ Readable[PathStr]       = Symlink("link")
)
