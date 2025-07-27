package pathlib

import (
	"errors"
	"io/fs"
	"os"
)

type File PathStr

func (f File) Open(flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(string(f), flag, perm)
}

// func (f File) withOpen(
// 	flag int,
// 	perm fs.FileMode,
// 	callback func(handle *os.File) error,
// ) (err error) {
// 	var handle *os.File
// 	handle, err = f.Open(flag, perm)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		err = errors.Join(handle.Close())
// 	}()
// 	if callback != nil {
// 		err = callback(handle)
// 	}
// 	return
// }

func (f File) withMake(callback func(*os.File) error) (err error) {
	var handle *os.File
	handle, err = f.Make()
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(handle.Close())
	}()
	if callback != nil {
		err = callback(handle)
	}
	return
}

func (f File) Touch() error {
	return f.withMake(nil)
}

// PurePath --------------------------------------------------------------------
var _ PurePath = File("./example.txt")

// BaseName implements PurePath.
func (f File) BaseName() string {
	return PathStr(f).BaseName()
}

// Ext implements PurePath.
func (f File) Ext() string {
	return PathStr(f).Ext()
}

// IsAbsolute implements PurePath.
func (f File) IsAbsolute() bool {
	return PathStr(f).IsAbsolute()
}

// IsLocal implements PurePath.
func (f File) IsLocal() bool {
	return PathStr(f).IsLocal()
}

// Join implements PurePath.
func (f File) Join(parts ...string) PathStr {
	return PathStr(f).Join(parts...)
}

// Parent implements PurePath.
func (f File) Parent() Dir {
	return PathStr(f).Parent()
}

// Parts implements PurePath.
func (f File) Parts() []string {
	return PathStr(f).Parts()
}

// Transformer -----------------------------------------------------------------
var _ Transformer[File] = File("./example")

// Abs implements Transformer.
func (f File) Abs() (File, error) {
	q, err := PathStr(f).Abs()
	return File(q), err
}

// Clean implements Transformer.
func (f File) Clean() File {
	return File(PathStr(f).Clean())
}

// Eq implements Transformer.
func (f File) Eq(other File) bool {
	return PathStr(f).Eq(PathStr(other))
}

// ExpandUser implements Transformer.
func (f File) ExpandUser() (File, error) {
	q, err := PathStr(f).ExpandUser()
	return File(q), err
}

// Localize implements Transformer.
func (f File) Localize() (File, error) {
	q, err := PathStr(f).Localize()
	return File(q), err
}

// Rel implements Transformer.
func (f File) Rel(target Dir) (File, error) {
	q, err := PathStr(f).Rel(target)
	return File(q), err
}

// Beholder --------------------------------------------------------------------
var _ Beholder[File] = File("./example")

// Exists implements Beholder.
func (f File) Exists() bool {
	return PathStr(f).Exists()
}

// Lstat implements Beholder.
func (f File) Lstat() (OnDisk[File], error) {
	info, err := os.Lstat(string(f))
	return onDisk[File]{info, time.Now()}, err
}

// OnDisk implements Beholder.
func (f File) OnDisk() (OnDisk[File], error) {
	return f.Lstat()
}

// Stat implements Beholder.
func (f File) Stat() (OnDisk[File], error) {
	info, err := os.Stat(string(f))
	return onDisk[File]{info, time.Now()}, err
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[File] = File("./example")

// Chmod implements Manipulator.
func (f File) Chmod(mode os.FileMode) (File, error) {
	return f, os.Chmod(string(f), mode)
}

// Chown implements Manipulator.
func (f File) Chown(uid int, gid int) (File, error) {
	return f, os.Chown(string(f), uid, gid)
}

// Remove implements Manipulator.
func (f File) Remove() error {
	return os.Remove(string(f))
}

// Rename implements Manipulator.
func (f File) Rename(newPath PathStr) (File, error) {
	return f, os.Rename(string(f), string(newPath))
}

// Maker -----------------------------------------------------------------------
var _ Maker[*os.File] = File("./example")

// Make implements Maker.
func (f File) Make(perm ...fs.FileMode) (*os.File, error) {
	const defaultMode fs.FileMode = 0777
main:
	{
		switch len(perm) {
		case 0:
			perm = append(perm, defaultMode)
			goto main
		case 1:
			return f.Open(os.O_RDWR|os.O_CREATE, perm[0])
		default:
			_, err := f.Parent().Make(perm[1:]...)
			if err != nil {
				return nil, err
			}
			perm = perm[0:1]
			goto main
		}
	}
}

// MustMake implements Maker.
func (f File) MustMake(perm ...fs.FileMode) *os.File {
	return expect(f.Make(perm...))
}

var _ Readable[[]byte]

// extra functions:
// .Touch() error
