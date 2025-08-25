package pathlib_test

import (
	"fmt"
	"path/filepath"
	"runtime"
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

func TestSymlink_linkChasing(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt")
	handle := expect(file.AsFile().Make(0644))
	content := "example"
	expect(handle.WriteString(content))
	enforce(handle.Close())
	link1 := expect(temp.Join("nested/dir").
		AsDir().
		MakeAll(0777, 0777)).
		Join("link1").AsSymlink()
	expect(link1.LinkTo("../../file.txt"))
	link2 := expect(temp.Join("link2").AsSymlink().LinkTo("nested/dir/link1"))

	data := expect(link2.Join().AsFile().Read())
	if string(data) != "example"{
		t.Fatal(string(data))
	}
	info := expect(link2.Stat())
	if info.Size() != int64(len(content)) {
		t.Fatalf("unexpected info: %#v", info)
	}
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

func TestSymlink_chown(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	file := dir.Join("file.txt")
	expect(file.AsFile().Make(0666))
	link := expect(dir.Join("link.to").AsSymlink().LinkTo(pathlib.PathStr(file)))

	testChown(t, link)
}

func TestSymlink_chmod(t *testing.T) {
	switch strings.Split(runtime.GOOS, "/")[0] {
	case "linux":
		t.Skip("On linux, symlinks cannot have permissions other than lrwxrwxrwx.")
		// see https://superuser.com/a/303063
	}

	dir := pathlib.Dir(t.TempDir())
	file := dir.Join("file.txt")
	expect(file.AsFile().Make(0666))
	link := expect(dir.Join("link.to").AsSymlink().LinkTo(pathlib.PathStr(file)))
	t.Fatalf("%o", expect(link.Lstat()).Mode())
	testChmod(t, link, 0600)
}

func TestSymlink_localize(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	link := expect(temp.Join("link").AsSymlink().Rel(temp))
	localized := expect(link.Localize())
	if localized != link {
		t.Fatalf("expected %q, got %q", link, localized)
	}
}

func TestSymlink_mover(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt")
	expect(file.AsFile().Make(0666))
	symlink := expect(temp.Join("link").AsSymlink().LinkTo(file))

	renamed := expect(symlink.Rename(temp.Join("renamed")))
	if symlink.Exists() {
		t.Fatal("renaming link should remove the original link")
	}
	if !renamed.Exists() {
		t.Fatal("renamed link should exist")
	}
	if !file.Exists() {
		t.Fatal("renaming symlink should not affect target file")
	}
	if expect(renamed.Read()) != file {
		t.Fatal("renamed link should point to the original target")
	}
	enforce(renamed.Remove())
	if renamed.Exists() {
		t.Fatal("removing link should remove the link")
	}
	if !file.Exists() {
		t.Fatal("removing link should not affect target file")
	}
}

// func TestSymlink_removeAll(t *testing.T) {
// 	temp := pathlib.Dir(t.TempDir())
// 	dir := expect(temp.Join("foo").AsDir().Make(0777))
// 	f := dir.Join("file.txt")
// 	expect(f.AsFile().Make(0666))
// 	link1 := expect(temp.Join("link1").AsSymlink().LinkTo(pathlib.PathStr(dir)))
// 	link2 := expect(temp.Join("link2").AsSymlink().LinkTo(pathlib.PathStr(link1)))

// 	expect(link2.RemoveAll())
// 	if !link1.Exists() {
// 		t.Fatal("removing link2 should not remove its target (link1)")
// 	}
// 	expect(link1.RemoveAll())
// 	if dir.Exists() {
// 		t.Fatal("removing link1 should not remove its target(dir)")
// 	}
// }
