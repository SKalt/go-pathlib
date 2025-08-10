package pathlib_test

import (
	"errors"
	"fmt"
	"io/fs"

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

func ExampleDir_Lstat() {
	tmpDir := pathlib.TempDir().
		Join("dir-lstat-example").
		AsDir().
		MakeAll(0o777, 0o777).
		Unwrap()
	defer tmpDir.RemoveAll()

	nonExistent := tmpDir.Join("dir").AsDir()
	_, err := nonExistent.Lstat().Unpack()
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("%T(%q).Lstat() => fs.ErrNotExist\n", nonExistent, nonExistent)
	}

	dir := nonExistent.Make(0o777).Unwrap()
	onDisk := dir.Lstat().Unwrap()
	fmt.Printf("%T(%q) created with mode %s\n", dir, dir, onDisk.Mode())

	file := tmpDir.Join("file.txt").AsFile()
	file.Make(0o755).Unwrap()

	_, err = pathlib.Dir(file.String()).Lstat().Unpack()

	if e, ok := err.(pathlib.WrongTypeOnDisk[pathlib.Dir]); ok {
		fmt.Println(e.Error())
	}
	// Output:
	// pathlib.Dir("/tmp/dir-lstat-example/dir").Lstat() => fs.ErrNotExist
	// pathlib.Dir("/tmp/dir-lstat-example/dir") created with mode drwxr-xr-x
	// pathlib.Dir("/tmp/dir-lstat-example/file.txt") unexpectedly has mode -rwxr-xr-x on-disk
}
