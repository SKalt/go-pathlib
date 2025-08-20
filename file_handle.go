package pathlib

import (
	"io"
	"io/fs"
	"os"
	"syscall"
	"time"
)

// An open file descriptor. Unlike an [os.File], it can only represent a logical
// file (as in a document on-disk), never a directory.
type Handle struct {
	inner *os.File
}

// PurePath --------------------------------------------------------------------
var _ PurePath = &Handle{}

func (h *Handle) Path() File {
	return File(h.inner.Name())
}

// BaseName implements PurePath.
func (h *Handle) BaseName() string {
	return h.Path().BaseName()
}

// Ext implements PurePath.
func (h *Handle) Ext() string {
	return h.Path().Ext()
}

// IsAbsolute implements PurePath.
func (h *Handle) IsAbsolute() bool {
	return h.Path().IsAbsolute()
}

// IsLocal implements PurePath.
func (h *Handle) IsLocal() bool {
	return h.Path().IsLocal()
}

// Join implements PurePath.
func (h *Handle) Join(segments ...string) PathStr {
	return h.Path().Join(segments...)
}

// Parent implements PurePath.
func (h *Handle) Parent() Dir {
	return h.Path().Parent()
}

// Parts implements PurePath.
func (h *Handle) Parts() []string {
	return h.Path().Parts()
}

// Beholder --------------------------------------------------------------------

var _ Beholder[File] = &Handle{}

// Lstat implements Beholder.
func (h *Handle) Lstat() (OnDisk[File], error) {
	return h.Path().Lstat()
}

// OnDisk implements Beholder.
func (h *Handle) OnDisk() (OnDisk[File], error) {
	return h.Path().OnDisk()
}

// Stat implements Beholder.
func (h *Handle) Stat() (OnDisk[File], error) {
	info, err := h.inner.Stat()
	if err != nil {
		return nil, err
	}
	result := onDisk[File]{h.Path(), info}
	if !info.Mode().IsRegular() {
		return result, WrongTypeOnDisk[File]{result}
	}
	return result, nil
}

// Exists implements Beholder.
func (h *Handle) Exists() bool {
	return h.Path().Exists()
}

// sorta: Manipulator ----------------------------------------------------------

// Chmod implements Manipulator.
func (h *Handle) Chmod(mode fs.FileMode) (*Handle, error) {
	return h, h.inner.Chmod(mode)
}

// Chown implements Manipulator.
func (h *Handle) Chown(uid int, gid int) (*Handle, error) {
	return h, h.inner.Chown(uid, gid)
}

// Remove implements Manipulator.
func (h *Handle) Remove() (File, error) {
	if err := h.Close(); err != nil {
		return h.Path(), err
	}
	return h.Path().Remove()
}

// Rename implements Manipulator.
func (h *Handle) Rename(newPath PathStr) (File, error) {
	if err := h.Close(); err != nil {
		return h.Path(), err
	}
	return h.Path().Rename(newPath)
}

// retained from *os.File ------------------------------------------------------

// See [os.File.Close].
func (h *Handle) Close() error {
	return h.inner.Close()
}

// See [os.File.Fd].
func (h *Handle) Fd() uintptr {
	return h.inner.Fd()
}

// See [os.File.Name].
func (h *Handle) Name() string {
	return h.inner.Name()
}

// See [os.File.Read].
func (h *Handle) Read(p []byte) (n int, err error) {
	return h.inner.Read(p)
}

// See [os.File.ReadAt].
func (h *Handle) ReadAt(p []byte, off int64) (n int, err error) {
	return h.inner.ReadAt(p, off)
}

// See [os.File.ReadFrom].
func (h *Handle) ReadFrom(r io.Reader) (n int64, err error) {
	return h.inner.ReadFrom(r)
}

// See [os.File.Seek].
func (h *Handle) Seek(offset int64, whence int) (int64, error) {
	return h.inner.Seek(offset, whence)
}

// See [os.File.SetDeadline].
func (h *Handle) SetDeadline(deadline time.Time) error {
	return h.inner.SetDeadline(deadline)
}

// See [os.File.SetReadDeadline].
func (h *Handle) SetReadDeadline(deadline time.Time) error {
	return h.inner.SetReadDeadline(deadline)
}

// See [os.File.Sync].
func (h *Handle) Sync() error {
	return h.inner.Sync()
}

// See [os.File.SyscallConn].
func (h *Handle) SyscallConn() (syscall.RawConn, error) {
	return h.inner.SyscallConn()
}

// See [os.File.Truncate].
func (h *Handle) Truncate(size int64) error {
	return h.inner.Truncate(size)
}

// See [os.File.Write].
func (h *Handle) Write(p []byte) (n int, err error) {
	return h.inner.Write(p)
}

// See [os.File.WriteAt].
func (h *Handle) WriteAt(p []byte, off int64) (n int, err error) {
	return h.inner.WriteAt(p, off)
}

// See [os.File.WriteTo].
func (h *Handle) WriteTo(w io.Writer) (n int64, err error) {
	return h.inner.WriteTo(w)
}

var _ io.StringWriter = &Handle{}

// WriteString implements io.StringWriter.
func (h *Handle) WriteString(s string) (n int, err error) {
	return h.inner.WriteString(s)
}

func (h *Handle) SetWriteDeadline(deadline time.Time) error {
	return h.inner.SetWriteDeadline(deadline)
}
