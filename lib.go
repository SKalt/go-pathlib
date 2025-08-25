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
	// Returns an absolute path, or an error if the path cannot be made absolute. Note that there may be more than one
	// absolute path for a given input path.
	//
	// See [path/filepath.Abs].
	Abs() (P, error)
	// Returns a relative path to the target directory, or an error if the path cannot be made relative.
	//
	// See [path/filepath.Rel].
	Rel(target Dir) (P, error)
	// See [path/filepath.Localize].
	Localize() (P, error)
	// Expand `~` into the home directory of the current user.
	ExpandUser() (P, error)
	// Remove ".", "..", and repeated slashes from a path.
	//
	// See [path/filepath.Clean].
	Clean() P
	// Returns true if the two paths represent the same path.
	Eq(other P) bool

	// Convenience method to cast get the untyped string representation of the path.
	String() string
}

// An observation of a path on-disk, including a constant observation timestamp.
type Info[P Kind] interface {
	fs.FileInfo
	PurePath
	Transformer[P]
	Changer
	Remover[P]
	// the typed version of [fs.FileInfo.Name]
	Path() P
}

// Behaviors for inspecting a path on-disk.
type Beholder[P Kind] interface {
	// Observe the file info of the path on-disk.
	OnDisk() (Info[P], error)
	// See [os.Stat].
	Stat() (Info[P], error)
	// See [os.Lstat].
	Lstat() (Info[P], error)
	// Returns true if the path exists on-disk.
	Exists() bool
}

type Maker[T any] interface {
	Make(perm fs.FileMode) (T, error)
	MakeAll(perm, parentPerm fs.FileMode) (T, error)
}

// Behaviors that cause something at a path to no longer be there.
type Remover[P Kind] interface {
	// see [os.Remove].
	Remove() error
	// see [os.Rename].
	Rename(newPath PathStr) (P, error)
}

// Methods that alter a filesystem object without changing its path or performing I/O.
type Changer interface {
	// see [os.Chmod].
	Chmod(fs.FileMode) error
	// see [os.Chown].
	Chown(uid, gid int) error
}
