package pathlib

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type Dir PathStr

// A wrapper around [path/filepath.WalkDir], which has the following properties:
//
// > The files are walked in lexical order, which makes the output deterministic but [reading the] entire directory into memory before proceeding to walk that directory.
//
// > WalkDir does not follow symbolic links.
//
// > WalkDir calls [callback] with paths that use the separator character appropriate for the operating system.
func (root Dir) Walk(
	callback func(path string, d fs.DirEntry, err error) error,
) error {
	return filepath.WalkDir(string(root), callback)
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

func (d Dir) RemoveAll() error {
	return os.RemoveAll(string(d))
}

func (d Dir) MustGlob(pattern string) []PathStr {
	return expect(d.Glob(pattern))
}

func (d Dir) Chdir() error {
	return os.Chdir(string(d))
}

// Readable --------------------------------------------------------------------
var _ Readable[[]os.DirEntry] = Dir(".")

// a wrapper around [os.ReadDir]. Read() returns all the entries of the directory sorted
// by filename. If an error occurred reading the directory, Read returns the entries it was
// able to read before the error, along with the error.
func (d Dir) Read() ([]os.DirEntry, error) {
	return os.ReadDir(string(d))
}

// PurePath --------------------------------------------------------------------
var _ PurePath = Dir(".")

// BaseName implements PurePath.
func (d Dir) BaseName() string {
	return PathStr(d).BaseName()
}

// IsAbsolute implements PurePath.
func (d Dir) IsAbsolute() bool {
	return PathStr(d).IsAbsolute()
}

// IsLocal implements PurePath.
func (d Dir) IsLocal() bool {
	return PathStr(d).IsLocal()
}

// Join implements PurePath.
func (d Dir) Join(parts ...string) PathStr {
	return PathStr(d).Join(parts...)
}

func (d Dir) Parts() []string {
	return PathStr(d).Parts()
}

// Parent implements PurePath.
func (d Dir) Parent() Dir {
	return PathStr(d).Parent()
}

// Ext implements PurePath.
func (d Dir) Ext() string {
	return PathStr(d).Ext()
}

// Transformer -----------------------------------------------------------------
var _ Transformer[Dir] = Dir(".")

// Abs implements Transformer.
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

// Localize implements Transformer.
func (d Dir) Localize() (Dir, error) {
	q, err := PathStr(d).Localize()
	return Dir(q), err // this will panic if d is not absolute
}

// Rel implements Transformer.
func (d Dir) Rel(target Dir) (Dir, error) {
	relative, err := PathStr(d).Rel(target)
	return Dir(relative), err
}

// ExpandUser implements TildeTransformer.
func (d Dir) ExpandUser() (Dir, error) {
	p, err := PathStr(d).ExpandUser()
	return Dir(p), err
}

// Beholder --------------------------------------------------------------------
var _ Beholder[Dir] = Dir(".")

func (d Dir) OnDisk() (OnDisk[Dir], error) {
	actual, err := PathStr(d).OnDisk()
	if err != nil {
		return nil, err
	}
	if !actual.IsDir() {
		return nil, WrongTypeOnDisk[Dir]{actual}
	}
	return onDisk[Dir]{actual, time.Now()}, nil
}

// Exists implements Beholder.
func (root Dir) Exists() bool {
	return PathStr(root).Exists()
}

// Lstat implements Beholder.
func (root Dir) Lstat() (OnDisk[Dir], error) {
	return root.OnDisk()
}

// Stat implements Beholder.
func (root Dir) Stat() (result OnDisk[Dir], err error) {
	var info fs.FileInfo
	info, err = os.Stat(string(root))
	if err != nil {
		return
	}
	if !info.IsDir() {
		err = WrongTypeOnDisk[Dir]{info}
		return
	}
	result = onDisk[Dir]{info, time.Now()}
	return
}

// Maker -----------------------------------------------------------------------
var _ Maker[Dir] = Dir("/example")

// Make implements Maker.
func (root Dir) Make(perm ...fs.FileMode) (result Dir, err error) {
	result = root
	const defaultMode fs.FileMode = 0775
main:
	{
		switch len(perm) {
		case 0:
			perm = append(perm, defaultMode)
			goto main
		case 1:
			// mkdir
			err = os.Mkdir(string(root), perm[0])
			return
		default:
			_, err = root.Parent().Make(perm[1:]...)
			if err != nil {
				return
			}
			goto main
		}
	}
}

// MustMake implements Maker.
func (root Dir) MustMake(perm ...fs.FileMode) Dir {
	return expect(root.Make(perm...))
}
