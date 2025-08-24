package pathlib

import (
	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

// SEe [os.Stat].
func stat[P Kind](p P) (Info[P], error) {
	info, err := os.Stat(string(p))
	return onDisk[P]{p, info}, err
}

func lstat[P Kind](p P) (Info[P], error) {
	info, err := os.Lstat(string(p))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return onDisk[P]{p, info}, err
}

func exists[P Kind](p P) bool {
	_, err := lstat(p)
	return !errors.Is(err, fs.ErrNotExist)
}

func join[P Kind](p P, segments ...string) PathStr {
	return PathStr(filepath.Join(append([]string{string(p)}, segments...)...))
}

// a wrapper around [path/filepath.Dir].
//
// Parent implements [PurePath].
func parent[P Kind](p P) Dir {
	s := string(p)
	s = filepath.Clean(s)
	s = filepath.Dir(s)
	return Dir(s)
}

func ancestors[P Kind](p P) iter.Seq[Dir] {
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

func baseName[P Kind](p P) string {
	return filepath.Base(string(p))
}

func ext[P Kind](p P) string {
	return filepath.Ext(string(p))
}

func isLocal[P Kind](p P) bool {
	return filepath.IsLocal(string(p))
}

func isAbsolute[P Kind](p P) bool {
	return filepath.IsAbs(string(p))
}

func clean[P Kind](p P) P {
	return P(filepath.Clean(string(p)))
}

func abs[P Kind](p P) (P, error) {
	q, err := filepath.Abs(string(p))
	return P(q), err
}

func localize[P Kind](p P) (P, error) {
	q, err := filepath.Localize(string(p))
	return P(q), err
}

// See [path/filepath.Rel]
func rel[P Kind](base Dir, p P) (P, error) {
	result, err := filepath.Rel(string(base), string(p))
	return P(result), err
}

func expandUser[P Kind](p P) (q P, err error) {
	if len(p) == 0 || p[0] != '~' || (len(p) > 1 && !os.IsPathSeparator(p[1])) {
		q = p
		return
	}
	var home Dir
	if home, err = UserHomeDir(); err != nil {
		return
	}
	q = P(home) + p[1:]
	return
}

func chmod[P Kind](p P, mode os.FileMode) error {
	return os.Chmod(string(p), mode)
}

func chown[P Kind](p P, uid int, gid int) error {
	return os.Chown(string(p), uid, gid)
}

func rename[P Kind](p P, newPath PathStr) (result P, err error) {
	result = p
	err = os.Rename(string(p), string(newPath))
	if err == nil {
		result = P(newPath)
	}
	return
}

func remove[P Kind](p P) error {
	return os.Remove(string(p))
}

// See [os.RemoveAll]
func removeAll[P Kind](p P) (P, error) {
	return p, os.RemoveAll(string(p))
}
