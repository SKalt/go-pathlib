package pathlib

import (
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
func (p PathStr) Stat() (Info[PathStr], error) {
	return stat(p)
}

// See [os.Lstat].
//
// Lstat implements [Beholder].
func (p PathStr) Lstat() (Info[PathStr], error) {
	return lstat(p)
}

// OnDisk implements [Beholder].
func (p PathStr) OnDisk() (Info[PathStr], error) {
	return lstat(p)
}

// Exists implements [Beholder].
func (p PathStr) Exists() bool {
	return exists(p)
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
// [Symlink.Read] for the possible return types.
//
// Read implements [Readable].
func (p PathStr) Read() (any, error) {
	// can't define this switch as a method of OnDisk[P] since OnDisk[P] has to handle
	// any kind of path
	var actual os.FileInfo
	var val any
	var err error
	actual, err = p.OnDisk()
	if err != nil {
		return nil, err
	}

	if actual.Mode().IsDir() {
		val, err = Dir(p).Read()
	} else if actual.Mode()&fs.ModeSymlink == fs.ModeSymlink {
		val, err = Symlink(p).Read()
	} else {
		val, err = File(p).Read()
	}
	return val, err
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
func (p PathStr) Rel(base Dir) (PathStr, error) {
	return rel(base, p)
}

// Expand a leading "~" into the user's home directory. If the home directory cannot be
// determined, the path is returned unchanged.
func (p PathStr) ExpandUser() (PathStr, error) {
	return expandUser(p)
}

// Changer ----------------------------------------------------------------------
var _ Changer = PathStr(".")

// Chmod implements [Changer].
func (p PathStr) Chmod(mode os.FileMode) error {
	return chmod(p, mode)
}

// Change Ownership of the path.
//
// Chown implements [Changer].
func (p PathStr) Chown(uid int, gid int) error {
	return chown(p, uid, gid)
}

// Mover ------------------------------------------------------------------------
var _ Remover[PathStr] = PathStr(".")

// Remove implements [Remover].
func (p PathStr) Remove() error {
	return remove(p)
}

// Rename implements [Remover].
func (p PathStr) Rename(newPath PathStr) (PathStr, error) {
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
	x, err := p.Abs()
	if err != nil {
		return false
	}
	y, err := q.Abs()
	if err != nil {
		return false
	}
	// TODO: check that this still works with UNC strings on windows
	return x == y
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
