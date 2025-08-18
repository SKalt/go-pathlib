package pathlib_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExampleSymlink_Lstat() {
	tempDir := expect(pathlib.Cwd()).
		Join("temp", "symlink-lstat").
		AsDir()
	expect(tempDir.MakeAll(0o755, 0o755))
	defer func() { expect(tempDir.RemoveAll()) }()

	file := tempDir.Join("file.txt").AsFile()
	expect(file.Make(0o666))

	{
		link := tempDir.Join("link")
		if link.String() != tempDir.String()+"/link" {
			panic(link)
		}
	}
	link := expect(tempDir.Join("link").
		AsSymlink().
		LinkTo(pathlib.PathStr(file.String())))

	onDisk := expect(link.Lstat())

	fmt.Printf(
		"%s -> %s",
		expect(onDisk.Path().Rel(tempDir)),
		expect(expect(link.Read()).Rel(tempDir)),
	)
	// Output:
	// link -> file.txt
}

func TestSymlink_beholder(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt")
	expect(file.AsFile().Make(0644))
	symlink := expect(temp.Join("nested/dir").
		AsDir().
		MakeAll(0777, 0777)).
		Join("link").AsSymlink()

	expect(symlink.LinkTo("../../file.txt"))
	stat := expect(symlink.Stat())
	if !stat.Mode().IsRegular() || stat.Mode().Perm() != 0644 {
		t.Fatalf("stat: %o", stat.Mode())
	}
	lstat := expect(symlink.Lstat())
	if lstat.Mode().IsRegular() {
		t.Fail()
	}
	if !symlink.Exists() {
		t.Fail()
	}
	if expect(symlink.OnDisk()).Mode().IsRegular() {
		t.Fail()
	}
	_, err := file.AsSymlink().Lstat()
	if err == nil {
		t.Fail()
	}
	// if _, ok := err.(pathlib.WrongTypeOnDisk[pathlib.Symlink]); !ok {
	// 	t.Fail()
	// }
}

func TestSymlink_transformer(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	file := dir.Join("file.txt").AsFile()
	if err := expect(file.Make(0666)).Close(); err != nil {
		panic(err)
	}
	link := expect(dir.Join("link.to").AsSymlink().LinkTo(pathlib.PathStr(file)))

	_, err := link.Join("foo").Lstat()
	if err == nil {
		t.Fail()
	}

	if !link.IsAbsolute() {
		t.Fail()
	}
	if link.IsLocal() {
		t.Fail()
	}
	if expect(link.Abs()) != link {
		t.Fail()
	}
	if expect(link.ExpandUser()) != link {
		t.Fail()
	}
	if !link.Eq(link.Clean()) {
		t.Fatal(expect(link.Abs()), expect(link.Clean().Abs()))
	}
	if !strings.HasPrefix(link.String(), "/") {
		t.Fail()
	}
	if expect(link.Rel(link.Parent())).String() != link.BaseName() {
		t.Fail()
	}
}

func TestSymlink_purePath(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	file := dir.Join("file.txt").AsFile()
	if err := expect(file.Make(0666)).Close(); err != nil {
		panic(err)
	}
	link := expect(dir.Join("link.to").AsSymlink().LinkTo(pathlib.PathStr(file)))
	ext := link.Ext()
	if ext != ".to" {
		t.Fail()
	}
	if filepath.Ext(link.String()) != ext {
		t.Fail()
	}
	_, err := link.Join("foo").Lstat()
	if err == nil {
		t.Fail()
	}

	if !link.IsAbsolute() {
		t.Fail()
	}
	if link.IsLocal() {
		t.Fail()
	}
	if len(expect(link.Rel(link.Parent())).Parts()) != 1 {
		t.Fail()
	}
}
