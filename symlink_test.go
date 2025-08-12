package pathlib_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExampleSymlink_Lstat() {
	tempDir := pathlib.Cwd().
		Unwrap().
		Join("temp", "symlink-lstat").
		AsDir().
		MakeAll(0o755, 0o755).
		Unwrap()
	defer tempDir.RemoveAll()

	file := tempDir.Join("file.txt").AsFile()
	file.Make(0o666).Unwrap()

	{
		link := tempDir.Join("link")
		if link.String() != tempDir.String()+"/link" {
			panic(link)
		}
	}
	link := tempDir.Join("link").
		AsSymlink().
		LinkTo(pathlib.PathStr(file.String())).
		Unwrap()

	onDisk := link.Lstat().Unwrap()
	fmt.Printf(
		"%s -> %s",
		onDisk.Path().Rel(tempDir).Unwrap(),
		link.Read().Unwrap().Rel(tempDir).Unwrap(),
	)
	// Output:
	// link -> file.txt
}

func TestSymlink_beholder(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt")
	file.AsFile().Make(0644).Unwrap()
	symlink := temp.Join("nested/dir").
		AsDir().
		MakeAll(0777, 0777).Unwrap().
		Join("link").AsSymlink()

	symlink.LinkTo("../../file.txt").Unwrap()
	stat := symlink.Stat().Unwrap()
	if !stat.Mode().IsRegular() || stat.Mode().Perm() != 0644 {
		t.Fatalf("stat: %o", stat.Mode())
	}
	lstat := symlink.Lstat().Unwrap()
	if lstat.Mode().IsRegular() {
		t.Fail()
	}
	if !symlink.Exists() {
		t.Fail()
	}
	if symlink.OnDisk().Unwrap().Mode().IsRegular() {
		t.Fail()
	}
	_, err := file.AsSymlink().Lstat().Unpack()
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
	if err := file.Make(0666).Unwrap().Close(); err != nil {
		panic(err)
	}
	link := dir.Join("link.to").AsSymlink().LinkTo(pathlib.PathStr(file)).Unwrap()

	_, err := link.Join("foo").Lstat().Unpack()
	if err == nil {
		t.Fail()
	}

	if !link.IsAbsolute() {
		t.Fail()
	}
	if link.IsLocal() {
		t.Fail()
	}
	if link.Abs().Unwrap() != link {
		t.Fail()
	}
	if link.ExpandUser().Unwrap() != link {
		t.Fail()
	}
	if !link.Eq(link.Clean()) {
		t.Fatal(link.Abs().Unwrap(), link.Clean().Abs().Unwrap())
	}
	if !strings.HasPrefix(link.String(), "/") {
		t.Fail()
	}
	if link.Rel(link.Parent()).Unwrap().String() != link.BaseName() {
		t.Fail()
	}
}

func TestSymlink_purePath(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	file := dir.Join("file.txt").AsFile()
	if err := file.Make(0666).Unwrap().Close(); err != nil {
		panic(err)
	}
	link := dir.Join("link.to").AsSymlink().LinkTo(pathlib.PathStr(file)).Unwrap()
	ext := link.Ext()
	if ext != ".to" {
		t.Fail()
	}
	if filepath.Ext(link.String()) != ext {
		t.Fail()
	}
	_, err := link.Join("foo").Lstat().Unpack()
	if err == nil {
		t.Fail()
	}

	if !link.IsAbsolute() {
		t.Fail()
	}
	if link.IsLocal() {
		t.Fail()
	}
	if len(link.Rel(link.Parent()).Unwrap().Parts()) != 1 {
		t.Fail()
	}
}
