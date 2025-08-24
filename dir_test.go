package pathlib_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExampleDir_Eq() {
	cwd := expect(pathlib.Cwd())
	fmt.Println(cwd.Eq("."))
	fmt.Println(pathlib.Dir("/foo").Eq("/foo"))
	fmt.Println(cwd.Join("./relative").Eq("relative"))
	fmt.Println(pathlib.Dir("/a/b").Eq("b"))
	// output:
	// true
	// true
	// true
	// false
}

func ExampleDir_Join() {
	result := pathlib.Dir("/tmp").Join("a/b")
	fmt.Printf("%s is a %T", result, result)
	// output:
	// /tmp/a/b is a pathlib.PathStr
}

func ExampleDir_Glob() {
	demoDir := pathlib.
		TempDir().
		Join("glob-example").
		AsDir()
	demoDir = expect(demoDir.MakeAll(0o777, 0o777))

	defer func() { _, _ = demoDir.RemoveAll() }()

	for _, name := range []string{"x", "y", "z", "a", "b", "c"} {
		expect(demoDir.Join(name + ".txt").AsFile().Make(0o777))
	}

	for _, match := range expect(demoDir.Glob("*.txt")) {
		fmt.Println(match)
	}

	// output:
	// /tmp/glob-example/a.txt
	// /tmp/glob-example/b.txt
	// /tmp/glob-example/c.txt
	// /tmp/glob-example/x.txt
	// /tmp/glob-example/y.txt
	// /tmp/glob-example/z.txt
}

func TestDir_glob(t *testing.T) {
	_, err := pathlib.TempDir().Glob("[*")
	if err == nil {
		t.Fatal("somehow, invalid glob syntax made it through")
	}
}

func ExampleDir_Lstat() {
	tmpDir := expect(pathlib.TempDir().
		Join("dir-lstat-example").
		AsDir().
		MakeAll(0o777, 0o777))
	defer func() { _, _ = tmpDir.RemoveAll() }()

	nonExistent := tmpDir.Join("dir").AsDir()
	_, err := nonExistent.Lstat()
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("%T(%q).Lstat() => fs.ErrNotExist\n", nonExistent, nonExistent)
	}

	dir := expect(nonExistent.Make(0o777))

	file := tmpDir.Join("file.txt").AsFile()
	expect(file.Make(0o755))

	_, err = pathlib.Dir(file.String()).Lstat()
	if e, ok := err.(pathlib.WrongTypeOnDisk[pathlib.Dir]); ok {
		fmt.Println(e.Error())
	}

	link := expect(tmpDir.Join("link").AsSymlink().LinkTo(pathlib.PathStr(dir)))
	disguised := pathlib.Dir(link)
	unmasked := expect(disguised.OnDisk()) // works
	for _, name := range []string{"a", "b", "c"} {
		expect(dir.Join(name).AsFile().Make(0666))
	}
	fmt.Println(unmasked.Path())
	for _, entry := range expect(unmasked.Path().Read()) {
		fmt.Println(entry)
	}
	// Output:
	// pathlib.Dir("/tmp/dir-lstat-example/dir").Lstat() => fs.ErrNotExist
	// pathlib.Dir("/tmp/dir-lstat-example/file.txt") unexpectedly has mode -rwxr-xr-x on-disk
	// /tmp/dir-lstat-example/link
	// - a
	// - b
	// - c
}

func ExampleDir_Walk() {
	dir := expect(pathlib.TempDir().Join("dir-walk-example").AsDir().Make(0755))
	defer func() { _, _ = dir.RemoveAll() }()

	_ = expect(dir.Join("foo/bar/baz").AsDir().MakeAll(0755, 0755))
	_ = expect(dir.Join("a/b/c").AsDir().MakeAll(0755, 0755))

	err := dir.Walk(func(path pathlib.PathStr, d fs.DirEntry, err error) error {
		fmt.Println(expect(path.Rel(dir)))
		return nil
	})
	if err != nil {
		panic(err)
	}

	// output:
	// .
	// a
	// a/b
	// a/b/c
	// foo
	// foo/bar
	// foo/bar/baz
}

// TODO: find a reliably unreadable directory to test
// func TestDir_failRead(t *testing.T) {
// 	dir := pathlib.TempDir().Join("dir-walk-example").AsDir().Make(0333))
// 	defer func() { dir.RemoveAll()) }()

// 	dir.Join("foo/bar/baz").AsDir().MakeAll(0755, 0755))
// 	dir.Join("a/b/c").AsDir().MakeAll(0755, 0755)

// 	_, err := dir.Read()
// 	if err == nil {
// 		t.Fatal("expected read error")
// 	}
// }

func TestDir_Chdir(t *testing.T) {
	cwd := expect(pathlib.Cwd())
	defer func() { expect(cwd.Chdir()) }()
	dir := expect(
		pathlib.TempDir().
			Join("chdir-example/a/b/c").
			AsDir().
			MakeAll(0755, 0755),
	)

	expect(dir.Chdir())
	newCwd, _ := os.Getwd()
	if newCwd != dir.String() {
		t.Fatal("failed to chdir into " + dir)
	}
}

