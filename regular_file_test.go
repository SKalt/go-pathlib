package pathlib_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExampleFile_purePath() {
	f := pathlib.File("~/.config/tool/../other-tool/config.toml")

	fmt.Println("On Unix")
	fmt.Printf("%T(%q)\n", f, f)
	method := func(name string, val any) {
		fmt.Printf("\t.%s() => %T(%#v)\n", name, val, val)
	}
	method("Parts", f.Parts())
	method("Clean", f.Clean())
	method("Parent", f.Parent())
	method("BaseName", f.BaseName())
	method("IsAbsolute", f.IsAbsolute())
	method("IsLocal", f.IsLocal())
	method("Ext", f.Ext())
	// method("Abs().Unwrap", d.Abs()))
	// method("ExpandUser().Unwrap", d.ExpandUser()))

	// Output:
	// On Unix
	// pathlib.File("~/.config/tool/../other-tool/config.toml")
	// 	.Parts() => []string([]string{"~", ".config", "tool", "..", "other-tool", "config.toml"})
	// 	.Clean() => pathlib.File("~/.config/other-tool/config.toml")
	// 	.Parent() => pathlib.Dir("~/.config/other-tool")
	// 	.BaseName() => string("config.toml")
	// 	.IsAbsolute() => bool(false)
	// 	.IsLocal() => bool(true)
	// 	.Ext() => string(".toml")
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
	if err := renamed.Chmod(0777); err != nil {
		t.Fatal(err)
	}
	if file.Exists() {
		t.Fail()
	}
	if !renamed.Exists() {
		t.Fail()
	}
	enforce(renamed.Remove())
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

}

func TestFile_chown(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt").AsFile()
	expect(file.Make(0666))
	testChown(t, file)
}

func TestFile_open(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt").AsFile()
	if _, err := file.Open(os.O_RDONLY, 0666); err == nil {
		t.Fatal("expected error opening file that does not exist")
	}

}
func TestFile_remake(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt").AsFile()
	h := expect(file.Make(0666))
	expect(h.WriteString("hello"))
	enforce(h.Close())

	h2 := expect(file.Make(0600))
	enforce(h2.Close())
	if expect(file.Stat()).Mode() == 0600 {
		t.Fail()
	}
	content := string(expect(file.Read()))
	if content != "hello" {
		t.Fail()
	}
}

func TestFile_wrongType(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	actual := expect(temp.Join("dir").AsDir().Make(0777))
	incognito := pathlib.File(actual)
	info, err := incognito.Stat()
	if err == nil {
		t.Fatalf("got %v", info)
	}
	if _, ok := err.(pathlib.WrongTypeOnDisk[pathlib.File]); !ok {
		t.Fatalf("unexpected type : %T(%v)", err, err)
	}
	info, err = incognito.Lstat()
	if err == nil {
		t.Fatalf("got %v", info.Mode())
	}
	if _, ok := err.(pathlib.WrongTypeOnDisk[pathlib.File]); !ok {
		t.Fatalf("unexpected type : %T(%v)", err, err)
	}
}

func TestFile_localize(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := expect(temp.Join("file.txt").AsFile().Rel(temp))
	localized := expect(file.Localize())
	if localized != file {
		t.Fatalf("expected %q, got %q", file, localized)
	}

}

func TestFile_makeAll_fail(t *testing.T) {
	d, err := pathlib.File("/foo/bar.txt").MakeAll(0644, 0755)
	if err == nil {
		enforce(d.Remove())
		enforce(d.Parent().Remove())
		t.Error("expected error from making /foo/bar.txt")
	}

}
