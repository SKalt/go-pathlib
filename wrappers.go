package pathlib

import (
	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

func stat[P Kind](p P) Result[OnDisk[P]] {
	info, err := os.Stat(string(p))
	return Result[OnDisk[P]]{onDisk[P]{p, info}, err}
}

func lstat[P Kind](p P) Result[OnDisk[P]] {
	info, err := os.Lstat(string(p))
	if errors.Is(err, fs.ErrNotExist) {
		return Result[OnDisk[P]]{nil, err}
	}
	return Result[OnDisk[P]]{onDisk[P]{p, info}, err}
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

func abs[P Kind](p P) Result[P] {
	q, err := filepath.Abs(string(p))
	return Result[P]{P(q), err}
}

func localize[P Kind](p P) Result[P] {
	q, err := filepath.Localize(string(p))
	return Result[P]{P(q), err}
}

func rel[P Kind](base Dir, p P) Result[P] {
	result, err := filepath.Rel(string(base), string(p))
	if err != nil {
		return Result[P]{"", errors.Join(err, errors.New("unable to make relative path"))}
	}
	return Result[P]{P(result), nil}
}

func expandUser[P Kind](p P) Result[P] {
	if len(p) == 0 || p[0] != '~' || (len(p) > 1 && !os.IsPathSeparator(p[1])) {
		return Result[P]{p, nil}
	}

	home := UserHomeDir()
	if home.err != nil {
		return Result[P]{"", home.err}
	}
	return Result[P]{P(P(home.val) + p[1:]), nil}
}

func chmod[P Kind](p P, mode os.FileMode) Result[P] {
	return Result[P]{p, os.Chmod(string(p), mode)}
}

func chown[P Kind](p P, uid int, gid int) Result[P] {
	return Result[P]{p, os.Chown(string(p), uid, gid)}
}

func rename[P Kind](p P, newPath PathStr) Result[P] {
	err := os.Rename(string(p), string(newPath))
	if err != nil {
		return Result[P]{"", err}
	}
	return Result[P]{p, nil}
}

// See [os.RemoveAll]
func removeAll[P Kind](p P) Result[P] {
	return Result[P]{p, os.RemoveAll(string(p))}
}
