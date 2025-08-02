package pathlib

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Dir PathStr

// A wrapper around [path/filepath.WalkDir], which has the following properties:
//
// > The files are walked in lexical order, which makes the output deterministic but [reading the] entire directory into memory before proceeding to walk that directory.
//
// > WalkDir does not follow symbolic links.
//
// > WalkDir calls [callback] with paths that use the separator character appropriate for the operating system.
func (d Dir) Walk(
	callback func(path string, d fs.DirEntry, err error) error,
) error {
	return filepath.WalkDir(string(d), callback)
}

// See [path/filepath.Glob]:
//
// > Glob returns the names of all files matching pattern or nil if there is no
// matching file. The syntax of patterns is the same as in [path/filepath.Match].
// The pattern may describe hierarchical names such as /usr/*/bin/ed (assuming
// the [path/filepath.Separator] is '/').
//
// > Glob ignores file system errors such as I/O errors reading directories. The only possible returned error is [path/filepath.ErrBadPattern], when pattern is malformed.
func (d Dir) Glob(pattern string) ([]PathStr, error) {
	matches, err := filepath.Glob(filepath.Join(string(d), pattern))
	result := make([]PathStr, len(matches))
	for i, m := range matches {
		result[i] = PathStr(m)
	}
	return result, err
}

// Panics if [path/filepath.Glob] returns an error.
func (d Dir) MustGlob(pattern string) []PathStr {
	return expect(d.Glob(pattern))
}

// CHange DIRectory. See [os.Chdir].
func (d Dir) Chdir() error {
	return os.Chdir(string(d))
}

// Readable --------------------------------------------------------------------
var _ Readable[[]os.DirEntry] = Dir(".")

// a wrapper around [os.ReadDir]:
//
// > [os.ReadDir] returns all the entries of the directory sorted
// by filename. If an error occurred reading the directory, ReadDir returns the entries it was
// able to read before the error, along with the error.
func (d Dir) Read() ([]os.DirEntry, error) {
	return os.ReadDir(string(d))
}

// See [os.ReadDir].
func (d Dir) MustRead() []os.DirEntry {
	return expect(d.Read())
}

// PurePath --------------------------------------------------------------------
var _ PurePath = Dir(".")

// See [filepath.Base].
// BaseName implements [PurePath].
func (d Dir) BaseName() string {
	return PathStr(d).BaseName()
}

// Returns true if the path is absolute, false otherwise.
// See [filepath.IsAbs] for more details.
//
// IsAbsolute implements [PurePath].
func (d Dir) IsAbsolute() bool {
	return PathStr(d).IsAbsolute()
}

// IsLocal implements [PurePath].
func (d Dir) IsLocal() bool {
	return PathStr(d).IsLocal()
}

// Join implements [PurePath].
func (d Dir) Join(parts ...string) PathStr {
	return PathStr(d).Join(parts...)
}

func (d Dir) Parts() []string {
	return PathStr(d).Parts()
}

// Parent implements [PurePath].
func (d Dir) Parent() Dir {
	return PathStr(d).Parent()
}

// Ext implements [PurePath].
func (d Dir) Ext() string {
	return PathStr(d).Ext()
}

// Transformer -----------------------------------------------------------------
var _ Transformer[Dir] = Dir(".")
var _ InfallibleTransformer[Dir] = Dir(".")

// Abs implements [Transformer].
func (d Dir) Abs() (Dir, error) {
	abs, err := PathStr(d).Abs()
	return Dir(abs), err
}
func (d Dir) Eq(other Dir) (equivalent bool) {
	return PathStr(d).Eq(PathStr(other))
}
func (d Dir) Clean() Dir {
	return Dir(filepath.Clean(string(d)))
}

// Localize implements [Transformer].
func (d Dir) Localize() (Dir, error) {
	q, err := PathStr(d).Localize()
	return Dir(q), err // this will panic if d is not absolute
}

// Rel implements [Transformer].
func (d Dir) Rel(target Dir) (Dir, error) {
	relative, err := PathStr(d).Rel(target)
	return Dir(relative), err
}

// ExpandUser implements [Transformer].
func (d Dir) ExpandUser() (Dir, error) {
	p, err := PathStr(d).ExpandUser()
	return Dir(p), err
}

// MustExpandUser implements [InfallibleTransformer].
func (d Dir) MustExpandUser() Dir {
	return expect(d.ExpandUser())
}

// MustLocalize implements [InfallibleTransformer].
func (d Dir) MustLocalize() Dir {
	return expect(d.Localize())
}

// MustMakeAbs implements [InfallibleTransformer].
func (d Dir) MustMakeAbs() Dir {
	return expect(d.Abs())
}

