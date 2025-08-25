package pathlib

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"syscall"
	"time"
)

// An open file descriptor. Unlike an [os.File], it can only represent a logical
// file (as in a document on-disk), never a directory.
type FileHandle interface {
	Path() File
	String() string
	PurePath
	Beholder[File]
	Transformer[File]
	Changer
	Remover[File]

	// from *os.File
	Name() string
	Truncate(size int64) error
	SyscallConn() (syscall.RawConn, error)
	SetDeadline(deadline time.Time) error
	SetReadDeadline(deadline time.Time) error
	SetWriteDeadline(deadline time.Time) error
	Fd() uintptr

	io.Closer
	io.Seeker
	io.Reader
	io.Writer
	io.StringWriter
}

type handle struct{ *os.File }


var _ FileHandle = &handle{}

func (h *handle) Path() File {
	return File(h.Name())
}

func (h *handle) String() string {
	return h.Path().String()
}

// PurePath --------------------------------------------------------------------
var _ PurePath = &handle{}

// BaseName implements [PurePath].
func (h *handle) BaseName() string {
	return h.Path().BaseName()
}

// Ext implements [PurePath].
func (h *handle) Ext() string {
	return h.Path().Ext()
}

// IsAbsolute implements [PurePath].
func (h *handle) IsAbsolute() bool {
	return h.Path().IsAbsolute()
}

// IsLocal implements [PurePath].
func (h *handle) IsLocal() bool {
	return h.Path().IsLocal()
}

// Join implements [PurePath].
func (h *handle) Join(segments ...string) PathStr {
	return h.Path().Join(segments...)
}

// Parent implements [PurePath].
func (h *handle) Parent() Dir {
	return h.Path().Parent()
}

// Parts implements [PurePath].
func (h *handle) Parts() []string {
	return h.Path().Parts()
}

// Beholder --------------------------------------------------------------------

var _ Beholder[File] = &handle{}

// Lstat implements [Beholder].
func (h *handle) Lstat() (Info[File], error) {
	info, err := h.Path().Lstat()
	h.closeIfNonexistent(err)
	// FIXME: handle case where h is one or more symlinks pointing to a regular file?
	return info, err
}

// OnDisk implements [Beholder].
func (h *handle) OnDisk() (Info[File], error) {
	info, err := h.Path().OnDisk()
	h.closeIfNonexistent(err)
	return info, err
}
func (h *handle) closeIfNonexistent(err error) {
	if errors.Is(err, fs.ErrNotExist) {
		_ = h.Close()
	}
}

// Stat implements [Beholder].
func (h *handle) Stat() (Info[File], error) {
	info, err := h.Path().Stat()
	// it might be cheaper to use the `h.inner.Stat()` method, but that
	// seems to erroneously report that the file exists if the file has
	// been removed since the handle was opened.
	h.closeIfNonexistent(err)
	if err != nil {
		return nil, err
	}
	result := onDisk[File]{h.Path(), info}
	if !info.Mode().IsRegular() {
		return result, WrongTypeOnDisk[File]{result}
	}
	return result, nil
}

// Exists implements [Beholder].
func (h *handle) Exists() bool {
	return h.Path().Exists()
}

// Changer ---------------------------------------------------------------------
var _ Changer = &handle{}

// Chmod implements [Changer].
func (h *handle) Chmod(mode fs.FileMode) error {
	return h.Chmod(mode)
}

// Chown implements [Changer].
func (h *handle) Chown(uid int, gid int) error {
	return h.Chown(uid, gid)
}

// Mover -----------------------------------------------------------------------

// Remove implements [Remover].
func (h *handle) Remove() error {
	_ = h.Close()
	return h.Path().Remove()
}

// Rename implements [Manipulator].
func (h *handle) Rename(newPath PathStr) (File, error) {
	_ = h.Close()
	return h.Path().Rename(newPath)
}

// Transformer ------------------------------------------------------------------
var _ Transformer[File] = &handle{}

// Abs implements [Transformer].
func (h *handle) Abs() (File, error) {
	return h.Path().Abs()
}

// Clean implements [Transformer].
func (h *handle) Clean() File {
	return h.Path().Clean()
}

// Eq implements [Transformer].
func (h *handle) Eq(other File) bool {
	return h.Path().Eq(other)
}

// ExpandUser implements [Transformer].
func (h *handle) ExpandUser() (File, error) {
	return h.Path().ExpandUser()
}

// Localize implements [Transformer].
func (h *handle) Localize() (File, error) {
	return h.Path().Localize()
}

// Rel implements [Transformer].
func (h *handle) Rel(base Dir) (File, error) {
	return h.Path().Rel(base)
}