func ExampleDir_purePath() {
	d := pathlib.Dir("~/.config/git/..")

	fmt.Println("On Unix")
	fmt.Printf("%T(%q)\n", d, d)
	method := func(name string, val any) {
		fmt.Printf("\t.%s() => %T(%#v)\n", name, val, val)
	}
	method("Parts", d.Parts())
	method("Clean", d.Clean())
	method("Parent", d.Parent())
	method("BaseName", d.BaseName())
	method("IsAbsolute", d.IsAbsolute())
	method("IsLocal", d.IsLocal())
	method("Ext", d.Ext())
	// method("Abs().Unwrap", d.Abs()))
	// method("ExpandUser().Unwrap", d.ExpandUser()))

	// Output:
	// On Unix
	// pathlib.Dir("~/.config/git/..")
	// 	.Parts() => []string([]string{"~", ".config", "git", ".."})
	// 	.Clean() => pathlib.Dir("~/.config")
	// 	.Parent() => pathlib.Dir("~")
	// 	.BaseName() => string("..")
	// 	.IsAbsolute() => bool(false)
	// 	.IsLocal() => bool(true)
	// 	.Ext() => string(".")
}


func ExampleDir_Localize() {
	fmt.Println("On Unix")
	fmt.Println("Localized", expect(pathlib.Dir("foo/bar/baz").Localize()))
	// Output:
	// On Unix
	// Localized foo/bar/baz
}

func ExampleDir_Abs() {
	cwd := expect(pathlib.Cwd())
	roundTrip := func(d pathlib.Dir) {
		abs := expect(d.Abs())
		rel := expect(abs.Rel(cwd))
		if rel.Clean() != d.Clean() {
			fmt.Printf("dir=%q\nabs=%q\nrel=%q\n\n", d, abs, rel)
		} else {
			fmt.Println("true")
		}
	}
	roundTrip("foo/bar")
	roundTrip("./bar")

	// Output:
	// true
	// true
}

func TestDir_ExpandUser(t *testing.T) {
	dir := pathlib.Dir("~/foo/bar")
	intermediate := expect(dir.ExpandUser())
	roundTrip := expect(intermediate.Rel(expect(pathlib.UserHomeDir())))
	if roundTrip != "foo/bar" {
		t.Fail()
	}
}
func TestDir_badStat(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	defer func() { expect(temp.RemoveAll()) }()

	expect(temp.Join("file.txt").AsFile().Make(0666))
	expect(temp.Join("link").AsSymlink().LinkTo(temp.Join("file.txt")))
	_, err := temp.Join("link").AsDir().OnDisk()
	if _, ok := err.(pathlib.WrongTypeOnDisk[pathlib.Dir]); !ok {
		t.Fail()
	}
}

func TestDir_remove(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	dir := expect(temp.Join("foo").AsDir().Make(0777))
	renamed := expect(dir.Rename(temp.Join("bar")))
	if dir.Exists() {
		t.Fail()
	}
	enforce(renamed.Remove())
}

func testChmod[P interface {
	pathlib.Kind
	pathlib.Beholder[P]
	pathlib.Changer
}](t *testing.T, p P, mode os.FileMode) {
	if err := p.Chmod(mode); err != nil {
		t.Fatal(err)
	}
	perm := expect(p.Lstat()).Mode() & fs.ModePerm
	if perm&mode != perm {
		t.Fatalf("expected %o, got %o", mode, perm)
	}
}

func TestDir_chmod(t *testing.T) {
	dir := expect(pathlib.Dir(t.TempDir()).Join("dir").AsDir().Make(0755))
	testChmod(t, dir, 0777)
}

func TestDir_chown(t *testing.T) {
	err := expect(pathlib.Dir(t.TempDir()).
		Join(t.Name()).
		AsDir().
		Make(0755)).
		Chown(0, 0)

	if err != nil {
		if _, ok := err.(*fs.PathError); !ok {
			t.Fatalf("expected *fs.PathError, got %T", err)
		}
	}

}

func testChown[P interface {
	pathlib.Kind
	pathlib.Changer
	pathlib.Beholder[P]
}](t *testing.T, p P) {
	switch strings.Split(runtime.GOOS, "/")[0] {
	case "windows", "plan9":
		t.Skip()
	}
	user, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		t.Fatal(err)
	}
	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if err = p.Chown(uid, gid); err != nil {
		t.Fatal(err)
	}
}

func TestDir_chown_self(t *testing.T) {
	dir := expect(pathlib.Dir(t.TempDir()).
		Join(t.Name()).
		AsDir().
		Make(0755))
	testChown(t, dir)
}

func TestOnDisk_chown(t *testing.T) {
	dir := expect(pathlib.Dir(t.TempDir()).
		Join(t.Name()).
		AsDir().
		Make(0755))
	err := expect(dir.OnDisk()).Chown(0, 0)
	if err != nil {
		if _, ok := err.(*fs.PathError); !ok {
			t.Fatalf("expected *fs.PathError, got %T", err)
		}
	}
}

func TestOnDisk_chown_self(t *testing.T) {
	dir := pathlib.Dir(t.TempDir()).
		Join(t.Name()).
		AsDir()
	dir = expect(dir.Make(0755))
	onDisk := expect(dir.OnDisk())
	switch strings.Split(runtime.GOOS, "/")[0] {
	case "windows", "plan9":
		t.Skip()
	}
	user, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		t.Fatal(err)
	}
	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if err = onDisk.Chown(uid, gid); err != nil {
		t.Fatal(err)
	}
}

func TestDir_makeAll_fail(t *testing.T) {
	d, err := pathlib.Dir("/foo/bar").MakeAll(0755, 0755)
	if err == nil {
		enforce(d.Remove())
		enforce(d.Parent().Remove())
		t.Error("expected error from making /foo/bar")
	}
}
