package pathlib

import (
	"os"
)

type Symlink PathStr

// Readable --------------------------------------------------------------------
var _ Readable[PathStr] = Symlink("./link")

// Read implements [Readable].
func (s Symlink) Read() Result[PathStr] {
	link, err := os.Readlink(string(s))
	return Result[PathStr]{PathStr(link), err}
}

// See [os.Symlink]
func (s Symlink) LinkTo(target PathStr) Result[Symlink] {
	return Result[Symlink]{s, os.Symlink(target.String(), s.String())}
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

// String implements [Transformer].
func (s Symlink) String() string {
	return string(s)
}

func (s Symlink) Eq(other Symlink) bool {
	return PathStr(s).Eq(PathStr(other))
}

func (s Symlink) Clean() Symlink {
	return clean(s)
}

// Abs implements [Transformer].
func (s Symlink) Abs() Result[Symlink] {
	return abs(s)
}

// Localize implements [Transformer].
func (s Symlink) Localize() Result[Symlink] {
	return localize(s)
}

// Rel implements [Transformer].
func (s Symlink) Rel(target Dir) Result[Symlink] {
	return rel(s, target)
}

func (s Symlink) ExpandUser() Result[Symlink] {
	return expandUser(s)
}

func (s Symlink) Ext() string {
	return ext(s)
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Symlink] = Symlink("./link")

// OnDisk implements [Beholder].
func (s Symlink) OnDisk() (result Result[OnDisk[Symlink]]) {
	result = lstat(s)
	if result.IsOk() && !isSymLink(result.Val.Mode()) {
		result.Err = WrongTypeOnDisk[Symlink]{result.Val}
	}
	return
}

// Exists implements [Beholder].
func (s Symlink) Exists() bool {
	return PathStr(s).Exists()
}

// Lstat implements [Beholder].
func (s Symlink) Lstat() Result[OnDisk[Symlink]] {
	return lstat(s)
}

// Stat implements [Beholder].
func (s Symlink) Stat() Result[OnDisk[Symlink]] {
	return stat(s)
}

// // https://go.dev/play/p/mWNvcZLrjog
// // https://godbolt.org/z/1caPfvzfh

// Manipulator -----------------------------------------------------------------
var _ Manipulator[Symlink] = Symlink("./link")

// Chmod implements [Manipulator].
func (s Symlink) Chmod(mode os.FileMode) Result[Symlink] {
	return chmod(s, mode)
}

// Chown implements [Manipulator].
func (s Symlink) Chown(uid int, gid int) Result[Symlink] {
	return chown(s, uid, gid)
}

// Remove implements [Manipulator].
func (s Symlink) Remove() error {
	return os.Remove(string(s))
}

// Rename implements [Manipulator].
func (s Symlink) Rename(newPath PathStr) Result[Symlink] {
	return rename(s, newPath)
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = Symlink("./link")

// RemoveAll implements [Destroyer].
func (s Symlink) RemoveAll() error {
	return os.RemoveAll(string(s))
}
