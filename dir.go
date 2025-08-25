package pathlib

import (
	"io/fs"
	"os"
	"path/filepath"
)

// A string that represents a directory. The directory may or may not exist on-disk,
// and the string may or may not end in an [os.PathSeparator].
type Dir PathStr

// See [path/filepath.WalkDir].
func (d Dir) Walk(
	callback func(path PathStr, d fs.DirEntry, err error) error,
) error {
	return filepath.WalkDir(string(d), func(path string, d fs.DirEntry, err error) error {
		return callback(PathStr(path), d, err)
	})
}

// TODO: walk -> iter.Seq[PathStr]?

// See [path/filepath.Glob].
func (d Dir) Glob(pattern string) ([]PathStr, error) {
	matches, err := filepath.Glob(filepath.Join(string(d), pattern))
	if err != nil {
		return nil, err
	}
	result := make([]PathStr, len(matches))
	for i, m := range matches {
		result[i] = PathStr(m)
	}
	return result, nil
}

// CHange DIRectory. See [os.Chdir].
func (d Dir) Chdir() (Dir, error) {
	return d, os.Chdir(string(d))
}

// Readable --------------------------------------------------------------------
var _ Readable[[]fs.DirEntry] = Dir(".")

// See [os.ReadDir].
//
// Read implements [Reader].
func (d Dir) Read() ([]fs.DirEntry, error) {
	return os.ReadDir(string(d))
}

// PurePath --------------------------------------------------------------------
var _ PurePath = Dir(".")

// See [filepath.Base].
//
// BaseName implements [PurePath].
func (d Dir) BaseName() string {
	return baseName(d)
}

// Returns true if the path is absolute, false otherwise.
// See [filepath.IsAbs] for more details.
//
// IsAbsolute implements [PurePath].
func (d Dir) IsAbsolute() bool {
	return isAbsolute(d)
}

// IsLocal implements [PurePath].
func (d Dir) IsLocal() bool {
	return isLocal(d)
}

// Join implements [PurePath].
func (d Dir) Join(parts ...string) PathStr {
	return join(d, parts...)
}

func (d Dir) Parts() []string {
	return PathStr(d).Parts()
}

// Parent implements [PurePath].
func (d Dir) Parent() Dir {
	return parent(d)
}

// Ext implements [PurePath].
func (d Dir) Ext() string {
	return ext(d)
}

// Transformer -----------------------------------------------------------------
var _ Transformer[Dir] = Dir(".")

// String implements [Transformer].
func (d Dir) String() string {
	return string(d)
}

// Abs implements [Transformer].
func (d Dir) Abs() (Dir, error) {
	return abs(d)
}

// Eq implements [Transformer].
func (d Dir) Eq(other Dir) (equivalent bool) {
	return PathStr(d).Eq(PathStr(other))
}

// Clean implements [Transformer].
func (d Dir) Clean() Dir {
	return clean(d)
}

// See [path/filepath.Localize]
//
// Localize implements [Transformer].
func (d Dir) Localize() (Dir, error) {
	return localize(d)
}

// See [path/filepath.Rel]
//
// Rel implements [Transformer].
func (d Dir) Rel(base Dir) (Dir, error) {
	return rel(base, d.Clean())
}

// ExpandUser implements [Transformer].
func (d Dir) ExpandUser() (Dir, error) {
	return expandUser(d)
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Dir] = Dir(".")

// OnDisk implements [Beholder]
func (d Dir) OnDisk() (result Info[Dir], err error) {
	return d.Stat()
}

// Exists implements [Beholder].
func (d Dir) Exists() bool {
	_, err := d.Lstat()
	return err == nil
}

// See [os.Lstat].
//
// Lstat implements [Beholder].
func (d Dir) Lstat() (result Info[Dir], err error) {
	result, err = lstat(d)
	if err == nil && !result.IsDir() {
		err = WrongTypeOnDisk[Dir]{result}
	}
	return
}

// See [os.Stat].
//
// Stat implements [Beholder].
func (d Dir) Stat() (result Info[Dir], err error) {
	result, err = stat(d)
	if err == nil && !result.IsDir() {
		err = WrongTypeOnDisk[Dir]{result}
	}
	return
}

// Maker -----------------------------------------------------------------------
var _ Maker[Dir] = Dir("/example")

// Make implements [Maker].
func (d Dir) Make(perm fs.FileMode) (result Dir, err error) {
	return d, os.Mkdir(string(d), perm)
}

// MakeAll implements [Maker]
func (d Dir) MakeAll(perm, parentPerm fs.FileMode) (result Dir, err error) {
	result = d
	if d.Exists() {
		return
	}
	_, err = d.Parent().MakeAll(parentPerm, parentPerm)
	if err != nil {
		return
	}
	err = os.MkdirAll(string(d), perm)
	return
}

// Changer ----------------------------------------------------------------------
var _ Changer = Dir(".")

// See [os.Chmod].
// Chmod implements [Changer].
func (d Dir) Chmod(mode os.FileMode) error {
	return chmod(d, mode)
}

// See [os.Chown].
//
// Chown implements [Changer].
func (d Dir) Chown(uid int, gid int) error {
	return chown(d, uid, gid)
}

// Remover -----------------------------------------------------------------------

var _ Remover[Dir] = Dir(".")

// See [os.Remove].
//
// Remove implements [Remover].
func (d Dir) Remove() error {
	return remove(d)
}

// See [os.Rename].
//
// Rename implements [Remover].
func (d Dir) Rename(newPath PathStr) (Dir, error) {
	return rename(d, newPath)
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer[Dir] = Dir(".")

// See [os.RemoveAll].
//
// RemoveAll implements [Destroyer].
func (d Dir) RemoveAll() (Dir, error) {
	return removeAll(d)
}
