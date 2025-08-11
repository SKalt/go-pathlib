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

// Note: go's `os.Stat/Lstat` imitates `stat(2)` from POSIX's libc spec.

// See [os.Stat].
//
// Stat implements [Beholder].
func (p PathStr) Stat() Result[OnDisk[PathStr]] {
	return stat(p)
}

// See [os.Lstat].
//
// Lstat implements [Beholder].
func (p PathStr) Lstat() Result[OnDisk[PathStr]] {
	return lstat(p)
}

// OnDisk implements [Beholder].
func (p PathStr) OnDisk() Result[OnDisk[PathStr]] {
	return lstat(p)
}

// Exists implements [Beholder].
func (p PathStr) Exists() (exists bool) {
	return !errors.Is(p.OnDisk().err, fs.ErrNotExist)
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

// Read attempts to read what the path represents. See [File.Read], [Dir.Read], and
// [Symlink.Read] for the possibilities.
//
// Read implements [Readable].
func (p PathStr) Read() Result[any] {
	// can't define this switch as a method of OnDisk[P] since OnDisk[P] has to handle
	// any kind of path
	var actual os.FileInfo
	var val any
	var err error
	actual, err = p.OnDisk().Unpack()
	if err != nil {
		return Result[any]{nil, err}
	}

	if actual.Mode().IsDir() {
		val, err = Dir(p).Read().Unpack()
	} else if actual.Mode()&fs.ModeSymlink == fs.ModeSymlink {
		val, err = Symlink(p).Read().Unpack()
	} else {
		val, err = File(p).Read().Unpack()
	}
	return Result[any]{val, err}
}

// See [os.Open].
func (p PathStr) Open() Result[*os.File] {
	handle, err := os.Open(string(p))
	return Result[*os.File]{handle, err}
}

// Transformer -----------------------------------------------------------------
var _ Transformer[PathStr] = PathStr(".")

func (p PathStr) String() string {
	return string(p)
}

// See [path/filepath.Clean].
//
// Clean implements [Transformer].
func (p PathStr) Clean() PathStr {
	return clean(p)
}

// Abs implements [Transformer].
// See [path/filepath.Abs] for more details.
func (p PathStr) Abs() Result[PathStr] {
	return abs(p)
}

// See [path/filepath.Localize].
// Localize implements [Transformer].
func (p PathStr) Localize() Result[PathStr] {
	return localize(p)
}

// Rel implements [Transformer]. See [path/filepath.Rel]:
//
// See [path/filepath.Rel].
func (p PathStr) Rel(base Dir) Result[PathStr] {
	return rel(base, p)
}

// Expand a leading "~" into the user's home directory. If the home directory cannot be
// determined, the path is returned unchanged.
func (p PathStr) ExpandUser() Result[PathStr] {
	return expandUser(p)
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[PathStr] = PathStr(".")

// Chmod implements [Manipulator].
func (p PathStr) Chmod(mode os.FileMode) Result[PathStr] {
	return chmod(p, mode)
}

// Change Ownership of the path.
//
// Chown implements [Manipulator].
func (p PathStr) Chown(uid int, gid int) Result[PathStr] {
	return chown(p, uid, gid)
}

// Remove implements [Manipulator].
func (p PathStr) Remove() Result[PathStr] {
	return remove(p)
}

// Rename implements [Manipulator].
func (p PathStr) Rename(newPath PathStr) Result[PathStr] {
	return rename(p, newPath)
}

// -----------------------------------------------------------------------------

func (p PathStr) Eq(q PathStr) bool {
	// try to avoid panicking if Cwd() can't be obtained
	p = p.Clean()
	q = q.Clean()
	if p.IsLocal() && q.IsLocal() {
		return p == q
	}
	// TODO: check that this still works with UNC strings on windows
	return p.Abs().Unwrap() == q.Abs().Unwrap()
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer[PathStr] = PathStr(".")

// See [os.RemoveAll].
//
// RemoveAll implements [Destroyer].
func (p PathStr) RemoveAll() Result[PathStr] {
	return Result[PathStr]{p, os.RemoveAll(string(p))}
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
