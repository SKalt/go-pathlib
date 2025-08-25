package pathlib

import (
	"io/fs"
	"os"
)

type Symlink PathStr

// Readable --------------------------------------------------------------------
var _ Readable[PathStr] = Symlink("./link")

// Read implements [Readable].
func (s Symlink) Read() (PathStr, error) {
	link, err := os.Readlink(string(s))
	return PathStr(link), err
}

// See [os.Symlink]
func (s Symlink) LinkTo(target PathStr) (Symlink, error) {
	return s, os.Symlink(target.String(), s.String())
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
func (s Symlink) Abs() (Symlink, error) {
	return abs(s)
}

// Localize implements [Transformer].
func (s Symlink) Localize() (Symlink, error) {
	return localize(s)
}

// Rel implements [Transformer].
func (s Symlink) Rel(base Dir) (Symlink, error) {
	return rel(base, s)
}

func (s Symlink) ExpandUser() (Symlink, error) {
	return expandUser(s)
}

func (s Symlink) Ext() string {
	return ext(s)
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Symlink] = Symlink("./link")

// OnDisk implements [Beholder].
func (s Symlink) OnDisk() (result Info[Symlink], err error) {
	return s.Lstat()
}

// Exists implements [Beholder].
func (s Symlink) Exists() bool {
	return PathStr(s).Exists()
}

// Lstat implements [Beholder].
func (s Symlink) Lstat() (result Info[Symlink], err error) {
	result, err = lstat(s)
	if err == nil && ((result.Mode() & fs.ModeSymlink) != fs.ModeSymlink) {
		err = WrongTypeOnDisk[Symlink]{result}
	}
	return
}

// Stat implements [Beholder].
func (s Symlink) Stat() (Info[Symlink], error) {
	return stat(s)
}

// // https://go.dev/play/p/mWNvcZLrjog
// // https://godbolt.org/z/1caPfvzfh

// Changer ----------------------------------------------------------------------
var _ Changer = Symlink("./link")

// Chmod implements [Changer].
func (s Symlink) Chmod(mode os.FileMode) error {
	return chmod(s, mode)
}

// Chown implements [Changer].
func (s Symlink) Chown(uid int, gid int) error {
	return chown(s, uid, gid)
}

// Mover ------------------------------------------------------------------------
var _ Remover[Symlink] = Symlink("./link")

// Remove implements [Remover].
func (s Symlink) Remove() error {
	return remove(s)
}

// Rename implements [Remover].
func (s Symlink) Rename(newPath PathStr) (Symlink, error) {
	return rename(s, newPath)
}
