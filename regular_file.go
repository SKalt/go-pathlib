package pathlib

import (
	"io/fs"
	"os"
)

// A path that represents a file.
type File PathStr

// See [os.OpenFile].
func (f File) Open(flag int, perm fs.FileMode) (FileHandle, error) {
	ptr, err := os.OpenFile(string(f), flag, perm)
	if err != nil {
		return nil, err
	}
	return &handle{ptr}, nil
}

// PurePath --------------------------------------------------------------------
var _ PurePath = File("./example.txt")

// BaseName implements [PurePath].
func (f File) BaseName() string {
	return baseName(f)
}

// Ext implements [PurePath].
func (f File) Ext() string {
	return ext(f)
}

// IsAbsolute implements [PurePath].
func (f File) IsAbsolute() bool {
	return isAbsolute(f)
}

// IsLocal implements [PurePath].
func (f File) IsLocal() bool {
	return isLocal(f)
}

// Join implements [PurePath].
func (f File) Join(parts ...string) PathStr {
	return join(f, parts...)
}

// Parent implements [PurePath].
func (f File) Parent() Dir {
	return parent(f)
}

// Parts implements [PurePath].
func (f File) Parts() []string {
	return PathStr(f).Parts()
}

// Transformer -----------------------------------------------------------------
var _ Transformer[File] = File("./example")

// Convenience method to cast get the untyped string representation of the path.
//
// String implements [Transformer].
func (f File) String() string {
	return string(f)
}

// Returns an absolute path, or an error if the path cannot be made absolute. Note that there may be more than one
// absolute path for a given input path.
//
// See [path/filepath.Abs].
//
// Abs implements [Transformer].
func (f File) Abs() (File, error) {
	return abs(f)
}

// Remove ".", "..", and repeated slashes from a path.
//
// See [path/filepath.Clean].
//
// Clean implements [Transformer].
func (f File) Clean() File {
	return clean(f)
}

// Returns true if the two paths represent the same path.
//
// Eq implements [Transformer].
func (f File) Eq(other File) bool {
	return PathStr(f).Eq(PathStr(other))
}

// ExpandUser implements [Transformer].
func (f File) ExpandUser() (File, error) {
	return expandUser(f)
}

// Localize implements [Transformer].
func (f File) Localize() (File, error) {
	return localize(f)
}

// Returns a relative path to the target directory, or an error if the path cannot be made relative.
//
// See [path/filepath.Rel].
//
// Rel implements [Transformer].
func (f File) Rel(base Dir) (File, error) {
	return rel(base, f)
}

// Beholder --------------------------------------------------------------------
var _ Beholder[File] = File("./example")

// Returns true if the path exists on-disk after following symlinks.
//
// See [os.Stat], [fs.ErrNotExist].
//
// Exists implements [Beholder].
func (f File) Exists() bool {
	return PathStr(f).Exists()
}

// Observe the file info of the path on-disk. Does not follow symlinks.
// If the observed info is not a file or a symlink, Lstat returns a [WrongTypeOnDisk] error.
//
// See [os.Lstat].
//
// OnDisk implements [Beholder]
func (f File) Lstat() (info Info[File], err error) {
	info, err = lstat(f)
	if err == nil && !info.Mode().IsRegular() &&
		info.Mode()&fs.ModeSymlink != fs.ModeSymlink {
		err = WrongTypeOnDisk[File]{info}
	}
	return
}

// Observe the file info of the path on-disk. Follows symlinks.
// If the observed info is not a file Stat returns a [WrongTypeOnDisk] error.
//
// See [os.Stat].
//
// Stat implements [Beholder].
func (f File) Stat() (Info[File], error) {
	info, err := stat(f)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() && info.Mode()&fs.ModeSymlink != fs.ModeSymlink {
		return nil, WrongTypeOnDisk[File]{info}
	}
	return info, nil
}

// Changer ----------------------------------------------------------------------
var _ Changer = File("./example")

// See [os.Chmod].
//
// Chmod implements [Changer].
func (f File) Chmod(mode os.FileMode) error {
	return chmod(f, mode)
}

// See [os.Chown].
//
// Chown implements [Changer].
func (f File) Chown(uid int, gid int) error {
	return chown(f, uid, gid)
}

// Mover ------------------------------------------------------------------------
var _ Remover[File] = File("./example")

// See [os.Remove].
//
// Remove implements [Remover].
func (f File) Remove() error {
	return remove(f)
}

// See [os.Rename].
//
// Rename implements [Remover].
func (f File) Rename(newPath PathStr) (File, error) {
	return rename(f, newPath)
}

// Maker -----------------------------------------------------------------------
var _ Maker[FileHandle] = File("./example")

// Create the file if it doesn't exist. If it does, do nothing.
//
// See [os.Open], [os.O_CREATE].
//
// Make implements [Maker].
func (f File) Make(perm fs.FileMode) (FileHandle, error) {
	return f.Open(os.O_RDWR|os.O_CREATE, perm)
}

// Create the file and any missing parents. If the file exists, do nothing.
//
// See [os.Open], [os.O_CREATE].
//
// Make implements [Maker].
func (f File) MakeAll(perm, parentPerm fs.FileMode) (result FileHandle, err error) {
	_, err = f.Parent().MakeAll(parentPerm, parentPerm)
	if err != nil {
		return
	}
	return f.Make(perm)
}

// Readable --------------------------------------------------------------------
var _ Readable[[]byte] = File("./example")

// See [os.ReadFile].
//
// Read implements [Readable].
func (f File) Read() ([]byte, error) {
	return os.ReadFile(string(f))
}
