package pathlib

import (
	"io/fs"
	"os"
	"time"
)

type File PathStr

// See [os.OpenFile].
func (f File) Open(flag int, perm fs.FileMode) (*os.File, error) {
	return os.OpenFile(string(f), flag, perm)
}

// PurePath --------------------------------------------------------------------
var _ PurePath = File("./example.txt")

// BaseName implements [PurePath].
func (f File) BaseName() string {
	return PathStr(f).BaseName()
}

// Ext implements [PurePath].
func (f File) Ext() string {
	return PathStr(f).Ext()
}

// IsAbsolute implements [PurePath].
func (f File) IsAbsolute() bool {
	return PathStr(f).IsAbsolute()
}

// IsLocal implements [PurePath].
func (f File) IsLocal() bool {
	return PathStr(f).IsLocal()
}

// Join implements [PurePath].
func (f File) Join(parts ...string) PathStr {
	return PathStr(f).Join(parts...)
}

// Parent implements [PurePath].
func (f File) Parent() Dir {
	return PathStr(f).Parent()
}

// Parts implements [PurePath].
func (f File) Parts() []string {
	return PathStr(f).Parts()
}

// Transformer -----------------------------------------------------------------
var _ Transformer[File] = File("./example")
var _ InfallibleTransformer[File] = File("./example")

// Abs implements [Transformer].
func (f File) Abs() (File, error) {
	q, err := PathStr(f).Abs()
	return File(q), err
}

// Clean implements [Transformer].
func (f File) Clean() File {
	return File(PathStr(f).Clean())
}

// Eq implements [Transformer].
func (f File) Eq(other File) bool {
	return PathStr(f).Eq(PathStr(other))
}

// ExpandUser implements [Transformer].
func (f File) ExpandUser() (File, error) {
	q, err := PathStr(f).ExpandUser()
	return File(q), err
}

// Localize implements [Transformer].
func (f File) Localize() (File, error) {
	q, err := PathStr(f).Localize()
	return File(q), err
}

// Rel implements [Transformer].
func (f File) Rel(target Dir) (File, error) {
	q, err := PathStr(f).Rel(target)
	return File(q), err
}

// MustExpandUser implements [InfallibleTransformer].
func (f File) MustExpandUser() File {
	return expect(f.ExpandUser())
}

// MustLocalize implements [InfallibleTransformer].
func (f File) MustLocalize() File {
	return expect(f.Localize())
}

// MustMakeAbs implements [InfallibleTransformer].
func (f File) MustMakeAbs() File {
	return expect(f.Abs())
}

// MustMakeRel implements [InfallibleTransformer].
func (f File) MustMakeRel(target Dir) File {
	return expect(f.Rel(target))
}

// Beholder --------------------------------------------------------------------
var _ Beholder[File] = File("./example")
var _ InfallibleBeholder[File] = File("./example")

// Exists implements [Beholder].
func (f File) Exists() bool {
	return PathStr(f).Exists()
}

// Lstat implements [Beholder].
func (f File) Lstat() (OnDisk[File], error) {
	info, err := os.Lstat(string(f))
	return onDisk[File]{info, time.Now()}, err
}

// OnDisk implements [Beholder].
func (f File) OnDisk() (OnDisk[File], error) {
	return f.Lstat()
}

// Stat implements [Beholder].
func (f File) Stat() (OnDisk[File], error) {
	info, err := os.Stat(string(f))
	return onDisk[File]{info, time.Now()}, err
}

// MustBeOnDisk implements [InfallibleBeholder].
func (f File) MustBeOnDisk() OnDisk[File] {
	return expect(f.OnDisk())
}

// MustLstat implements [InfallibleBeholder].
func (f File) MustLstat() OnDisk[File] {
	return expect(f.Lstat())
}

// MustStat implements [InfallibleBeholder].
func (f File) MustStat() OnDisk[File] {
	return expect(f.Stat())
}

// Manipulator -----------------------------------------------------------------
var _ Manipulator[File] = File("./example")
var _ InfallibleManipulator[File] = File("./example")

// Chmod implements [Manipulator].
func (f File) Chmod(mode os.FileMode) (File, error) {
	return f, os.Chmod(string(f), mode)
}

// Chown implements [Manipulator].
func (f File) Chown(uid int, gid int) (File, error) {
	return f, os.Chown(string(f), uid, gid)
}

// Remove implements [Manipulator].
func (f File) Remove() error {
	return os.Remove(string(f))
}

// Rename implements [Manipulator].
func (f File) Rename(newPath PathStr) (File, error) {
	return f, os.Rename(string(f), string(newPath))
}

// MustChmod implements [InfallibleManipulator].
func (f File) MustChmod(perm os.FileMode) File {
	return expect(f.Chmod(perm))
}

// MustChown implements [InfallibleManipulator].
func (f File) MustChown(uid int, gid int) File {
	return expect(f.Chown(uid, gid))
}

// MustRemove implements [InfallibleManipulator].
func (f File) MustRemove() {
	expect[any](nil, f.Remove())
}

// MustRename implements [InfallibleManipulator].
func (f File) MustRename(newPath PathStr) File {
	return expect(f.Rename(newPath))
}

// Maker -----------------------------------------------------------------------
var _ Maker[*os.File] = File("./example")
var _ InfallibleMaker[*os.File] = File("./example")

// Make implements [Maker].
func (f File) Make(perm ...fs.FileMode) (handle *os.File, err error) {
	const defaultMode fs.FileMode = 0o666
	if len(perm) == 0 {
		perm = append(perm, defaultMode)
	} else if len(perm) > 1 {
		if parent := f.Parent(); !parent.Exists() {
			if _, err := parent.Make(perm[1:]...); err != nil {
				return nil, err
			}
		}
	}
	return f.Open(os.O_RDWR|os.O_CREATE, perm[0])
}

// MustMake implements [Maker].
func (f File) MustMake(perm ...fs.FileMode) *os.File {
	return expect(f.Make(perm...))
}

// Readable --------------------------------------------------------------------
var _ Readable[[]byte] = File("./example")
var _ InfallibleReader[[]byte] = File("./example")

func (f File) Read() (data []byte, err error) {
	return os.ReadFile(string(f))
}

func (f File) MustRead() []byte {
	return expect(f.Read())
}

// Destroyer -------------------------------------------------------------------
var _ Destroyer = File("./example")
var _ InfallibleDestroyer = File("./example")

// MustRemoveAll implements [InfallibleDestroyer].
func (f File) MustRemoveAll() {
	PathStr(f).MustRemoveAll()
}

// RemoveAll implements [Destroyer].
func (f File) RemoveAll() error {
	return os.RemoveAll(string(f))
}
