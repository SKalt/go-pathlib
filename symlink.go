package pathlib

import (
	"os"
)

type Symlink PathStr

// Readable --------------------------------------------------------------------
var _ Readable[PathStr] = Symlink("./link")
var _ InfallibleReader[PathStr] = Symlink("./link")

// Read implements [Readable].
func (s Symlink) Read() (PathStr, error) {
	link, err := os.Readlink(string(s))
	return PathStr(link), err
}

// MustRead implements [InfallibleReader].
func (s Symlink) MustRead() PathStr {
	return expect(s.Read())
}

// -----------------------------------------------------------------------------
var _ PurePath = Symlink("./link")

// BaseName implements [PurePath].
func (s Symlink) BaseName() string {
	return baseName(s)
}

// IsAbsolute implements [PurePath].
func (s Symlink) IsAbsolute() bool {
	return isAbsolute(s)
}

// IsLocal implements [PurePath].
func (s Symlink) IsLocal() bool {
	return isLocal(s)
}

// Join implements [PurePath].
func (s Symlink) Join(parts ...string) PathStr {
	return join(s, parts...)
}

func (s Symlink) Parts() []string {
	return PathStr(s).Parts()
}

// Parent implements [PurePath].
func (s Symlink) Parent() Dir {
	return parent(s)
}

// -----------------------------------------------------------------------------
var _ Transformer[Symlink] = Symlink("./link")
var _ InfallibleTransformer[Symlink] = Symlink("./link")

func (s Symlink) Eq(other Symlink) bool {
	return PathStr(s).Eq(PathStr(other))
}

func (s Symlink) Clean() Symlink {
	return clean(s)
}

// Abs implements [Transformer].
func (s Symlink) Abs() (Symlink, error) {
	return abs(s)
}

// Localize implements [Transformer].
func (s Symlink) Localize() (Symlink, error) {
	return localize(s)
}

// Rel implements [Transformer].
func (s Symlink) Rel(target Dir) (Symlink, error) {
	return rel(s, target)
}

func (s Symlink) ExpandUser() (Symlink, error) {
	return expandUser(s)
}

func (s Symlink) Ext() string {
	return ext(s)
}

// MustExpandUser implements [InfallibleTransformer].
func (s Symlink) MustExpandUser() Symlink {
	return expect(s.ExpandUser())
}

// MustLocalize implements [InfallibleTransformer].
func (s Symlink) MustLocalize() Symlink {
	return expect(s.Localize())
}

// MustMakeAbs implements [InfallibleTransformer].
func (s Symlink) MustMakeAbs() Symlink {
	return expect(s.Abs())
}

// MustMakeRel implements [InfallibleTransformer].
func (s Symlink) MustMakeRel(target Dir) Symlink {
	return expect(s.Rel(target))
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Symlink] = Symlink("./link")
var _ InfallibleBeholder[Symlink] = Symlink("./link")

// OnDisk implements [Beholder].
func (s Symlink) OnDisk() (result OnDisk[Symlink], err error) {
	result, err = lstat(s)
	if err != nil {
		return
	}
	if !isSymLink(result.Mode()) {
		err = WrongTypeOnDisk[Symlink]{result}
		result = nil
	}
	return
}

// Exists implements [Beholder].
func (s Symlink) Exists() bool {
	return PathStr(s).Exists()
}

// Lstat implements [Beholder].
func (s Symlink) Lstat() (OnDisk[Symlink], error) {
	return lstat(s)
}

// Stat implements [Beholder].
func (s Symlink) Stat() (OnDisk[Symlink], error) {
	return stat(s)
}

// MustBeOnDisk implements [InfallibleBeholder].
func (s Symlink) MustBeOnDisk() OnDisk[Symlink] {
	return expect(s.OnDisk())
}

// MustLstat implements [InfallibleBeholder].
func (s Symlink) MustLstat() OnDisk[Symlink] {
	return expect(s.Lstat())
}

// MustStat implements [InfallibleBeholder].
func (s Symlink) MustStat() OnDisk[Symlink] {
	return expect(s.Stat())
}

// // https://go.dev/play/p/mWNvcZLrjog
// // https://godbolt.org/z/1caPfvzfh

// Manipulator -----------------------------------------------------------------
var _ Manipulator[Symlink] = Symlink("./link")
var _ InfallibleManipulator[Symlink] = Symlink("./link")

// Chmod implements [Manipulator].
func (s Symlink) Chmod(mode os.FileMode) (Symlink, error) {
	return chmod(s, mode)
}

// Chown implements [Manipulator].
func (s Symlink) Chown(uid int, gid int) (Symlink, error) {
	return chown(s, uid, gid)
}

// Remove implements [Manipulator].
func (s Symlink) Remove() error {
	return os.Remove(string(s))
}

// Rename implements [Manipulator].
func (s Symlink) Rename(newPath PathStr) (Symlink, error) {
	return rename(s, newPath)
}

// Panics if [Symlink.Chmod] returns an error.
//
// MustChmod implements [InfallibleManipulator].
func (s Symlink) MustChmod(mode os.FileMode) Symlink {
	return expect(s.Chmod(mode))
}

// Panics if [Symlink.Chown] returns an error.
//
// MustChown implements [InfallibleManipulator].
func (s Symlink) MustChown(uid int, gid int) Symlink {
	return expect(s.Chown(uid, gid))
}

// Panics if [Symlink.Remove] returns an error.
//
// MustRemove implements [InfallibleManipulator].
func (s Symlink) MustRemove() {
	if err := s.Remove(); err != nil {
		panic(err)
	}
}

// Panics if [Symlink.Rename] returns an error.
//
// MustRename implements [InfallibleManipulator].
func (s Symlink) MustRename(newPath PathStr) Symlink {
	return expect(s.Rename(newPath))
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = Symlink("./link")
var _ InfallibleDestroyer = Symlink("./link")

// RemoveAll implements [Destroyer].
func (s Symlink) RemoveAll() error {
	return os.RemoveAll(string(s))
}

// Panics if [os.RemoveAll] returns an error.
//
// MustRemoveAll implements [InfallibleDestroyer]
func (s Symlink) MustRemoveAll() {
	if err := os.RemoveAll(string(s)); err != nil {
		panic(err)
	}
}
