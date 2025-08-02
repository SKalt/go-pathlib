package pathlib

import (
	"io/fs"
	"os"
)

type Kind interface {
	PurePath
	~string
}

type Readable[T any] interface {
	Read() (T, error)
}
type InfallibleReader[T any] interface {
	MustRead() T
}

// String-only infallible path operations that do not require filesystem access or syscalls.
type PurePath interface {
	// See [path/filepath.Join].
	Join(segments ...string) PathStr
	// Return the parent directory. Should have the same properties as [path/filepath.Dir].
	Parent() Dir
	// See [path/filepath.Base].
	BaseName() string
	// See [path/filepath.Ext].
	Ext() string
	// Split the path into multiple non-empty segments.
	Parts() []string

	// Returns true if the path is absolute. See [path/filepath.IsAbs].
	IsAbsolute() bool
	// Returns true if the path is local/relative. See [path/filepath.IsLocal].
	IsLocal() bool
}

// transforms the appearance of a path, but not what it represents.
type Transformer[Self Kind] interface {
	// See [path/filepath.Abs].
	// Returns an absolute path, or an error if the path cannot be made absolute.
	Abs() (Self, error)
	// See [path/filepath.Rel].
	// Returns a relative path to the target directory, or an error if the path cannot be made relative.
	Rel(target Dir) (Self, error)
	// See [path/filepath.Localize].
	Localize() (Self, error)
	// Expand `~` into the home directory of the current user.
	ExpandUser() (Self, error)
	// See [path/filepath.Clean].
	Clean() Self
	// Returns true if the two paths represent the same path.
	Eq(other Self) bool
}
type InfallibleTransformer[Self Kind] interface {
	// Makes the path absolute or panics if it cannot.
	MustMakeAbs() Self
	// Makes the path relative to the target directory or panics if it cannot.
	MustMakeRel(target Dir) Self
	// Localizes the path or panics if it cannot. See [path/filepath.Localize].
	MustLocalize() Self
	// Expands a leading tilde into the home directory. Panics if [os.UserHomeDir] is unable to resolve the user's home directory.
	MustExpandUser() Self
}

// An observation of a path on-disk, including a constant observation timestamp.
type OnDisk[PathKind Kind] interface {
	fs.FileInfo
	PurePath
	Transformer[PathKind]
}

type Beholder[PathKind Kind] interface {
	// Observe the file info of the path on-disk. Does not follow symlinks. See [os.Lstat].
	OnDisk() (OnDisk[PathKind], error)
	// See [os.Stat].
	Stat() (OnDisk[PathKind], error)
	// See [os.Lstat].
	Lstat() (OnDisk[PathKind], error)
	// Returns true if the path exists on-disk.
	Exists() bool
}
type InfallibleBeholder[PathKind Kind] interface {
	// Panics if the path does not exist on-disk.
	MustBeOnDisk() OnDisk[PathKind]
	// Panics if [os.Stat] fails.
	MustStat() OnDisk[PathKind]
	// Panics if [os.Lstat] fails.
	MustLstat() OnDisk[PathKind]
}

type Maker[T any] interface {
	Make(perm fs.FileMode) (T, error)
	MakeAll(perm, parentPerm fs.FileMode) (T, error)
}

type InfallibleMaker[T any] interface {
	// Panics if Make fails.
	MustMake(perm fs.FileMode) T
	MakeAll(perm, parentPerm fs.FileMode) (T, error)
}

type Manipulator[PathKind Kind] interface {
	// see [os.Remove].
	Remove() error
	// see [os.Chmod].
	Chmod(os.FileMode) (PathKind, error)
	// see [os.Chown].
	Chown(uid, gid int) (PathKind, error)
	// see [os.Rename].
	Rename(newPath PathStr) (PathKind, error)
}
type InfallibleManipulator[PathKind Kind] interface {
	// see [os.Remove]. Panics if Remove fails.
	MustRemove()
	// see [os.Chmod]. Panics if Chmod fails.
	MustChmod(mode os.FileMode) PathKind
	// see [os.Chown]. Panics if Chown fails.
	MustChown(uid, gid int) PathKind
	// see [os.Rename]. Panics if Rename fails.
	MustRename(newPath PathStr) PathKind
}

type Destroyer interface {
	// see [os.RemoveAll].
	RemoveAll() error
}
type InfallibleDestroyer interface {
	// Panics if [os.RemoveAll] fails.
	MustRemoveAll()
}
