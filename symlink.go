package pathlib

import (
	"os"
)



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

// Join implements PurePath.
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
