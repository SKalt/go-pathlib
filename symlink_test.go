package pathlib_test

import (
	"fmt"
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
	_ = temp.Join("file.txt").AsFile().Make(0644).Unwrap()
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
	if symlink.OnDisk().Unwrap().Mode().IsRegular() {
		t.Fail()
	}

}
