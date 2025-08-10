package pathlib

import (
	"io/fs"
	"os"
)

// A path that represents a file.
type File PathStr

// See [os.OpenFile].
func (f File) Open(flag int, perm fs.FileMode) Result[*os.File] {
	handle, err := os.OpenFile(string(f), flag, perm)
	return Result[*os.File]{handle, err}
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
func (f File) Abs() Result[File] {
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
func (f File) ExpandUser() Result[File] {
	return expandUser(f)
}

// Localize implements [Transformer].
func (f File) Localize() Result[File] {
	return localize(f)
}

// Rel implements [Transformer].
func (f File) Rel(base Dir) Result[File] {
	return rel(base, f)
}

// Beholder --------------------------------------------------------------------
var _ Beholder[File] = File("./example")

// Exists implements [Beholder].
func (f File) Exists() bool {
	return PathStr(f).Exists()
}

// Lstat implements [Beholder].
func (f File) Lstat() Result[OnDisk[File]] {
	return lstat(f)
}

// OnDisk implements [Beholder].
func (f File) OnDisk() Result[OnDisk[File]] {
	return lstat(f)
}

// Stat implements [Beholder].
func (f File) Stat() Result[OnDisk[File]] {
	return stat(f)
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[File] = File("./example")

// Chmod implements [Manipulator].
func (f File) Chmod(mode os.FileMode) Result[File] {
	return chmod(f, mode)
}

// Chown implements [Manipulator].
func (f File) Chown(uid int, gid int) Result[File] {
	return chown(f, uid, gid)
}

// Remove implements [Manipulator].
func (f File) Remove() error {
	return os.Remove(string(f))
}

// Rename implements [Manipulator].
func (f File) Rename(newPath PathStr) Result[File] {
	return rename(f, newPath)
}

// Maker -----------------------------------------------------------------------
var _ Maker[*os.File] = File("./example")

// Make implements [Maker].
func (f File) Make(perm fs.FileMode) Result[*os.File] {
	return f.Open(os.O_RDWR|os.O_CREATE, perm)
}

// MakeAll implements [Maker].
func (f File) MakeAll(perm, parentPerm fs.FileMode) (result Result[*os.File]) {
	result.Err = f.Parent().MakeAll(parentPerm, parentPerm).Err
	if !result.IsOk() {
		return
	}
	result = f.Make(perm)
	return
}

// Readable --------------------------------------------------------------------
var _ Readable[[]byte] = File("./example")

func (f File) Read() Result[[]byte] {
	data, err := os.ReadFile(string(f))
	return Result[[]byte]{data, err}
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = File("./example")

// RemoveAll implements [Destroyer].
func (f File) RemoveAll() error {
	return os.RemoveAll(string(f))
}
