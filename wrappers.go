package pathlib

import (
	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

func stat[P Kind](p P) (OnDisk[P], error) {
	info, err := os.Stat(string(p))
	return onDisk[P]{info}, err
}

func lstat[P Kind](p P) (actual OnDisk[P], err error) {
	var info os.FileInfo
	info, err = os.Lstat(string(p))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	actual = onDisk[P]{info}
	return
}

func join[P Kind](p P, segments ...string) PathStr {
	return PathStr(filepath.Join(append([]string{string(p)}, segments...)...))
}

// a wrapper around [path/filepath.Dir].
//
// Parent implements [PurePath].
func parent[P Kind](p P) Dir {
	return Dir(filepath.Dir(string(p)))
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

func rel[P Kind](p P, target Dir) (P, error) {
	result, err := filepath.Rel(string(p), string(target))
	if err != nil {
		return "", errors.Join(err, errors.New("unable to make relative path"))
	}
	return P(result), nil
}

func expandUser[P Kind](p P) (result P, err error) {
	if len(p) == 0 || p[0] != '~' || (len(p) > 1 && !os.IsPathSeparator(p[1])) {
		result = p
		return
	}

	var home Dir
	home, err = UserHomeDir()
	if err != nil {
		return
	}

	result = P(P(home) + p[1:])
	return
}

func chmod[P Kind](p P, mode os.FileMode) (result P, err error) {
	return p, os.Chmod(string(p), mode)
}

func chown[P Kind](p P, uid int, gid int) (result P, err error) {
	result = p
	err = os.Chown(string(p), uid, gid)
	return
}

func rename[P Kind](p P, newPath PathStr) (result P, err error) {
	result = p
	err = os.Rename(string(p), string(newPath))
	if err != nil {
		return
	}
	result = P(newPath)
	return
}
