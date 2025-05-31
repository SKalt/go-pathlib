package pathlib

import (
	"io/fs"
	"os"
)

type PathStr string
type Dir PathStr
type Symlink PathStr

// TODO: type Fifo
// TODO: type Device      PathStr
// TODO: type Regularfile PathStr

type PathOnDisk[P PathStr | Dir | Symlink] struct {
	original P // retained for error reporting in case Info is nil
	result[fs.FileInfo]
}
type result[T any] struct {
	value *T
	err   error
}
// type optional[T any] *T

// wrappers around [os.Lstat]/[os.Stat] and operations on the resulting [os.FileInfo]
type IFilePath interface {
	Stat() (os.FileInfo, error)
	LStat() (os.FileInfo, error)

	Exists() bool

	IsRegular() bool
	IsDir() bool
	IsSymLink() bool
	IsDevice() bool
	IsCharDevice() bool
	IsSocket() bool
	IsTemporary() bool
	IsFifo() bool
}

type Readable[T any] interface {
	Read() (T, error) // TODO: ReadAll() (T, error)?
}

// String-only path operations that do not require filesystem access.
type PurePath[Self PathStr | Dir | Symlink] interface {
	// Navigation
	Join(...string) PathStr
	Parent() Dir
	NearestDir() Dir

	IsAbsolute() bool
	IsLocal() bool

	Abs(cwd Dir) Self
	Rel(target Dir) (Self, error)
	Localize() Self
	// TODO: ToLocal(cwd string)

	BaseName() string
	// Ext() string
}
type ExistingManipulator[P PathStr | Dir] interface {
	Remove() error
	Rename(newName string) PathOnDisk[P]
	Move(newPath PathStr) error
	Info() (os.FileInfo, error)
}
type DirManipulator interface {
	ExistingManipulator[Dir]
	Mkdir(perm os.FileMode) PathOnDisk[Dir]
	MkdirAll(perm os.FileMode) PathOnDisk[Dir]
	MkdirTemp(pattern string) PathOnDisk[Dir]
}

type FileManipulator[P PathStr] interface {
	ExistingManipulator[P]
	Open() (*os.File, error)
}

var (
	_ PurePath[PathStr] = PathStr(".")
	_ PurePath[Dir] = Dir(".")
	_ PurePath[Symlink] = Symlink("link")

	_ Readable[[]os.DirEntry] = Dir(".")
	_ Readable[PathStr]       = Symlink("link")
)