// MustMakeRel implements [InfallibleTransformer].
func (d Dir) MustMakeRel(target Dir) Dir {
	return expect(d.Rel(target))
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Dir] = Dir(".")
var _ InfallibleBeholder[Dir] = Dir(".")

func (d Dir) OnDisk() (OnDisk[Dir], error) {
	actual, err := PathStr(d).OnDisk()
	if err != nil {
		return nil, err
	}
	if !actual.IsDir() {
		return nil, WrongTypeOnDisk[Dir]{actual}
	}
	return onDisk[Dir]{actual, }, nil
}

// Exists implements [Beholder].
func (d Dir) Exists() bool {
	return PathStr(d).Exists()
}

// Lstat implements [Beholder].
func (d Dir) Lstat() (OnDisk[Dir], error) {
	return d.OnDisk()
}

// Stat implements [Beholder].
func (d Dir) Stat() (result OnDisk[Dir], err error) {
	var info fs.FileInfo
	info, err = os.Stat(string(d))
	if err != nil {
		return
	}
	if !info.IsDir() {
		err = WrongTypeOnDisk[Dir]{info}
		return
	}
	result = onDisk[Dir]{info, }
	return
}

// MustLstat implements [InfallibleBeholder].
func (d Dir) MustLstat() OnDisk[Dir] {
	return expect(d.Lstat())
}

// MustOnDisk implements [InfallibleBeholder].
func (d Dir) MustBeOnDisk() OnDisk[Dir] {
	return expect(d.OnDisk())
}

// MustStat implements [InfallibleBeholder].
func (d Dir) MustStat() OnDisk[Dir] {
	return expect(d.Stat())
}

// Maker -----------------------------------------------------------------------
var _ InfallibleMaker[Dir] = Dir("/example")
var _ Maker[Dir] = Dir("/example")

// Make implements [Maker].
func (d Dir) Make(perm fs.FileMode) (result Dir, err error) {
	err = os.Mkdir(string(d), perm)
	if err != nil {
		result = d
	}
	return
}

// MakeAll implements [Maker]
func (d Dir) MakeAll(perm, parentPerm fs.FileMode) (result Dir, err error) {
	if d.Exists() {
		result = d
		return
	}
	_, err = d.Parent().MakeAll(parentPerm, parentPerm)
	if err != nil {
		return
	}
	return d.Make(perm)
}

// MustMake implements [InfallibleMaker].
func (root Dir) MustMake(perm fs.FileMode) Dir {
	return expect(root.Make(perm))
}

// Panics if [Dir.MakeAll].returns an error.
//
// MustMakeAll implements [InfallibleMaker]
func (d Dir) MustMakeAll(perm, parentPerm fs.FileMode) Dir {
	return expect(d.MakeAll(perm, parentPerm))
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[Dir] = Dir(".")
var _ InfallibleManipulator[Dir] = Dir(".")

// See [os.Chmod].
// Chmod implements [Manipulator].
func (d Dir) Chmod(mode os.FileMode) (Dir, error) {
	result, err := PathStr(d).Chmod(mode)
	return Dir(result), err
}

// See [os.Chown].
// Chown implements [Manipulator].
func (d Dir) Chown(uid int, gid int) (Dir, error) {
	result, err := PathStr(d).Chown(uid, gid)
	return Dir(result), err

}

// See [os.Remove].
// Remove implements [Manipulator].
func (d Dir) Remove() error {
	return os.Remove(string(d))
}

// See [os.Rename].
// Rename implements [Manipulator].
func (d Dir) Rename(newPath PathStr) (Dir, error) {
	result, err := PathStr(d).Rename(newPath)
	return Dir(result), err
}

// See [os.Chmod].
// MustChmod implements [InfallibleManipulator].
func (d Dir) MustChmod(mode os.FileMode) Dir {
	return expect(d.Chmod(mode))
}

// See [os.Chown].
// MustChown implements [InfallibleManipulator].
func (d Dir) MustChown(uid int, gid int) Dir {
	return expect(d.Chown(uid, gid))
}

// See [os.Remove].
// MustRemove implements [InfallibleManipulator].
func (d Dir) MustRemove() {
	expect[any](nil, d.Remove())
}

// See [os.Rename].
// MustRename implements [InfallibleManipulator].
func (d Dir) MustRename(newPath PathStr) Dir {
	return expect(d.Rename(newPath))
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = Dir(".")
var _ InfallibleDestroyer = Dir(".")

// See [os.RemoveAll].
// RemoveAll implements [Destroyer].
func (d Dir) RemoveAll() error {
	return os.RemoveAll(string(d))
}

// See [os.RemoveAll].
// MustRemoveAll implements [InfallibleDestroyer].
func (d Dir) MustRemoveAll() {
	PathStr(d).MustRemoveAll()
}
