package pathlib_test

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExamplePathStr_IsLocal() {
	p := pathlib.PathStr("/absolute/path/to/file.txt")
	fmt.Println(p.IsLocal()) // false

	p = pathlib.PathStr("relative/path/to/file.txt")
	fmt.Println(p.IsLocal()) // true

	p = pathlib.PathStr("./local/path/to/file.txt")
	fmt.Println(p.IsLocal()) // true
	// Output:
	// false
	// true
	// true
}
func ExamplePathStr_Join() {
	fmt.Println("On Unix:")
	example := func(p pathlib.PathStr, segments ...string) {
		q := p.Join(segments...)
		x := fmt.Sprintf("%q", segments)
		fmt.Printf("%T(%q).Join(%s) => %q\n", p, p, x[1:len(x)-1], q)
	}
	example(pathlib.PathStr("a"), "b", "c")
	example(pathlib.PathStr("a"), "b/c")
	example(pathlib.PathStr("a/b"), "c")
	example(pathlib.PathStr("a/b"), "/c")

	example(pathlib.PathStr("a/b"), "../../../xyz")
	// Output:
	// On Unix:
	// pathlib.PathStr("a").Join("b" "c") => "a/b/c"
	// pathlib.PathStr("a").Join("b/c") => "a/b/c"
	// pathlib.PathStr("a/b").Join("c") => "a/b/c"
	// pathlib.PathStr("a/b").Join("/c") => "a/b/c"
	// pathlib.PathStr("a/b").Join("../../../xyz") => "../xyz"
}

func ExamplePathStr_Rel() {
	example := func(a pathlib.PathStr, b pathlib.Dir) {
		val, err := a.Rel(b)
		if err == nil {
			fmt.Printf("%T(%q).Rel(%q) => %T(%q)\n", a, a, b, val, val)
		} else {
			fmt.Printf("%T(%q).Rel(%q) => Err(%v)\n", a, a, b, err)
		}
	}
	fmt.Println("On Unix")
	example("/a/b/c", "/a")
	example("/b/c", "/a")
	example("./b/c", "/a")
	// Output:
	// On Unix
	// pathlib.PathStr("/a/b/c").Rel("/a") => pathlib.PathStr("b/c")
	// pathlib.PathStr("/b/c").Rel("/a") => pathlib.PathStr("../b/c")
	// pathlib.PathStr("./b/c").Rel("/a") => Err(Rel: can't make ./b/c relative to /a)
}

func ExamplePathStr_IsAbsolute() {
	fmt.Println("On Unix:")
	fmt.Println(pathlib.PathStr("/home/gopher").IsAbsolute())
	fmt.Println(pathlib.PathStr(".bashrc").IsAbsolute())
	fmt.Println(pathlib.PathStr("..").IsAbsolute())
	fmt.Println(pathlib.PathStr(".").IsAbsolute())
	fmt.Println(pathlib.PathStr("/").IsAbsolute())
	fmt.Println(pathlib.PathStr("").IsAbsolute())

	// Output:
	// On Unix:
	// true
	// false
	// false
	// false
	// true
	// false
}

func ExamplePathStr_Ancestors_absolute() {
	for d := range pathlib.PathStr("/foo/bar/baz").Ancestors() {
		fmt.Println(d)
	}
	// output:
	// /foo/bar
	// /foo
	// /
}

func ExamplePathStr_Ancestors_relative() {
	for d := range pathlib.PathStr("./foo/bar/baz").Ancestors() {
		fmt.Println(d)
	}
	// this is the same as:
	for d := range pathlib.PathStr("foo/bar/baz").Ancestors() {
		fmt.Println(d)
	}
	// output:
	// foo/bar
	// foo
	// .
	// foo/bar
	// foo
	// .
}

