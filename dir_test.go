package pathlib_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExampleDir_Eq() {
	cwd := pathlib.Cwd().Unwrap()
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
		AsDir().
		MakeAll(0o777, 0o777).Unwrap()

	defer func() { _ = demoDir.RemoveAll() }()

	for _, name := range []string{"x", "y", "z", "a", "b", "c"} {
		demoDir.Join(name + ".txt").AsFile().Make(0o777).Unwrap()
	}

	for _, match := range demoDir.Glob("*.txt").Unwrap() {
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
	_, err := pathlib.
		TempDir().Glob("[*").Unpack()
	if err == nil {
		t.Fatal("somehow, invalid glob syntax made it through")
	}
}

func ExampleDir_Lstat() {
	tmpDir := pathlib.TempDir().
		Join("dir-lstat-example").
		AsDir().
		MakeAll(0o777, 0o777).
		Unwrap()
	defer func() { tmpDir.RemoveAll() }()

	nonExistent := tmpDir.Join("dir").AsDir()
	_, err := nonExistent.Lstat().Unpack()
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("%T(%q).Lstat() => fs.ErrNotExist\n", nonExistent, nonExistent)
	}

	dir := nonExistent.Make(0o777).Unwrap()

	file := tmpDir.Join("file.txt").AsFile()
	file.Make(0o755).Unwrap()

	_, err = pathlib.Dir(file.String()).Lstat().Unpack()
	if e, ok := err.(pathlib.WrongTypeOnDisk[pathlib.Dir]); ok {
		fmt.Println(e.Error())
	}

	link := tmpDir.Join("link").AsSymlink().LinkTo(pathlib.PathStr(dir)).Unwrap()
	disguised := pathlib.Dir(link)
	unmasked := disguised.OnDisk().Unwrap() // works
	for _, name := range []string{"a", "b", "c"} {
		dir.Join(name).AsFile().Make(0666).Unwrap()
	}
	fmt.Println(unmasked.Path())
	for _, entry := range unmasked.Path().Read().Unwrap() {
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
	dir := pathlib.TempDir().Join("dir-walk-example").AsDir().Make(0755).Unwrap()
	defer func() { dir.RemoveAll().Unwrap() }()

	dir.Join("foo/bar/baz").AsDir().MakeAll(0755, 0755).Unwrap()
	dir.Join("a/b/c").AsDir().MakeAll(0755, 0755)

	err := dir.Walk(func(path pathlib.PathStr, d fs.DirEntry, err error) error {
		fmt.Println(path.Rel(dir).Unwrap())
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
// 	dir := pathlib.TempDir().Join("dir-walk-example").AsDir().Make(0333).Unwrap()
// 	defer func() { dir.RemoveAll().Unwrap() }()

// 	dir.Join("foo/bar/baz").AsDir().MakeAll(0755, 0755).Unwrap()
// 	dir.Join("a/b/c").AsDir().MakeAll(0755, 0755)

// 	_, err := dir.Read().Unpack()
// 	if err == nil {
// 		t.Fatal("expected read error")
// 	}
// }

func TestDir_Chdir(t *testing.T) {
	cwd := pathlib.Cwd().Unwrap()
	defer func() { cwd.Chdir().Unwrap() }()
	dir := pathlib.TempDir().
		Join("chdir-example/a/b/c").
		AsDir().
		MakeAll(0755, 0755).
		Unwrap()

	dir.Chdir().Unwrap()
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
	method("BaseName", d.BaseName())
	method("IsAbsolute", d.IsAbsolute())
	method("IsLocal", d.IsLocal())
	method("Ext", d.Ext())
	method("Parent", d.Parent())
	method("Parts", d.Parts())
	method("Clean", d.Clean())
	// method("Abs().Unwrap", d.Abs().Unwrap())
	// method("ExpandUser().Unwrap", d.ExpandUser().Unwrap())

	// Output:
	// On Unix
	// pathlib.Dir("~/.config/git/..")
	// 	.BaseName() => string("..")
	// 	.IsAbsolute() => bool(false)
	// 	.IsLocal() => bool(true)
	// 	.Ext() => string(".")
	// 	.Parent() => pathlib.Dir("~/.config/git")
	// 	.Parts() => []string([]string{"~", ".config", "git", ".."})
	// 	.Clean() => pathlib.Dir("~/.config")
}

// Adapted from https://pkg.go.dev/path/filepath#example-Rel
// func TestDir_Rel(t *testing.T) {
// 	dirs := []pathlib.Dir{
// 		"/a/b/c",
// 		"/b/c",
// 		"./b/c",
// 	}
// 	var base pathlib.Dir = "/a"

// 	fmt.Println("On Unix:")
// 	for _, d := range dirs {
// 		rel, err := d.Rel(base).Unpack()
// 		fmt.Printf("%q.Rel(%q) => %q %v\n", d, base, rel, err)
// 	}
// 	// Output:
// 	// 	"/a/b/c": "b/c" <nil>
// 	// "/b/c": "../b/c" <nil>
// 	// "./b/c": "" Rel: can't make ./b/c relative to /a
// }

func ExampleDir_transformer() {
	fmt.Println("On Unix")
	fmt.Println("Localized", pathlib.Dir("foo/bar/baz").Localize().Unwrap())
	// Output:
	// On Unix
	// Localized foo/bar/baz
}

func ExampleDir_Abs() {
	cwd := pathlib.Cwd().Unwrap()
	roundTrip := func(d pathlib.Dir) {
		abs := d.Abs().Unwrap()
		rel := abs.Rel(cwd).Unwrap()
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
	roundTrip := pathlib.Dir("~/foo/bar").
		ExpandUser().
		Unwrap().
		Rel(pathlib.UserHomeDir().Unwrap()).
		Unwrap()
	if roundTrip != "foo/bar" {
		t.Fail()
	}
}
func TestDir_badStat(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	defer func() { temp.RemoveAll().Unwrap() }()

	temp.Join("file.txt").AsFile().Make(0666).Unwrap()
	temp.Join("link").AsSymlink().LinkTo(temp.Join("file.txt")).Unwrap()
	_, err := temp.Join("link").AsDir().OnDisk().Unpack()
	if _, ok := err.(pathlib.WrongTypeOnDisk[pathlib.Dir]); !ok {
		t.Fail()
	}
}

func TestDir_remove(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	dir.
		Join(t.Name()).
		AsDir().
		Make(0777).Unwrap().
		Rename(dir.Join(t.Name() + "__foo")).Unwrap().
		Remove().Unwrap()
}

func TestDir_chmod(t *testing.T) {
	dir := pathlib.Dir(t.TempDir()).Join("dir").AsDir().Make(0755).Unwrap()
	if dir.Chmod(0777).Unwrap().Stat().Unwrap().Mode()&fs.ModePerm != 0777 {
		t.Fail()
	}
}
