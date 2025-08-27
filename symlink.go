package pathlib

import (
	"errors"
	"io/fs"
	"os"
)

type Symlink PathStr

// See [os.Symlink].
func (s Symlink) LinkTo(target PathStr) (Symlink, error) {
	return s, os.Symlink(target.String(), s.String())
}

// TODO: iterator over chained symlinks

// Readable --------------------------------------------------------------------
var _ Readable[PathStr] = Symlink("./link")

// Returns the target of the symlink. Note that the target may not exist.
//
// See [os.Readlink].
//
// Read implements [Readable].
func (s Symlink) Read() (PathStr, error) {
	link, err := os.Readlink(string(s))
	return PathStr(link), err
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

// Ext implements [PurePath]
func (s Symlink) Ext() string {
	return ext(s)
}

// -----------------------------------------------------------------------------
var _ Transformer[Symlink] = Symlink("./link")

// Convenience method to cast get the untyped string representation of the path.
//
// String implements [Transformer].
func (s Symlink) String() string {
	return string(s)
}

// Returns true if the two paths represent the same path. This does not take into account any links.
//
// Eq implements [Transformer].
func (s Symlink) Eq(other Symlink) bool {
	return PathStr(s).Eq(PathStr(other))
}

// Remove ".", "..", and repeated slashes from a path.
//
// See [path/filepath.Clean].
//
// Clean implements [Transformer]
func (s Symlink) Clean() Symlink {
	return clean(s)
}

// Returns an absolute path, or an error if the path cannot be made absolute. Note that there may be more than one
// absolute path for a given input path.
//
// See [path/filepath.Abs].
//
// Abs implements [Transformer].
func (s Symlink) Abs() (Symlink, error) {
	return abs(s)
}

// Localize implements [Transformer].
func (s Symlink) Localize() (Symlink, error) {
	return localize(s)
}

// Returns a relative path to the target directory, or an error if the path cannot be made relative.
//
// See [path/filepath.Rel].
//
// Rel implements [Transformer].
func (s Symlink) Rel(base Dir) (Symlink, error) {
	return rel(base, s)
}

// ExpandUser implements [Transformer]
func (s Symlink) ExpandUser() (Symlink, error) {
	return expandUser(s)
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Symlink] = Symlink("./link")

// Returns true if the path exists on-disk WITHOUT following symlinks.
//
// See [os.Stat], [fs.ErrNotExist].
//
// Exists implements [Beholder].
func (s Symlink) Exists() bool {
	_, err := lstat(s)
	return !errors.Is(err, fs.ErrNotExist)
}

// Looks up the symlink's info on-disk. Note that this returns information about the
// symlink itself, not its target. If the file info's mode not match [fs.ModeSymlink],
// Lstat returns a [WrongTypeOnDisk] error.
//
// See [os.Lstat].
//
// Lstat implements [Beholder].
func (s Symlink) Lstat() (result Info[Symlink], err error) {
	result, err = lstat(s)
	if err == nil && ((result.Mode() & fs.ModeSymlink) != fs.ModeSymlink) {
		err = WrongTypeOnDisk[Symlink]{result}
	}
	return
}

// Since Stat follows symlinks, it doesn't perform any validation of returned [Info]'s file mode.
//
// See [os.Stat].
//
// Stat implements [Beholder].
func (s Symlink) Stat() (Info[Symlink], error) {
	return stat(s)
}

// // https://go.dev/play/p/mWNvcZLrjog
// // https://godbolt.org/z/1caPfvzfh

// Changer ----------------------------------------------------------------------
var _ Changer = Symlink("./link")

// See [os.Chmod].
//
// Chmod implements [Changer].
func (s Symlink) Chmod(mode os.FileMode) error {
	return chmod(s, mode)
}

// See [os.Chown].
//
// Chown implements [Changer].
func (s Symlink) Chown(uid int, gid int) error {
	return chown(s, uid, gid)
}

// Mover ------------------------------------------------------------------------
var _ Remover[Symlink] = Symlink("./link")

// Remove the link without affecting the link target.
//
// See [os.Remove].
//
// Remove implements [Remover].
func (s Symlink) Remove() error {
	return remove(s)
}

// Rename the link without affecting the link target.
//
// see [os.Rename].
//
// Rename implements [Remover].
func (s Symlink) Rename(newPath PathStr) (Symlink, error) {
	return rename(s, newPath)
}