func ExamplePathStr_Parts() {
	example := func(p pathlib.PathStr) {
		fmt.Printf("%q => %#v\n", p, p.Parts())
	}
	fmt.Println("On Unix:")
	example("/a/b")
	example("./a/b")
	example("a/b")
	example("a/../b")
	example("a//b")
	example("")
	// output:
	// On Unix:
	// "/a/b" => []string{"/", "a", "b"}
	// "./a/b" => []string{".", "a", "b"}
	// "a/b" => []string{"a", "b"}
	// "a/../b" => []string{"a", "..", "b"}
	// "a//b" => []string{"a", "b"}
	// "" => []string(nil)
}

func ExamplePathStr_Ext() {
	example := func(p pathlib.PathStr) {
		fmt.Printf("%q => %q\n", p, p.Ext())
	}
	example("index")
	example("index.js")
	example("main.test.js")
	// Output:
	// "index" => ""
	// "index.js" => ".js"
	// "main.test.js" => ".js"
}

func ExamplePathStr_BaseName() {
	example := func(p pathlib.PathStr) {
		fmt.Printf("%q => %q\n", p, p.BaseName())
	}
	fmt.Println("On Unix:")
	example("/foo/bar/baz.js")
	example("/foo/bar/baz")
	example("/foo/bar/baz/")
	example("dev.txt")
	example("../todo.txt")
	example("..")
	example(".")
	example("/")
	example("")
	// Output:
	// On Unix:
	// "/foo/bar/baz.js" => "baz.js"
	// "/foo/bar/baz" => "baz"
	// "/foo/bar/baz/" => "baz"
	// "dev.txt" => "dev.txt"
	// "../todo.txt" => "todo.txt"
	// ".." => ".."
	// "." => "."
	// "/" => "/"
	// "" => "."
}

func ExamplePathStr_ExpandUser() {
	home := expect(pathlib.UserHomeDir())
	example := func(p pathlib.PathStr) {
		expanded := expect(p.ExpandUser())
		fmt.Printf(
			"%q => %q\n",
			p,
			strings.Replace(string(expanded), string(home), "$HOME", 1),
		)
	}
	fmt.Println("On Unix:")
	example("~")
	example("~/foo/bar.txt")
	example("foo/~/bar")
	// Output:
	// On Unix:
	// "~" => "$HOME"
	// "~/foo/bar.txt" => "$HOME/foo/bar.txt"
	// "foo/~/bar" => "foo/~/bar"
}

func ExamplePathStr_beholder() {
	temp := expect(pathlib.TempDir().Join("path-str-beholder").AsDir().Make(0777))
	defer func() { _, _ = temp.RemoveAll() }()

	file := temp.Join("file.txt")
	expect(file.AsFile().Make(0644))

	rel := expect(file.Rel(temp))
	fmt.Printf("OnDisk: %q %s\n", rel, expect(file.OnDisk()).Mode())
	fmt.Printf(" Lstat: %q %s\n", rel, expect(file.Lstat()).Mode())
	fmt.Printf("  Stat: %q %s\n", rel, expect(file.Stat()).Mode())

	// Output:
	// OnDisk: "file.txt" -rw-r--r--
	//  Lstat: "file.txt" -rw-r--r--
	//   Stat: "file.txt" -rw-r--r--
}

func ExamplePathStr_read() {
	tmpDir := expect(pathlib.TempDir().Join("example-pathStr-read").AsDir().Make(0777))
	example := tmpDir.Join("example")

	{
		expect(example.AsDir().Make(0777))
		expect(example.Join("foo").AsFile().Make(0644))
		expect(example.Join("bar").AsDir().Make(0777))
		entries := expect(example.Read()).([]fs.DirEntry)
		for entry := range entries {
			fmt.Println(entry)
		}
		enforce(example.Remove())
	}
	{
		_, err := expect(example.AsFile().Make(0644)).WriteString("text")
		if err != nil {
			panic(err)
		}
		fmt.Println(expect(example.Read()))
		enforce(example.Remove())
	}
	{
		target := tmpDir.Join("target")
		if _, err := expect(target.AsFile().Make(0644)).WriteString("target"); err != nil {
			panic(err)
		}
		expect(example.AsSymlink().LinkTo(target))
		fmt.Println(expect(example.Read()))
	}
}

