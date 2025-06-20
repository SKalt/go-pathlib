package pathlib

import (
	"io"
	"io/fs"
	"os"
)

// TODO: type Fifo
// TODO: type Device      PathStr
// TODO: type RegularFile PathStr

// type _pathKind uint8
// const (
// 	kindUnknown _pathKind = iota
// 	kindDir
// 	kindSymlink
// )
// type kinder interface {
// 	PurePath
// 	kind() _pathKind
// }

type kind interface {
	PurePath
	~string // TODO: re-restrict to Dir | File | Symlink?
	// see https://blog.chewxy.com/2018/03/18/golang-interfaces/#sealed-interfaces
}

type ContextManager[T io.Closer] interface {
	With(func(T) error) error
}

type Readable[T any] interface {
	Read() (T, error)
}

// String-only path operations that do not require filesystem access.
type PurePath interface {
	// Navigation
	Join(...string) PathStr // TODO: handle joining an absolute path
	Parent() Dir
	BaseName() string
	Ext() string
	Parts() []string
	// -> volume (windows-only)
	// -> iter.Seq[AnyPath]

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

type Beholder[PathKind kind] interface {
	OnDisk() (OnDisk[PathKind], error)
	Stat() (OnDisk[PathKind], error)
	Lstat() (OnDisk[PathKind], error)
	Exists() bool
}

type OnDisk[PathKind kind] interface {
	fs.FileInfo
	PurePath
	Transformer[PathKind]
	// Readable[any] // refining the type of what gets read would
	// require passing an additional type parameter, which
	// causes weird type-states to become possible, like OnDisk[Dir, struct{...}]
	// TODO: Observed() time.Time?
}

type Maker[T any] interface { // ~Fallible
	// see [os.Create]
	Make(perm ...fs.FileMode) (T, error)
	MustMake(perm ...fs.FileMode) T
	// TODO: add mode, parents args?
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

// TODO: Dir and Symlink should add a .RemoveAll()

type DirCreator interface { // ~Fallible x3
	// see [os.Mkdir]
	Mkdir(perm fs.FileMode) (Dir, error)
	// see [os.MkdirAll]
	MkdirAll(perm fs.FileMode) (Dir, error)
	// see [os.MkdirTemp]
	MkdirTemp(pattern string) (Dir, error)
}

type FileManipulator[P kind] interface {
	Manipulator[P]
	Open(flag int, mode os.FileMode) (*os.File, error) // ~Fallible
}

func Cwd() (Dir, error) {
	dir, err := os.Getwd()
	return Dir(dir), err
}

func UserHomeDir() (Dir, error) {
	dir, err := os.UserHomeDir()
	return Dir(dir), err
}
func UserCacheDir() (Dir, error) {
	dir, err := os.UserCacheDir()
	return Dir(dir), err
}
func UserConfigDir() (Dir, error) {
	dir, err := os.UserConfigDir()
	return Dir(dir), err
}

// returns the process/os-wide temporary directory
func TempDir() Dir {
	return Dir(os.TempDir())
}

// TODO: type SymlinkManipulator interface {}

// utility function
func expect[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
