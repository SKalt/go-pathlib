package pathlib

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type PathStr string

var (
	_ PurePath             = PathStr(".")
	_ Transformer[PathStr] = PathStr(".")
	_ Readable[any]        = PathStr(".")
)

func (p PathStr) OnDisk() (onDisk *OnDisk[PathStr], err error) {
	var info os.FileInfo
	info, err = lstat(p)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	onDisk = &OnDisk[PathStr]{info}
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
	return isAbsolute(p)
}

// returns true if the path is local/relative, false otherwise.
// see [filepath.IsLocal] for more details.
func (p PathStr) IsLocal() bool {
	return isLocal(p)
}

func (p PathStr) Read() (result any, err error) {
	// can't define this switch as a method of OnDisk[P] since OnDisk[P] has to handle
	// any kind of path
	var onDisk *OnDisk[PathStr]
	onDisk, err = p.OnDisk()
	if err != nil {
		return
	}
	mode := (*onDisk).Mode()

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

// Abs implements PurePath.
func (p PathStr) Abs(cwd Dir) PathStr {
	if p.IsAbsolute() {
		return p
	}
	// assume that cwd is absolute
	if !cwd.IsAbsolute() {
		panic("cwd must be absolute")
	}
	return cwd.Join(string(p)) // join the path with the current working directory
}

// Localize implements PurePath.
func (p PathStr) Localize() PathStr {
	filepath.Localize(string(p))
	panic("unimplemented") // TODO: implement this
}

// Rel implements PurePath.
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
// > specific [path/filepath.Separator]. Empty elements are ignored. The result passed
// > through [path/filepath.Clean]. However, if the argument list is empty or all its
// > elements are empty, Join returns an empty string. On Windows, the result will only be
// > a UNC path if the first non-empty element is a UNC path.
func (p PathStr) Join(segments ...string) PathStr {
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

// Either the parent of the path or the path itself, if it's a directory
func (p PathStr) NearestDir() Dir {
	if onDisk, err := p.OnDisk(); err == nil && onDisk.IsDir() {
		return Dir(p) // p is a directory, return it as a Dir
	}
	return p.Parent()

}

var homeDir string

// caches the user's home directory. Returns an empty string if it cannot be determined.
func getHomeDir() string {
	if homeDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		homeDir = home
	}
	return homeDir
}

// Expand a leading "~" into the user's home directory. If the home directory cannot be
// determined, the path is returned unchanged.
func (p PathStr) ExpandUser() PathStr {
	if len(p) > 0 && p[0] == '~' {
		if home := getHomeDir(); home != "" {
			return PathStr(PathStr(home) + p[1:]) // FIXME: check p[2] == "/"
		}
	}
	return p

}
