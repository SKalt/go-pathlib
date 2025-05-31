package pathlib

import (
	"os"
	"path/filepath"
)

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
	if p.IsDir() {
		return Dir(p)
	} else {
		return p.Parent()
	}
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
