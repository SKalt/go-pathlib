package pathlib

import (
	"errors"
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

// invariant: if result.err == nil, then result.value != nil
type result[T any] struct {
	value *T
	err   error
}

// Map implements Residual.
func (r result[T]) Map(f func(T) (*T, error)) Residual[T] {
	if r.err != nil {
		return r // propagate the error
	}
	next := result[T]{}
	next.value, next.err = f(*r.value)
	return next
}

// Or implements Residual.
func (r result[T]) Or(T) T {
	panic("unimplemented")
}

// OrElse implements Residual.
func (r result[T]) OrElse(func(err error) T) T {
	panic("unimplemented")
}

// Unwrap implements Residual.
func (r result[T]) Unwrap() T {
	panic("unimplemented")
}
func (r result[T]) Ok() error {
	return r.err
}

type optional[T any] struct{ ptr *T }

// Map implements Residual.
func (o optional[T]) Map(f func(T) (*T, error)) Residual[T] {
	if o.ptr == nil {
		return o
	} else {
		next := optional[T]{}
		val, err := f(*o.ptr)
		if err == nil {
			next.ptr = val
		}
		return next
	}
}
func (o optional[T]) Or(defaultValue T) T {
	if o.ptr == nil {
		return defaultValue
	}
	return *o.ptr
}
func (o optional[T]) OrElse(f func(err error) T) T {
	if o.ptr == nil {
		return f(nil)
	}
	return *o.ptr
}
func (o optional[T]) Unwrap() T {
	if o.ptr == nil {
		panic("attempted to unwrap an empty optional")
	}
	return *o.ptr
}
func (o optional[T]) Ok() (error) {
	if o.ptr == nil {
		return missing
	}
	return nil
}
var missing = errors.New("missing value")

type Residual[T any] interface {
	Map(func(T) (*T, error)) Residual[T]
	Or(T) T
	OrElse(func(err error) T) T
	Unwrap() T
	Ok() error
}

var _ Residual[string] = optional[string]{}
var _ Residual[int] = result[int]{}

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
	Ext() string
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
	Open(flag int, mode os.FileMode) (*os.File, error)
}

var (
	_ PurePath[PathStr] = PathStr(".")
	_ PurePath[Dir]     = Dir(".")
	_ PurePath[Symlink] = Symlink("link")

	_ Readable[[]os.DirEntry] = Dir(".")
	_ Readable[PathStr]       = Symlink("link")
)
