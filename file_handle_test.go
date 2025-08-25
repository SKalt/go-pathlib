package pathlib_test

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/skalt/pathlib.go"
)

func ExampleFileHandle_Read() {
	file := pathlib.TempDir().Join("example.txt").AsFile()
	handle := expect(file.Make(0o644))
	defer func() { _ = handle.Close() }()
	if _, err := handle.WriteString("Hello, World!"); err != nil {
		panic(err)
	}
	if err := handle.Close(); err != nil {
		panic(err)
	}
	handle = expect(file.Open(os.O_RDONLY, 0644))
	data := expect(io.ReadAll(handle))
	fmt.Printf("%s\n", string(data))
	// Output:
	// Hello, World!

}

func ExampleFileHandle_Path() {
	handle, err := pathlib.TempDir().Join("example.txt").AsFile().Make(0o644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = handle.Close() }()
	fmt.Println(handle.Path())
	// Output:
	// /tmp/example.txt
}

func ExampleFileHandle_Parts() {
	handle, err := pathlib.TempDir().Join("example.txt").AsFile().Make(0o644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = handle.Close() }()
	fmt.Println(handle.Parts())
	// Output:
	// [/ tmp example.txt]
}

func ExampleFileHandle_Join() {
	handle, err := pathlib.TempDir().Join("example.txt").AsFile().Make(0o644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = handle.Close() }()
	joined := handle.Join("../other.txt")
	fmt.Printf("%T(%s)", joined, joined)
	// Output:
	// pathlib.PathStr(/tmp/other.txt)
}
func tempDir(t *testing.T) pathlib.Dir {
	t.Helper()
	dir := pathlib.Dir(t.TempDir())
	return dir
}

func assertStrEq[S ~string](t *testing.T, expected, actual S) {
	if expected != actual {
		t.Errorf("expected %q\nactual  %q", actual, expected)
	}
}
func assertEq[S comparable](t *testing.T, expected, actual S) {
	if expected != actual {
		t.Errorf("expected %v\nactual  %v", actual, expected)
	}
}

func TestHandle_purePath(t *testing.T) {
	temp := tempDir(t)
	handle, err := temp.Join("example.txt").AsFile().Make(0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = handle.Close() }()
	var _ pathlib.PurePath = handle

	assertStrEq(t, handle.BaseName(), "example.txt")
	assertStrEq(t, temp, handle.Parent())
	assertEq(t, true, handle.IsAbsolute())
	assertEq(t, false, handle.IsLocal())
	assertEq(t, ".txt", handle.Ext())
}

func TestHandle_beholder(t *testing.T) {
	temp := tempDir(t)
	handle := expect(temp.Join("example.txt").AsFile().Make(0o644))
	defer func() { _ = handle.Close() }()
	var _ pathlib.Beholder[pathlib.File] = handle

	info := expect(handle.Stat())

	assertEq(t, true, info.Mode().IsRegular())
	assertEq(t, false, info.Mode().IsDir())
	assertStrEq(t, "example.txt", info.Name())

	info = expect(handle.OnDisk())
	assertEq(t, true, info.Mode().IsRegular())
	assertEq(t, false, info.Mode().IsDir())
	assertStrEq(t, "example.txt", info.Name())

	info = expect(handle.Lstat())
	assertEq(t, true, info.Mode().IsRegular())
	assertEq(t, false, info.Mode().IsDir())
	assertStrEq(t, "example.txt", info.Name())
}

func ExampleFileHandle_Remove() {
	handle, err := pathlib.TempDir().Join("example.txt").AsFile().Make(0o644)
	if err != nil {
		panic(err)
	}
	if err := handle.Remove(); err != nil {
		panic(err)
	}
	data := []byte{}
	_, err = handle.Read(data)
	if !errors.Is(err, os.ErrClosed) {
		panic("expected closed file")
	}
	fmt.Printf("%s exists: %t\n", handle, handle.Exists())
	// Output:
	// /tmp/example.txt exists: false
}

func TestHandle_nonexistent_stat(t *testing.T) {
	temp := tempDir(t)
	f := temp.Join("example.txt").AsFile()
	handle := expect(f.Make(0644))
	enforce(f.Remove())
	_, err := handle.Stat()
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected nonexistent file, got %v", err)
	}
	if _, err := (handle.Write([]byte("data"))); err == nil {
		t.Fatalf("expected write to fail on removed file")
	}
}

func TestHandle_nonexistent_stat_rename(t *testing.T) {
	temp := tempDir(t)
	f := temp.Join("example.txt").AsFile()
	handle := expect(f.Make(0644))
	_ = expect(handle.Rename(temp.Join("other.sh")))
	_, err := handle.Stat()
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected nonexistent file, got %v", err)
	}
	if _, err := (handle.Write([]byte("data"))); err == nil {
		t.Fatalf("expected write to fail on removed file")
	}
}

func TestFileHandle_chown(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt").AsFile()
	handle := expect(file.Make(0666))
	testChown(t, handle)
}

func TestFileHandle_chmod(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	file := temp.Join("file.txt").AsFile()
	handle := expect(file.Make(0666))
	enforce(handle.Chmod(0600))
	info := expect(handle.OnDisk())
	if info.Mode() != 0600 {
		t.Errorf("expected %s\ngot %ss", fs.FileMode(0600), info.Mode())
	}
}

func TestFileHandle_transformer(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	handle := expect(dir.Join("file.txt").AsFile().Make(0666))

	_, err := handle.Localize()
	if err == nil {
		t.Error("expected non-nil err for absolute path")
	}

	if !handle.IsAbsolute() {
		t.Fail()
	}
	if handle.IsLocal() {
		t.Fail()
	}
	if expect(handle.Abs()) != handle.Path() {
		t.Fail()
	}
	if expect(handle.ExpandUser()) != handle.Path() {
		t.Fail()
	}
	if !handle.Eq(handle.Clean()) {
		t.Fatal(expect(handle.Abs()), expect(handle.Clean().Abs()))
	}
	if !strings.HasPrefix(handle.String(), "/") {
		t.Fail()
	}
	if expect(handle.Rel(handle.Parent())).String() != handle.BaseName() {
		t.Fail()
	}
}

func TestFileHandle_openSymlink(t *testing.T) {
	temp := tempDir(t)
	h1 := expect(temp.Join("foo/bar/baz.txt").AsFile().MakeAll(0666, 0777))
	content := "example"
	expect(h1.WriteString(content))
	enforce(h1.Close())

	link := expect(temp.Join("link").AsSymlink().LinkTo("./foo/bar/baz.txt"))
	handle := expect(pathlib.PathStr(link).AsFile().Open(os.O_RDONLY, 0666))
	data := expect(io.ReadAll(handle))
	if string(data) != content {
		t.Fatalf("expected %q\ngot %q", content, string(data))
	}
	info := expect(handle.Stat())
	if info.Size() != int64(len(content)) {
		t.Fatalf("expected %dB, got %dB", len(content), info.Size())
	}
}
