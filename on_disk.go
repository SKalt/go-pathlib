package pathlib

import (
	"io/fs"
	"os"
	"time"
)

var (
	_ fs.FileInfo = OnDisk[PathStr]{}
	_ PurePath    = OnDisk[PathStr]{}
)

func (p OnDisk[P]) Parent() Dir {
	return p.Path().Parent()
}

func (p OnDisk[P]) Path() P {
	return P(p.info.Name())
}

func (p *OnDisk[P]) Chmod(mode os.FileMode) error {
	return os.Chmod(p.info.Name(), mode)
}

func (p OnDisk[P]) IsRegular() (isRegular bool) {
	return p.info.Mode().IsRegular()
}

// BaseName implements PurePath.
func (p OnDisk[P]) BaseName() string {
	return p.Path().BaseName()
}

// Ext implements PurePath.
func (p OnDisk[P]) Ext() string {
	return p.Path().Ext()
}

// IsAbsolute implements PurePath.
func (p OnDisk[P]) IsAbsolute() bool {
	return p.Path().IsAbsolute()
}

// IsLocal implements PurePath.
func (p OnDisk[P]) IsLocal() bool {
	return p.Path().IsLocal()
}

// Join implements PurePath.
func (p OnDisk[P]) Join(parts ...string) PathStr {
	return p.Path().Join(parts...)
}

// NearestDir implements PurePath.
func (p OnDisk[P]) NearestDir() Dir {
	return p.Path().NearestDir()
}

// IsDir implements fs.FileInfo.
func (p OnDisk[P]) IsDir() bool {
	return p.info.IsDir()
}

// ModTime implements fs.FileInfo.
func (p OnDisk[P]) ModTime() time.Time {
	return p.info.ModTime()
}

// Mode implements fs.FileInfo.
func (p OnDisk[P]) Mode() fs.FileMode {
	return p.info.Mode()
}

// Name implements fs.FileInfo.
func (p OnDisk[P]) Name() string {
	return p.info.Name()
}

// Size implements fs.FileInfo.
func (p OnDisk[P]) Size() int64 {
	return p.info.Size()
}

// Sys implements fs.FileInfo.
func (p OnDisk[P]) Sys() any {
	return p.info.Sys()
}
