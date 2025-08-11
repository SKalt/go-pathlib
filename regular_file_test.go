package pathlib_test

import (
	"fmt"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExampleFile_purePath() {
	d := pathlib.File("~/.config/git/..")

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
	// pathlib.File("~/.config/git/..")
	// 	.BaseName() => string("..")
	// 	.IsAbsolute() => bool(false)
	// 	.IsLocal() => bool(true)
	// 	.Ext() => string(".")
	// 	.Parent() => pathlib.Dir("~/.config/git")
	// 	.Parts() => []string([]string{"~", ".config", "git", ".."})
	// 	.Clean() => pathlib.File("~/.config")
}

func TestFile_Join(t *testing.T) {
	assertEq := func(actual, expected pathlib.PathStr) {
		if actual != expected {
			t.Fatalf("- %q\n+ %q\n", expected, actual)
		}
	}
	assertEq(pathlib.File("~/foo").Join("bar"), "~/foo/bar")
}

func TestFile_manipulator(t *testing.T) {
	prevWorkingDirectory := pathlib.Cwd().Unwrap()
	dir := pathlib.Dir(t.TempDir())
	dir.Chdir().Unwrap()
	defer func() {
		prevWorkingDirectory.Chdir()
		dir.RemoveAll().Unwrap()
	}()

	file := dir.Join("foo.txt").AsFile()
	file.Make(0655).Unwrap()
	if !file.Exists() {
		t.Fail()
	}
	file.Rename("foo.sh").Unwrap().Chmod(0777).Unwrap().Remove()
	if file.Exists() {
		t.Fail()
	}
}

func TestFile_Transformer(t *testing.T) {
	var file pathlib.Transformer[pathlib.File] = pathlib.File("~/foo/bar/../baz.txt")
	file.Abs().Unwrap()
	file.Rel("~").Unwrap()
	if file.Clean() != "~/foo/baz.txt" {
		t.Fatal("clean", file.Clean())
	}
	if !file.ExpandUser().Unwrap().IsAbsolute() {
		t.Fatal("ExpandUser", file.ExpandUser().Unwrap())
	}
	if !file.Eq(file.Clean()) {
		t.Fatal(file.Abs().Unwrap(), file.Clean().Abs().Unwrap())
	}
}

func TestFile_Beholder(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt").AsFile()
	if file.Exists() {
		t.Fail()
	}

	file.Make(0666).Unwrap()
	file.Stat().Unwrap()
	file.Lstat().Unwrap()
	file.OnDisk().Unwrap()

}
