package pathlib

import (
	"io/fs"
	"os"
)

// A path that represents a file.
type File PathStr

// See [os.OpenFile].
func (f File) Open(flag int, perm fs.FileMode) (*Handle, error) {
	ptr, err := os.OpenFile(string(f), os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return nil, err
	}
	return &Handle{ptr}, nil
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
func (f File) Lstat() (OnDisk[File], error) {
	return lstat(f)
}

// OnDisk implements [Beholder].
func (f File) OnDisk() (OnDisk[File], error) {
	return lstat(f)
}

// Stat implements [Beholder].
func (f File) Stat() (OnDisk[File], error) {
	return stat(f)
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[File] = File("./example")

// Chmod implements [Manipulator].
func (f File) Chmod(mode os.FileMode) (File, error) {
	return chmod(f, mode)
}

// Chown implements [Manipulator].
func (f File) Chown(uid int, gid int) (File, error) {
	return chown(f, uid, gid)
}

// Remove implements [Manipulator].
func (f File) Remove() (File, error) {
	return remove(f)
}

// Rename implements [Manipulator].
func (f File) Rename(newPath PathStr) (File, error) {
	return rename(f, newPath)
}

// Maker -----------------------------------------------------------------------
var _ Maker[*Handle] = File("./example")

// Make implements [Maker].
func (f File) Make(perm fs.FileMode) (*Handle, error) {
	return f.Open(os.O_RDWR|os.O_CREATE, perm)
}

// MakeAll implements [Maker].
func (f File) MakeAll(perm, parentPerm fs.FileMode) (result *Handle, err error) {
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

// FIXME: file handle wrapper type