func ExamplePathStr_Read() {
	var dir = pathlib.TempDir().Join("PathStr.Read")
	var file = dir.Join("a.txt")
	var link = dir.Join("link")

	_, err := expect(file.AsFile().MakeAll(0666, 0777)).WriteString("file text")
	if err != nil {
		panic(err)
	}
	defer func() { expect(dir.AsDir().RemoveAll()) }()
	expect(link.AsSymlink().LinkTo(file))

	var contents any
	contents = expect(pathlib.PathStr(file).Read())
	fmt.Printf("file.Read() => %T(%q)\n", contents, contents)

	contents = expect(pathlib.PathStr(dir).Read())
	fmt.Printf("dir.Read() => %T(%q)\n", contents, contents)

	contents = expect(pathlib.PathStr(link).Read())
	fmt.Printf("link.Read() => %T(%q)", contents, contents)
	// Output:
	// file.Read() => []uint8("file text")
	// dir.Read() => []fs.DirEntry(["- a.txt" "L link"])
	// link.Read() => pathlib.PathStr("/tmp/PathStr.Read/a.txt")
}

func TestPathStr_Rename(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	a := dir.Join("a")
	expect(a.AsFile().Make(0666))
	b := expect(a.Rename(dir.Join("b")))
	if a.Exists() {
		t.Fail()
	}
	if !b.Exists() {
		t.Fail()
	}
}

func TestPathStr_chmod(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	f := dir.Join("a")
	expect(f.AsFile().Make(0666))
	testChmod(t, f, 0640)
}

func TestPathStr_chown(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	f := dir.Join("foo")
	expect(f.AsFile().Make(0666))
	testChown(t, pathlib.PathStr(f))
}

func ExamplePathStr_Remove() {
	dir := expect(pathlib.TempDir().Join("pathStr-remove").AsDir().Make(0777))
	defer func() { expect(dir.RemoveAll()) }()

	f := dir.Join("foo.txt")
	expect(f.AsFile().Make(0666))
	link := dir.Join("link")
	expect(link.AsSymlink().LinkTo(f))
	err := pathlib.PathStr(dir).Remove()
	if err != nil {
		fmt.Println("Unable to delete non-empty directory " + dir)
	}

	enforce(link.Remove())
	if !link.Exists() {
		fmt.Println("removed symlink " + link)
	}
	if f.Exists() {
		fmt.Println("removing " + link + " didn't affect link target " + f)
	} else {
		panic("removing " + link + " removed " + f)
	}

	enforce(f.Remove())
	if !f.Exists() {
		fmt.Println("removed file " + f)
	}
	// Output:
	// Unable to delete non-empty directory /tmp/pathStr-remove
	// removed symlink /tmp/pathStr-remove/link
	// removing /tmp/pathStr-remove/link didn't affect link target /tmp/pathStr-remove/foo.txt
	// removed file /tmp/pathStr-remove/foo.txt

}

func TestPathStr_localize(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	t.Chdir(temp.String())

	f := temp.Join("file.txt")
	expect(f.AsFile().Make(0666))

	f = expect(f.Rel(temp))
	localized := expect(f.Localize())
	if localized != f {
		t.Fatalf("expected %q, got %q", f, localized)
	}
	// Output:
	// "/tmp/pathlib_test/link" -> "/tmp/pathlib_test/file.txt"
}

func TestPathStr_read(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())

	f := temp.Join("file.txt")
	_, err := f.OnDisk()
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected fs.ErrNotExist, got %v", err)
	}
	_, err = f.Read()
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected fs.ErrNotExist, got %v", err)
	}
}

func TestPathStr_ExpandUser_homeless(t *testing.T) {
	unsetHome(t)
	if _, err := pathlib.PathStr("~/file.txt").ExpandUser(); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
