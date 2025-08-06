package pathlib

import (
	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

type PathStr string

// Beholder --------------------------------------------------------------------
var _ Beholder[PathStr] = PathStr(".")
var _ InfallibleBeholder[PathStr] = PathStr(".")

// Note: go's `os.Stat/Lstat` imitates `stat(2)` from POSIX's libc spec.

// See [os.Stat].
//
// Stat implements [Beholder].
func (p PathStr) Stat() (OnDisk[PathStr], error) {
	return stat(p)
}

// See [os.Lstat].
//
// Lstat implements [Beholder].
func (p PathStr) Lstat() (OnDisk[PathStr], error) {
	return lstat(p)
}

// OnDisk implements [Beholder].
func (p PathStr) OnDisk() (actual OnDisk[PathStr], err error) {
	return lstat(p)
}

// Exists implements [Beholder].
func (p PathStr) Exists() (exists bool) {
	_, err := p.OnDisk()
	return !errors.Is(err, fs.ErrNotExist)
}

// Panics if [PathStr.OnDisk] returns an error.
//
// MustBeOnDisk implements [InfallibleBeholder].
func (p PathStr) MustBeOnDisk() OnDisk[PathStr] {
	return expect(p.OnDisk())
}

// Panics if [PathStr.Lstat] returns an error.
//
// MustLstat implements [InfallibleBeholder].
func (p PathStr) MustLstat() OnDisk[PathStr] {
	return expect(p.Lstat())
}

// Panics if [PathStr.Stat] returns an error.
//
// MustStat implements [InfallibleBeholder].
func (p PathStr) MustStat() OnDisk[PathStr] {
	return expect(p.Stat())
}

// PurePath --------------------------------------------------------------------
var _ PurePath = PathStr(".")

// A wrapper around [path/filepath.Join].
func (p PathStr) Join(segments ...string) PathStr {
	return join(p, segments...)
}

// Splits the path at every character that's a path separator. Omits empty segments.
// See [os.IsPathSeparator].
func (p PathStr) Parts() (parts []string) {
	input := string(p)
	vol := filepath.VolumeName(input)
	if vol != "" {
		parts = append(parts, vol)
		input = input[len(vol):]
	}
	if input == "" {
		return
	}
	var i, last = 1, 0
	if os.IsPathSeparator(input[0]) {
		parts = append(parts, input[:1])
		last = 1
	}
	for i < len(input) {
		if os.IsPathSeparator(input[i]) {
			part := input[last:i]
			if part != "" {
				parts = append(parts, part)
			}
			last = i + 1
		}
		i++
	}
	if last < i {
		parts = append(parts, input[last:])
	}

	return
}

// a wrapper around [path/filepath.Dir].
//
// Parent implements [PurePath].
func (p PathStr) Parent() Dir {
	return parent(p)
}

// experimental
func (p PathStr) Ancestors() iter.Seq[Dir] {
	return ancestors(p)
}

// A wrapper around [path/filepath.Base].
//
// BaseName implements [PurePath].
func (p PathStr) BaseName() string {
	return baseName(p)
}

// A wrapper around [path/filepath.Ext].
//
// Ext implements [PurePath]
func (p PathStr) Ext() string {
	return ext(p)
}

// Returns true if the path is absolute, false otherwise.
// See [path/filepath.IsAbs] for more details.
//
// IsAbsolute implements [PurePath].
func (p PathStr) IsAbsolute() bool {
	return isAbsolute(p)
}

// returns true if the path is local/relative, false otherwise.
// see [path/filepath.IsLocal] for more details.
//
// IsLocal implements [PurePath].
func (p PathStr) IsLocal() bool {
	return isLocal(p)
}

// Readable --------------------------------------------------------------------
var _ Readable[any] = PathStr(".")
var _ InfallibleReader[any] = PathStr(".")

// Read attempts to read what the path represents. See [File.Read], [Dir.Read], and
// [Symlink.Read] for the possibilities.
//
// Read implements [Readable].
func (p PathStr) Read() (result any, err error) {
	// can't define this switch as a method of OnDisk[P] since OnDisk[P] has to handle
	// any kind of path
	var actual OnDisk[PathStr]
	actual, err = p.OnDisk()
	if err != nil {
		return
	}

	if actual.Mode().IsDir() {
		result, err = Dir(p).Read()
	} else if isSymLink(actual.Mode()) {
		result, err = Symlink(p).Read()
	} else {
		result, err = os.ReadFile(string(p))
	}
	return
}

// Panics if [PathStr.Read] returns an error.
//
// MustRead implements [InfallibleReader].
func (p PathStr) MustRead() any {
	return expect(p.Read())
}

// See [os.Open].
func (p PathStr) Open() (*os.File, error) {
	return os.Open(string(p))
}

// Transformer -----------------------------------------------------------------
var _ Transformer[PathStr] = PathStr(".")
var _ InfallibleTransformer[PathStr] = PathStr(".")

// See [path/filepath.Clean].
//
// Clean implements [Transformer].
func (p PathStr) Clean() PathStr {
	return clean(p)
}

// Abs implements [Transformer].
// See [path/filepath.Abs] for more details.
func (p PathStr) Abs() (PathStr, error) {
	return abs(p)
}

// See [path/filepath.Localize].
// Localize implements [Transformer].
func (p PathStr) Localize() (PathStr, error) {
	return localize(p)
}

// Rel implements [Transformer]. See [path/filepath.Rel]:
//
// See [path/filepath.Rel].
func (p PathStr) Rel(target Dir) (PathStr, error) {
	return rel(p, target)
}

// Expand a leading "~" into the user's home directory. If the home directory cannot be
// determined, the path is returned unchanged.
func (p PathStr) ExpandUser() (result PathStr, err error) {
	return expandUser(p)
}

// MustExpandUser implements [InfallibleTransformer].
func (p PathStr) MustExpandUser() PathStr {
	return expect(p.ExpandUser())
}

// MustLocalize implements [InfallibleTransformer].
func (p PathStr) MustLocalize() PathStr {
	return expect(p.Localize())
}

// MustMakeAbs implements [InfallibleTransformer].
func (p PathStr) MustMakeAbs() PathStr {
	return expect(p.Abs())
}

// MustMakeRel implements [InfallibleTransformer].
func (p PathStr) MustMakeRel(target Dir) PathStr {
	return expect(p.Rel(target))
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[PathStr] = PathStr(".")
var _ InfallibleManipulator[PathStr] = PathStr(".")

// Chmod implements [Manipulator].
func (p PathStr) Chmod(mode os.FileMode) (result PathStr, err error) {
	return chmod(p, mode)
}

// Change Ownership of the path.
//
// Chown implements [Manipulator].
func (p PathStr) Chown(uid int, gid int) (result PathStr, err error) {
	return chown(p, uid, gid)
}

// Remove implements [Manipulator].
func (p PathStr) Remove() error {
	return os.Remove(string(p))
}

// Rename implements [Manipulator].
func (p PathStr) Rename(newPath PathStr) (result PathStr, err error) {
	return rename(p, newPath)
}

// MustChmod implements [InfallibleManipulator].
func (p PathStr) MustChmod(mode os.FileMode) PathStr {
	return expect(p.Chmod(mode))
}

// MustChown implements [InfallibleManipulator].
func (p PathStr) MustChown(uid int, gid int) PathStr {
	return expect(p.Chown(uid, gid))
}

// MustRemove implements [InfallibleManipulator].
func (p PathStr) MustRemove() {
	if err := p.Remove(); err != nil {
		panic(err)
	}
}

// MustRename implements [InfallibleManipulator].
func (p PathStr) MustRename(newPath PathStr) PathStr {
	return expect(p.Rename(newPath))
}

// -----------------------------------------------------------------------------

func (p PathStr) Eq(q PathStr) bool {
	// try to avoid panicking if Cwd() can't be obtained
	if p.IsLocal() && q.IsLocal() {
		return p == q
	}
	// TODO: check that this still works with UNC strings on windows
	return p.MustMakeAbs() == q.MustMakeAbs()
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = PathStr(".")
var _ InfallibleDestroyer = PathStr(".")

// See [os.RemoveAll].
//
// RemoveAll implements [Destroyer].
func (p PathStr) RemoveAll() error {
	return os.RemoveAll(string(p))
}

// Panics if [PathStr.RemoveAll] returns an error.
//
// MustRemoveAll implements [InfallibleDestroyer].
func (p PathStr) MustRemoveAll() {
	if err := p.RemoveAll(); err != nil {
		panic(err)
	}
}

// casts -----------------------------------------------------------------------

// Utility function to declare that the PathStr represents a directory
func (p PathStr) AsDir() Dir {
	return Dir(p)
}

// Utility function to declare that the PathStr represents a file
func (p PathStr) AsFile() File {
	return File(p)
}

// Utility function to declare that the PathStr represents a symlink.
func (p PathStr) AsSymlink() Symlink {
	return Symlink(p)
}
