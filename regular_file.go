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

// String implements [Transformer].
func (f File) String() string {
	return string(f)
}

// Abs implements [Transformer].
func (f File) Abs() (File, error) {
	return abs(f)
}

// Clean implements [Transformer].
func (f File) Clean() File {
	return clean(f)
}

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

// Rel implements [Transformer].
func (f File) Rel(base Dir) (File, error) {
	return rel(base, f)
}

// Beholder --------------------------------------------------------------------
var _ Beholder[File] = File("./example")

// Exists implements [Beholder].
func (f File) Exists() bool {
	return PathStr(f).Exists()
}

// Lstat implements [Beholder].
func (f File) Lstat() (Info[File], error) {
	// FIXME: note complication where Dir/File MAY also be a symlink
	return lstat(f)
}

// OnDisk implements [Beholder].
func (f File) OnDisk() (Info[File], error) {
	return lstat(f)
}

// Stat implements [Beholder].
func (f File) Stat() (Info[File], error) {
	return stat(f)
}

// Changer ----------------------------------------------------------------------
var _ Changer = File("./example")

// Chmod implements [Changer].
func (f File) Chmod(mode os.FileMode) error {
	return chmod(f, mode)
}

// Chown implements [Changer].
func (f File) Chown(uid int, gid int) error {
	return chown(f, uid, gid)
}

// Mover ------------------------------------------------------------------------
var _ Mover[File] = File("./example")

// Remove implements [Mover].
func (f File) Remove() error {
	return remove(f)
}

// Rename implements [Mover].
func (f File) Rename(newPath PathStr) (File, error) {
	return rename(f, newPath)
}

// Maker -----------------------------------------------------------------------
var _ Maker[FileHandle] = File("./example")

// Make implements [Maker].
func (f File) Make(perm fs.FileMode) (FileHandle, error) {
	return f.Open(os.O_RDWR|os.O_CREATE, perm)
}

// MakeAll implements [Maker].
func (f File) MakeAll(perm, parentPerm fs.FileMode) (result FileHandle, err error) {
	_, err = f.Parent().MakeAll(parentPerm, parentPerm)
	if err != nil {
		return
	}
	return f.Make(perm)
}

// Readable --------------------------------------------------------------------
var _ Readable[[]byte] = File("./example")

func (f File) Read() ([]byte, error) {
	return os.ReadFile(string(f))
}
