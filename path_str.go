package pathlib

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type PathStr string

// Note: go's `os.Stat/Lstat` imitates `stat(2)` from POSIX's libc spec.

// Stat implements Beholder.
func (p PathStr) Stat() (OnDisk[PathStr], error) {
	info, err := os.Stat(string(p))
	return onDisk[PathStr]{info}, err
}

// Lstat implements Beholder.
func (p PathStr) Lstat() (OnDisk[PathStr], error) {
	return p.OnDisk()
}

var (
	_ PurePath             = PathStr(".")
	_ Transformer[PathStr] = PathStr(".")
	_ Readable[any]        = PathStr(".")
	_ Beholder[PathStr]    = PathStr(".")
)

func (p PathStr) OnDisk() (actual OnDisk[PathStr], err error) {
	var info os.FileInfo
	info, err = os.Lstat(string(p))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	actual = onDisk[PathStr]{info}
	return
}

// Note: single-field structs have the same size as their field

func (p PathStr) Exists() (exists bool) {
	_, err := p.OnDisk()
	return !errors.Is(err, fs.ErrNotExist)
}

// Returns true if the path is absolute, false otherwise.
// See [filepath.IsAbs] for more details.
func (p PathStr) IsAbsolute() bool {
	return filepath.IsAbs(string(p))
}

// returns true if the path is local/relative, false otherwise.
// see [filepath.IsLocal] for more details.
func (p PathStr) IsLocal() bool {
	return filepath.IsLocal(string(p))
}

func (p PathStr) Read() (result any, err error) {
	// can't define this switch as a method of OnDisk[P] since OnDisk[P] has to handle
	// any kind of path
	var actual OnDisk[PathStr]
	actual, err = p.OnDisk()
	if err != nil {
		return
	}
	mode := actual.Mode()

	if mode.IsRegular() {
		result, err = os.ReadFile(string(p))
	} else if mode.IsDir() {
		result, err = Dir(p).Read()
	} else if isSymLink(mode) {
		result, err = Symlink(p).Read()
	} else if isCharDevice(mode) {
		// TODO
	} else if isDevice(mode) {
		// TODO
	} else if isFifo(mode) {
		// TODO
	} else if isSocket(mode) {
		// TODO
	} else if isTemporary(mode) {
		// TODO
	}
	return
}

func (p PathStr) Open() (*os.File, error) {
	return os.Open(string(p))
}

func (p PathStr) WithOpen(cb func(*os.File) error) error { // FIXME: name
	f, err := p.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	return cb(f)
}

func (p PathStr) Clean() PathStr {
	return PathStr(filepath.Clean(string(p)))
}

// Abs implements PurePath.
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

// Localize implements PurePath.
func (p PathStr) Localize() (PathStr, error) {
	q, err := filepath.Localize(string(p))
	return PathStr(q), err
}

// Rel implements PurePath. See [path/filepath.Rel]:
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

// A wrapper around [path/filepath.Join]:
//
// > Join joins any number of path elements into a single path, separating them with an OS
// specific [path/filepath.Separator]. Empty elements are ignored. The result passed
// through [path/filepath.Clean]. However, if the argument list is empty or all its
// elements are empty, Join returns an empty string. On Windows, the result will only be
// a UNC path if the first non-empty element is a UNC path.
func (p PathStr) Join(segments ...string) PathStr {
// FIXME: handle joining absolute paths
	return PathStr(filepath.Join(append([]string{string(p)}, segments...)...))
}

// a wrapper around [path/filepath.Dir]:
//
// > returns all but the last element of path [...]  If the path is empty, Dir returns ".".
// If the path consists entirely of separators, [path/filepath.Dir] returns a single separator. The
// returned path does not end in a separator unless it is the root directory.
func (p PathStr) Parent() Dir {
	return Dir(filepath.Dir(string(p)))
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

var homeDir string

// caches the user's home directory. Returns an empty string if it cannot be determined.
func getHomeDir() (home string, err error) {
	if homeDir != "" {
		home = homeDir
	} else {
		home, err = os.UserHomeDir()
		if err == nil {
			homeDir = home
		}
	}
	return
}

// Expand a leading "~" into the user's home directory. If the home directory cannot be
// determined, the path is returned unchanged.
func (p PathStr) ExpandUser() (PathStr, error) {
	if len(p) > 0 && p[0] == '~' {
		if home, err := getHomeDir(); home != "" && err == nil {
			return PathStr(PathStr(home) + p[1:]), nil // FIXME: check p[2] == "/"
		}
	}
	return p, nil

}
