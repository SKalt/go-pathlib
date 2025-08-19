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
	// method("Abs().Unwrap", d.Abs()))
	// method("ExpandUser().Unwrap", d.ExpandUser()))

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
	prevWorkingDirectory := expect(pathlib.Cwd())
	dir := pathlib.Dir(t.TempDir())
	expect(dir.Chdir())
	defer func() {
		expect(dir.RemoveAll())
		expect(prevWorkingDirectory.Chdir())
	}()

	file := dir.Join("foo.txt").AsFile()
	expect(file.Make(0655))
	if !file.Exists() {
		t.Fail()
	}
	renamed := expect(file.Rename("foo.sh"))
	renamed = expect(renamed.Chmod(0777))
	if file.Exists() {
		t.Fail()
	}
	if !renamed.Exists() {
		t.Fail()
	}
	expect(renamed.Remove())
}

func TestFile_Transformer(t *testing.T) {
	var file pathlib.Transformer[pathlib.File] = pathlib.File("~/foo/bar/../baz.txt")
	expect(file.Abs())
	expect(file.Rel("~"))
	if file.Clean() != "~/foo/baz.txt" {
		t.Fatal("clean", file.Clean())
	}
	if !expect(file.ExpandUser()).IsAbsolute() {
		t.Fatal("ExpandUser", expect(file.ExpandUser()))
	}
	if !file.Eq(file.Clean()) {
		t.Fatal(expect(file.Abs()), expect(file.Clean().Abs()))
	}
}

func TestFile_Beholder(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt").AsFile()
	if file.Exists() {
		t.Fail()
	}

	expect(file.Make(0666))
	expect(file.Stat())
	expect(file.Lstat())
	expect(file.OnDisk())

}

func TestFile_chown(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt").AsFile()
	expect(file.Make(0666))
	testChown(t, file)
}
