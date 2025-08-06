package pathlib

import (
	"io/fs"
	"os"
	"path/filepath"
)

// A string that represents a directory. The directory may or may not exist on-disk,
// and the string may or may not end in an [os.PathSeparator].
type Dir PathStr

func (d Dir) failure(err error, op string) error {
	return Error[Dir]{
		Path: d,
		Op:   op,
		Err:  err,
	}
}

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
	if err != nil {
		return nil, d.failure(err, "glob")
	}
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
type DirEntry[P Kind] struct {
	fs.DirEntry
}

var _ Readable[[]fs.DirEntry] = Dir(".")

// a wrapper around [os.ReadDir]:
//
// > [os.ReadDir] returns all the entries of the directory sorted
// by filename. If an error occurred reading the directory, ReadDir returns the entries it was
// able to read before the error, along with the error.
func (d Dir) Read() (result []fs.DirEntry, err error) {
	result, err = os.ReadDir(string(d))
	if err != nil {
		err = d.failure(err, "read")
	}
	return
}

// See [os.ReadDir].
func (d Dir) MustRead() []fs.DirEntry {
	return expect(d.Read())
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
var _ InfallibleTransformer[Dir] = Dir(".")

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

// Localize implements [Transformer].
func (d Dir) Localize() (Dir, error) {
	return localize(d)
}

// Rel implements [Transformer].
func (d Dir) Rel(target Dir) (Dir, error) {
	return rel(d, target)
}

// ExpandUser implements [Transformer].
func (d Dir) ExpandUser() (Dir, error) {
	return expandUser(d)
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

// OnDisk implements [Beholder]
func (d Dir) OnDisk() (OnDisk[Dir], error) {
	actual, err := PathStr(d).OnDisk()
	if err != nil {
		return nil, d.failure(err, "observe on-disk")
	}
	if !actual.IsDir() {
		return nil, WrongTypeOnDisk[Dir]{onDisk[Dir]{actual}}
	}
	return onDisk[Dir]{actual}, nil
}

// Exists implements [Beholder].
func (d Dir) Exists() bool {
	return PathStr(d).Exists()
}

// Lstat implements [Beholder].
func (d Dir) Lstat() (OnDisk[Dir], error) {
	return lstat(d)
}

// Stat implements [Beholder].
func (d Dir) Stat() (result OnDisk[Dir], err error) {
	result, err = stat(d)
	if err != nil {
		return
	}
	if !result.IsDir() {
		err = WrongTypeOnDisk[Dir]{result}
		return
	}
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
	result = d
	err = os.Mkdir(string(d), perm)
	if err != nil {
		err = d.failure(err, "make")
	}
	return
}

// MakeAll implements [Maker]
func (d Dir) MakeAll(perm, parentPerm fs.FileMode) (result Dir, err error) {
	result = d
	if d.Exists() {
		return
	}
	_, err = d.Parent().MakeAll(parentPerm, parentPerm)
	if err != nil {
		err = d.failure(err, "make parents of")
		return
	}
	err = os.MkdirAll(string(d), perm)
	if err != nil {
		err = d.failure(err, "make all of")
	}
	return
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
	return chmod(d, mode)
}

// See [os.Chown].
//
// Chown implements [Manipulator].
func (d Dir) Chown(uid int, gid int) (Dir, error) {
	return chown(d, uid, gid)
}

// See [os.Remove].
//
// Remove implements [Manipulator].
func (d Dir) Remove() error {
	return os.Remove(string(d))
}

// See [os.Rename].
//
// Rename implements [Manipulator].
func (d Dir) Rename(newPath PathStr) (Dir, error) {
	return rename(d, newPath)
}

// Panics if [Dir.Chmod] returns an error.
// See [os.Chmod].
//
// MustChmod implements [InfallibleManipulator].
func (d Dir) MustChmod(mode os.FileMode) Dir {
	return expect(d.Chmod(mode))
}

// Panics if [Dir.Chown] returns an error.
// See [os.Chown].
//
// MustChown implements [InfallibleManipulator].
func (d Dir) MustChown(uid int, gid int) Dir {
	return expect(d.Chown(uid, gid))
}

// Panics if [Dir.Remove] returns an error.
// See [os.Remove].
//
// MustRemove implements [InfallibleManipulator].
func (d Dir) MustRemove() {
	expect[any](nil, d.Remove())
}

// Panics if [Dir.Rename] returns an error. See also: [os.Rename].
//
// MustRename implements [InfallibleManipulator].
func (d Dir) MustRename(newPath PathStr) Dir {
	return expect(d.Rename(newPath))
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = Dir(".")
var _ InfallibleDestroyer = Dir(".")

// See [os.RemoveAll].
//
// RemoveAll implements [Destroyer].
func (d Dir) RemoveAll() error {
	return os.RemoveAll(string(d))
}

// Panics if [PathStr.RemoveAll] returns an error. See also: [os.RemoveAll].
//
// MustRemoveAll implements [InfallibleDestroyer].
func (d Dir) MustRemoveAll() {
	PathStr(d).MustRemoveAll()
}
