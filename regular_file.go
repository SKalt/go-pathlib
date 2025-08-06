package pathlib

import (
	"io/fs"
	"os"
)

// A path that represents a file.
type File PathStr

// See [os.OpenFile].
func (f File) Open(flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(string(f), flag, perm)
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
var _ InfallibleTransformer[File] = File("./example")

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
func (f File) Rel(target Dir) (File, error) {
	return rel(f, target)
}

// MustExpandUser implements [InfallibleTransformer].
func (f File) MustExpandUser() File {
	return expect(f.ExpandUser())
}

// MustLocalize implements [InfallibleTransformer].
func (f File) MustLocalize() File {
	return expect(f.Localize())
}

// MustMakeAbs implements [InfallibleTransformer].
func (f File) MustMakeAbs() File {
	return expect(f.Abs())
}

// MustMakeRel implements [InfallibleTransformer].
func (f File) MustMakeRel(target Dir) File {
	return expect(f.Rel(target))
}

// Beholder --------------------------------------------------------------------
var _ Beholder[File] = File("./example")
var _ InfallibleBeholder[File] = File("./example")

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

// MustBeOnDisk implements [InfallibleBeholder].
func (f File) MustBeOnDisk() OnDisk[File] {
	return expect(f.OnDisk())
}

// MustLstat implements [InfallibleBeholder].
func (f File) MustLstat() OnDisk[File] {
	return expect(f.Lstat())
}

// MustStat implements [InfallibleBeholder].
func (f File) MustStat() OnDisk[File] {
	return expect(f.Stat())
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[File] = File("./example")
var _ InfallibleManipulator[File] = File("./example")

// Chmod implements [Manipulator].
func (f File) Chmod(mode os.FileMode) (File, error) {
	return chmod(f, mode)
}

// Chown implements [Manipulator].
func (f File) Chown(uid int, gid int) (File, error) {
	return chown(f, uid, gid)
}

// Remove implements [Manipulator].
func (f File) Remove() error {
	return os.Remove(string(f))
}

// Rename implements [Manipulator].
func (f File) Rename(newPath PathStr) (File, error) {
	return rename(f, newPath)
}

// MustChmod implements [InfallibleManipulator].
func (f File) MustChmod(perm os.FileMode) File {
	return expect(f.Chmod(perm))
}

// MustChown implements [InfallibleManipulator].
func (f File) MustChown(uid int, gid int) File {
	return expect(f.Chown(uid, gid))
}

// MustRemove implements [InfallibleManipulator].
func (f File) MustRemove() {
	expect[any](nil, f.Remove())
}

// MustRename implements [InfallibleManipulator].
func (f File) MustRename(newPath PathStr) File {
	return expect(f.Rename(newPath))
}

// Maker -----------------------------------------------------------------------
var _ Maker[*os.File] = File("./example")
var _ InfallibleMaker[*os.File] = File("./example")

// Make implements [Maker].
func (f File) Make(perm fs.FileMode) (handle *os.File, err error) {
	return f.Open(os.O_RDWR|os.O_CREATE, perm)
}

// MakeAll implements [Maker].
func (f File) MakeAll(perm, parentPerm fs.FileMode) (handle *os.File, err error) {
	_, err = f.Parent().MakeAll(parentPerm, parentPerm)
	if err != nil {
		return
	}
	return f.Make(perm)
}

// Panics if [File.Make] returns an error.
//
// MustMake implements [Maker].
func (f File) MustMake(perm fs.FileMode) *os.File {
	return expect(f.Make(perm))
}

// Panics if [File.Make] returns an error.
//
// MustMake implements [Maker].
func (f File) MustMakeAll(perm, parentPerm fs.FileMode) *os.File {
	return expect(f.MakeAll(perm, parentPerm))
}

// Readable --------------------------------------------------------------------
var _ Readable[[]byte] = File("./example")
var _ InfallibleReader[[]byte] = File("./example")

func (f File) Read() (data []byte, err error) {
	return os.ReadFile(string(f))
}

func (f File) MustRead() []byte {
	return expect(f.Read())
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = File("./example")
var _ InfallibleDestroyer = File("./example")

// MustRemoveAll implements [InfallibleDestroyer].
func (f File) MustRemoveAll() {
	PathStr(f).MustRemoveAll()
}

// RemoveAll implements [Destroyer].
func (f File) RemoveAll() error {
	return os.RemoveAll(string(f))
}
