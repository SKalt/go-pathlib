package pathlib

import (
	"errors"
	"os"
)

type Symlink PathStr

// -----------------------------------------------------------------------------
var _ Readable[PathStr] = Symlink("")

// Read implements Readable.
func (s Symlink) Read() (PathStr, error) {
	link, err := os.Readlink(string(s))
	return PathStr(link), err
}

// -----------------------------------------------------------------------------
var _ PurePath = Symlink("")

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

// -----------------------------------------------------------------------------
var _ Transformer[Symlink] = Symlink("")

// Abs implements Transformer.
func (s Symlink) Abs() (Symlink, error) {
	q, err := PathStr(s).Abs()
	return Symlink(q), err
}

// Localize implements Transformer.
func (s Symlink) Localize() (Symlink, error) {
	q, err := PathStr(s).Localize()
	return Symlink(q), err
}

// Rel implements Transformer.
func (s Symlink) Rel(target Dir) (Symlink, error) {
	result, err := PathStr(s).Rel(target)
	return Symlink(result), err
}

func (s Symlink) ExpandUser() (Symlink, error) {
	q, err := PathStr(s).ExpandUser()
	return Symlink(q), err
}

func (s Symlink) Ext() string {
	return PathStr(s).Ext()
}

func (s Symlink) OnDisk() (*onDisk[Symlink], error) {
	actual, err := PathStr(s).OnDisk()
	if err != nil {
		return nil, err
	}
	if !isSymLink(actual.Mode()) {
		return nil, errors.New("not a symlink: " + actual.Name())
	}
	return &onDisk[Symlink]{*actual}, nil
}
