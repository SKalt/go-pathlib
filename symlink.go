package pathlib

import (
	"errors"
	"os"
)

type Symlink PathStr

var (
	_ PurePath                = Symlink("")
	_ Transformer[Symlink] = Symlink("")
	_ Readable[PathStr]       = Symlink("")
)

// Read implements Readable.
func (s Symlink) Read() (PathStr, error) {
	link, err := os.Readlink(string(s))
	return PathStr(link), err
}

// BaseName implements PurePath.
func (s Symlink) BaseName() string {
	return PathStr(s).BaseName()
}

// IsAbsolute implements PurePath.
func (s Symlink) IsAbsolute() bool {
	return PathStr(s).IsAbsolute()
}

// IsLocal implements PurePath.
func (s Symlink) IsLocal() bool {
	return PathStr(s).IsLocal()
}

// Join implements [PurePath].
func (s Symlink) Join(parts ...string) PathStr {
	return PathStr(s).Join(parts...)
}

// NearestDir implements PurePath.
func (s Symlink) NearestDir() Dir {
	return PathStr(s).NearestDir()
}

// Parent implements PurePath.
func (s Symlink) Parent() Dir {
	return PathStr(s).Parent()
}

// Abs implements [PurePath].
func (s Symlink) Abs(cwd Dir) Symlink {
	return Symlink(PathStr(s).Abs(cwd)) // this will panic if cwd is not absolute
}

// Localize implements PurePath.
func (s Symlink) Localize() Symlink {
	return Symlink(PathStr(s).Localize())
}

// Rel implements PurePath.
func (s Symlink) Rel(target Dir) (Symlink, error) {
	result, err := PathStr(s).Rel(target)
	return Symlink(result), err
}

func (s Symlink) Ext() string {
	return PathStr(s).Ext()
}

func (s Symlink) OnDisk() (*OnDisk[Symlink], error) {
	onDisk, err := PathStr(s).OnDisk()
	if err != nil {
		return nil, err
	}
	if !isSymLink(onDisk.Mode()) {
		return nil, errors.New("not a symlink: " + onDisk.Name())
	}
	return &OnDisk[Symlink]{*onDisk}, nil
}
