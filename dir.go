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
func (d Dir) Glob(pattern string) Result[[]PathStr] {
	matches, err := filepath.Glob(filepath.Join(string(d), pattern))
	if err != nil {
		return Result[[]PathStr]{nil, d.failure(err, "glob")}
	}
	result := make([]PathStr, len(matches))
	for i, m := range matches {
		result[i] = PathStr(m)
	}
	return Result[[]PathStr]{result, nil}
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
func (d Dir) Read() Result[[]fs.DirEntry] {
	result, err := os.ReadDir(string(d))
	if err != nil {
		err = d.failure(err, "read")
	}
	return Result[[]fs.DirEntry]{result, err}
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
func (d Dir) Abs() Result[Dir] {
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
func (d Dir) Localize() Result[Dir] {
	return localize(d)
}

// Rel implements [Transformer].
func (d Dir) Rel(base Dir) Result[Dir] {
	return rel(base, d)
}

// ExpandUser implements [Transformer].
func (d Dir) ExpandUser() Result[Dir] {
	return expandUser(d)
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Dir] = Dir(".")

// OnDisk implements [Beholder]
func (d Dir) OnDisk() (result Result[OnDisk[Dir]]) {
	result = lstat(d)
	if result.IsOk() && !result.Val.IsDir() {
		result.Err = WrongTypeOnDisk[Dir]{onDisk[Dir]{result.Val}}
	}
	return

}

// Exists implements [Beholder].
func (d Dir) Exists() bool {
	return d.Lstat().IsOk()
}

// Lstat implements [Beholder].
func (d Dir) Lstat() Result[OnDisk[Dir]] {
	return lstat(d)
}

// Stat implements [Beholder].
func (d Dir) Stat() (result Result[OnDisk[Dir]]) {
	result = stat(d)
	if result.IsOk() && !result.Val.IsDir() {
		result.Err = WrongTypeOnDisk[Dir]{result.Val}
	}
	return
}

// Maker -----------------------------------------------------------------------
var _ Maker[Dir] = Dir("/example")

// Make implements [Maker].
func (d Dir) Make(perm fs.FileMode) (result Result[Dir]) {
	result = Result[Dir]{d, os.Mkdir(string(d), perm)}
	return
}

// MakeAll implements [Maker]
func (d Dir) MakeAll(perm, parentPerm fs.FileMode) (result Result[Dir]) {
	result = Result[Dir]{Val: d}
	if d.Exists() {
		return
	}
	result.Err = d.Parent().MakeAll(parentPerm, parentPerm).Err
	if !result.IsOk() {
		return
	}
	result.Err = os.MkdirAll(string(d), perm)
	return
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[Dir] = Dir(".")

// See [os.Chmod].
// Chmod implements [Manipulator].
func (d Dir) Chmod(mode os.FileMode) Result[Dir] {
	return chmod(d, mode)
}

// See [os.Chown].
//
// Chown implements [Manipulator].
func (d Dir) Chown(uid int, gid int) Result[Dir] {
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
func (d Dir) Rename(newPath PathStr) Result[Dir] {
	return rename(d, newPath)
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer[Dir] = Dir(".")

// See [os.RemoveAll].
//
// RemoveAll implements [Destroyer].
func (d Dir) RemoveAll() Result[Dir] {
	return removeAll(d)
}
