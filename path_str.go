package pathlib

import (
	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type PathStr string

// Beholder --------------------------------------------------------------------
var _ Beholder[PathStr] = PathStr(".")
var _ InfallibleBeholder[PathStr] = PathStr(".")

// Note: go's `os.Stat/Lstat` imitates `stat(2)` from POSIX's libc spec.

// Stat implements Beholder.
func (p PathStr) Stat() (OnDisk[PathStr], error) {
	info, err := os.Stat(string(p))
	return onDisk[PathStr]{info, time.Now()}, err
}

// Lstat implements Beholder.
func (p PathStr) Lstat() (OnDisk[PathStr], error) {
	return p.OnDisk()
}

func (p PathStr) OnDisk() (actual OnDisk[PathStr], err error) {
	var info os.FileInfo
	info, err = os.Lstat(string(p))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	actual = onDisk[PathStr]{info, time.Now()}
	return
}

// Note: single-field structs have the same size as their field

func (p PathStr) Exists() (exists bool) {
	_, err := p.OnDisk()
	return !errors.Is(err, fs.ErrNotExist)
}

// MustBeOnDisk implements InfallibleBeholder.
func (p PathStr) MustBeOnDisk() OnDisk[PathStr] {
	return expect(p.OnDisk())
}

// MustLstat implements InfallibleBeholder.
func (p PathStr) MustLstat() OnDisk[PathStr] {
	return expect(p.Lstat())
}

// MustStat implements InfallibleBeholder.
func (p PathStr) MustStat() OnDisk[PathStr] {
	return expect(p.Stat())
}

// PurePath --------------------------------------------------------------------
var _ PurePath = PathStr(".")

// A wrapper around [path/filepath.Join]:
//
// > Join joins any number of path elements into a single path, separating them with an OS
// specific [path/filepath.Separator]. Empty elements are ignored. The result passed
// through [path/filepath.Clean]. However, if the argument list is empty or all its
// elements are empty, Join returns an empty string. On Windows, the result will only be
// a UNC path if the first non-empty element is a UNC path.
//
// Note that this method inherits [path/filepath.Join]'s behavior of ignoring leading
// path separators.
func (p PathStr) Join(segments ...string) PathStr {
	return PathStr(filepath.Join(append([]string{string(p)}, segments...)...))
}
func (p PathStr) Parts() (parts []string) {
	if p == "" {
		return
	}
	var dir, file string
	for {
		dir, file = filepath.Split(dir)
		parts = append(parts, file)
		if dir == "" {
			break
		}
	}
	slices.Reverse(parts)
	return
}

// a wrapper around [path/filepath.Dir]:
//
// > returns all but the last element of path [...]  If the path is empty, Dir returns ".".
// If the path consists entirely of separators, [path/filepath.Dir] returns a single separator. The
// returned path does not end in a separator unless it is the root directory.
func (p PathStr) Parent() Dir {
	return Dir(filepath.Dir(string(p)))
}

// experimental
func (p PathStr) Ancestors() iter.Seq[Dir] {
	return func(yield func(pp Dir) bool) {
		q := p.Parent()
		for yield(q) {
			parent := q.Parent()
			if q == parent {
				break
			}
			q = parent
		}
	}
}

// A wrapper around [path/filepath.Base]:
//
// > Base returns the last element of path. Trailing path separators are removed before
// extracting the last element. If the path is empty, [path/filepath.Base] returns ".".
// If the path consists entirely of separators, [path/filepath.Base] returns a single
// separator.
func (p PathStr) BaseName() string {
	return filepath.Base(string(p))
}

// A wrapper around [path/filepath.Ext]:
//
// > Ext returns the file name extension used by path. The extension is the suffix
// beginning at the final dot in the final element of path; it is empty if there is no
// dot.
func (p PathStr) Ext() string {
	return filepath.Ext(string(p))
}

// Returns true if the path is absolute, false otherwise.
// See [path/filepath.IsAbs] for more details.
func (p PathStr) IsAbsolute() bool {
	return filepath.IsAbs(string(p))
}

// returns true if the path is local/relative, false otherwise.
// see [path/filepath.IsLocal] for more details.
func (p PathStr) IsLocal() bool {
	return filepath.IsLocal(string(p))
}

// Readable --------------------------------------------------------------------
var _ Readable[any] = PathStr(".")
var _ InfallibleReader[any] = PathStr(".")

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

func (p PathStr) MustRead() any {
	return expect(p.Read())
}

func (p PathStr) Open() (*os.File, error) {
	return os.Open(string(p))
}

// Transformer -----------------------------------------------------------------
var _ Transformer[PathStr] = PathStr(".")
var _ InfallibleTransformer[PathStr] = PathStr(".")

// Clean implements Transformer
func (p PathStr) Clean() PathStr {
	return PathStr(filepath.Clean(string(p)))
}

// Abs implements Transformer.
// See [path/filepath.Abs] for more details:
//
// > Abs returns an absolute representation of path. If the path is not absolute
// it will be joined with the current working directory to turn it into an
// absolute path. The absolute path name for a given file is not guaranteed to
// be unique. Abs calls [path/filepath.Clean] on the result.
func (p PathStr) Abs() (PathStr, error) {
	q, err := filepath.Abs(string(p))
	return PathStr(q), err
}

// Localize implements Transformer.
func (p PathStr) Localize() (PathStr, error) {
	q, err := filepath.Localize(string(p))
	return PathStr(q), err
}

// Rel implements Transformer. See [path/filepath.Rel]:
//
// > Rel returns a relative path that is lexically equivalent to [target] when
// joined to basepath with an intervening separator. That is,
// [path/filepath.Join](basepath, Rel(basepath, target)) is equivalent to target
// itself. On success, the returned path will always be relative to basepath,
// even if basepath and target share no elements. An error is returned if target
// can't be made relative to basepath or if knowing the current working directory
// would be necessary to compute it. Rel calls [path/filepath.Clean] on the result.
func (p PathStr) Rel(target Dir) (PathStr, error) {
	result, err := filepath.Rel(string(p), string(target))
	if err != nil {
		return "", errors.Join(err, errors.New("unable to make relative path"))
	}
	return PathStr(result), nil
}

// Expand a leading "~" into the user's home directory. If the home directory cannot be
// determined, the path is returned unchanged.
func (p PathStr) ExpandUser() (result PathStr, err error) {
	if len(p) == 0 || p[0] != '~' || (len(p) > 1 && !os.IsPathSeparator(p[1])) {
		result = p
		return
	}

	var home Dir
	home, err = UserHomeDir()
	if err != nil {
		return
	}

	result = PathStr(PathStr(home) + p[1:])
	return
}

// MustExpandUser implements InfallibleTransformer.
func (p PathStr) MustExpandUser() PathStr {
	return expect(p.ExpandUser())
}

// MustLocalize implements InfallibleTransformer.
func (p PathStr) MustLocalize() PathStr {
	return expect(p.Localize())
}

// MustMakeAbs implements InfallibleTransformer.
func (p PathStr) MustMakeAbs() PathStr {
	return expect(p.Abs())
}

// MustMakeRel implements InfallibleTransformer.
func (p PathStr) MustMakeRel(target Dir) PathStr {
	return expect(p.Rel(target))
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[PathStr] = PathStr(".")
var _ InfallibleManipulator[PathStr] = PathStr(".")

// Chmod implements Manipulator.
func (root PathStr) Chmod(mode os.FileMode) (result PathStr, err error) {
	if err = os.Chmod(string(root), mode); err != nil {
		return
	}
	result = root
	return
}

// Chown implements Manipulator.
func (root PathStr) Chown(uid int, gid int) (result PathStr, err error) {
	if err = os.Chown(string(root), uid, gid); err != nil {
		return
	}
	result = root
	return
}

// Remove implements Manipulator.
func (root PathStr) Remove() error {
	return os.Remove(string(root))
}

// Rename implements Manipulator.
func (root PathStr) Rename(newPath PathStr) (result PathStr, err error) {
	err = os.Rename(string(root), string(newPath))
	if err != nil {
		return
	}
	result = PathStr(newPath)
	return
}

// MustChmod implements InfallibleManipulator.
func (root PathStr) MustChmod(mode os.FileMode) PathStr {
	return expect(root.Chmod(mode))
}

// MustChown implements InfallibleManipulator.
func (root PathStr) MustChown(uid int, gid int) PathStr {
	return expect(root.Chown(uid, gid))
}

// MustRemove implements InfallibleManipulator.
func (root PathStr) MustRemove() {
	if err := root.Remove(); err != nil {
		panic(err)
	}
}

// MustRename implements InfallibleManipulator.
func (root PathStr) MustRename(newPath PathStr) PathStr {
	return expect(root.Rename(newPath))
}

// -----------------------------------------------------------------------------

func (p PathStr) Eq(q PathStr) bool {
	// try to avoid panicking if Cwd() can't be obtained
	if p.IsLocal() && q.IsLocal() {
		return p == q
	}
	// TODO: check that this still works with UNC strings on windows
	return p.MustMakeAbs().MustLocalize() == q.MustMakeAbs().MustLocalize()
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = PathStr(".")
var _ InfallibleDestroyer = PathStr(".")

// RemoveAll implements Destroyer.
func (p PathStr) RemoveAll() error {
	return os.RemoveAll(string(p))
}

// MustRemoveAll implements InfallibleDestroyer.
func (p PathStr) MustRemoveAll() {
	if err := p.RemoveAll(); err != nil {
		panic(err)
	}
}

// casts -----------------------------------------------------------------------

func (p PathStr) AsDir() Dir {
	return Dir(p)
}
func (p PathStr) AsFile() File {
	return File(p)
}
func (p PathStr) AsSymlink() Symlink {
	return Symlink(p)
}
