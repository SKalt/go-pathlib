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
	return PathStr(s).BaseName()
}

// IsAbsolute implements [PurePath].
func (s Symlink) IsAbsolute() bool {
	return PathStr(s).IsAbsolute()
}

// IsLocal implements [PurePath].
func (s Symlink) IsLocal() bool {
	return PathStr(s).IsLocal()
}

// Join implements [PurePath].
func (s Symlink) Join(parts ...string) PathStr {
	return PathStr(s).Join(parts...)
}

func (s Symlink) Parts() []string {
	return PathStr(s).Parts()
}

// Parent implements [PurePath].
func (s Symlink) Parent() Dir {
	return PathStr(s).Parent()
}

// -----------------------------------------------------------------------------
var _ Transformer[Symlink] = Symlink("./link")
var _ InfallibleTransformer[Symlink] = Symlink("./link")

func (s Symlink) Eq(other Symlink) bool {
	return PathStr(s).Eq(PathStr(other))
}

func (s Symlink) Clean() Symlink {
	return Symlink(PathStr(s).Clean())
}

// Abs implements [Transformer].
func (s Symlink) Abs() (Symlink, error) {
	q, err := PathStr(s).Abs()
	return Symlink(q), err
}

// Localize implements [Transformer].
func (s Symlink) Localize() (Symlink, error) {
	q, err := PathStr(s).Localize()
	return Symlink(q), err
}

// Rel implements [Transformer].
func (s Symlink) Rel(target Dir) (Symlink, error) {
	result, err := PathStr(s).Rel(target)
	return Symlink(result), err
}

func (s Symlink) ExpandUser() (Symlink, error) {
	q, err := PathStr(s).ExpandUser()
	return Symlink(q), err
}

func (s Symlink) Ext() string {
	return PathStr(s).Ext()
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
func (s Symlink) OnDisk() (OnDisk[Symlink], error) {
	actual, err := PathStr(s).OnDisk()
	if err != nil {
		return nil, err
	}
	if !isSymLink(actual.Mode()) {
		return nil, WrongTypeOnDisk[Symlink]{actual}
	}
	return onDisk[Symlink]{actual}, nil
}

// Exists implements [Beholder].
func (s Symlink) Exists() bool {
	panic("unimplemented")
}

// Lstat implements [Beholder].
func (s Symlink) Lstat() (OnDisk[Symlink], error) {
	panic("unimplemented")
}

// Stat implements [Beholder].
func (s Symlink) Stat() (OnDisk[Symlink], error) {
	panic("unimplemented")
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
	result, err := PathStr(s).Chmod(mode)
	return Symlink(result), err
}

// Chown implements [Manipulator].
func (s Symlink) Chown(uid int, gid int) (Symlink, error) {
	result, err := PathStr(s).Chown(uid, gid)
	return Symlink(result), err
}

// Remove implements [Manipulator].
func (s Symlink) Remove() error {
	return os.Remove(string(s))
}

// Rename implements [Manipulator].
func (s Symlink) Rename(newPath PathStr) (Symlink, error) {
	result, err := PathStr(s).Rename(newPath)
	return Symlink(result), err
}

// MustChmod implements [InfallibleManipulator].
func (s Symlink) MustChmod(mode os.FileMode) Symlink {
	return expect(s.Chmod(mode))
}

// MustChown implements [InfallibleManipulator].
func (s Symlink) MustChown(uid int, gid int) Symlink {
	return expect(s.Chown(uid, gid))
}

// MustRemove implements [InfallibleManipulator].
func (s Symlink) MustRemove() {
	if err := s.Remove(); err != nil {
		panic(err)
	}
}

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

func (s Symlink) MustRemoveAll() {
	if err := os.RemoveAll(string(s)); err != nil {
		panic(err)
	}
}
