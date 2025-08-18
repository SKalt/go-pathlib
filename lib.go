package pathlib

import (
	"io/fs"
)

// Any type constraint: any string type that represents a path
type Kind interface {
	PurePath
	~string
}

type Readable[T any] interface {
	Read() (T, error)
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
type Transformer[P Kind] interface {
	// See [path/filepath.Abs].
	// Returns an absolute path, or an error if the path cannot be made absolute.
	Abs() (P, error)
	// See [path/filepath.Rel].
	// Returns a relative path to the target directory, or an error if the path cannot be made relative.
	Rel(target Dir) (P, error)
	// See [path/filepath.Localize].
	Localize() (P, error)
	// Expand `~` into the home directory of the current user.
	ExpandUser() (P, error)
	// See [path/filepath.Clean].
	Clean() P
	// Returns true if the two paths represent the same path.
	Eq(other P) bool
	String() string
}

// An observation of a path on-disk, including a constant observation timestamp.
type OnDisk[PathKind Kind] interface {
	fs.FileInfo
	PurePath
	Transformer[PathKind]
	Manipulator[PathKind]
	// the typed version of [fs.FileInfo.Name]
	Path() PathKind
}

// Behaviors for inspecting a path on-disk.
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

type Maker[T any] interface {
	Make(perm fs.FileMode) (T, error)
	MakeAll(perm, parentPerm fs.FileMode) (T, error)
}

type Manipulator[P Kind] interface {
	// see [os.Remove].
	Remove() (P, error)
	// see [os.Chmod].
	Chmod(fs.FileMode) (P, error)
	// see [os.Chown].
	Chown(uid, gid int) (P, error)
	// see [os.Rename].
	Rename(newPath PathStr) (P, error)
}

type Destroyer[P Kind] interface {
	// see [os.RemoveAll].
	RemoveAll() (P, error)
}
