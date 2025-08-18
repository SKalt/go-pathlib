package pathlib_test

import (
	"strings"
	"testing"

	"github.com/skalt/pathlib.go"
)

func TestOnDisk_ensemble(t *testing.T) {
	dir := expect(pathlib.TempDir().
		Join("test/on-disk/ensemble.dir").
		AsDir().
		MakeAll(0777, 0777))

	assertEq := func(a, b any) {
		if a != b {
			t.Fatalf("- %#v\n+ %#v\n", a, b)
		}
	}
	info := expect(dir.OnDisk())
	assertEq(info.Parent(), pathlib.TempDir().Join("test/on-disk").AsDir())
	assertEq(info.BaseName(), "ensemble.dir")
	assertEq(info.Ext(), ".dir")
	assertEq(info.IsAbsolute(), true)
	assertEq(info.IsLocal(), false)
	assertEq(info.Join("foo"), pathlib.TempDir().Join("/test/on-disk/ensemble.dir/foo"))
	assertEq(
		strings.Join(info.Parts(), "\n"),
		strings.Join(
			[]string{
				"/",
				pathlib.TempDir().BaseName(),
				"test",
				"on-disk",
				"ensemble.dir",
			},
			"\n",
		),
	)
}

func TestOnDisk_Transformer(t *testing.T) {
	onDisk := expect(pathlib.Dir(t.TempDir()).OnDisk())

	if !onDisk.IsAbsolute() {
		t.Fail()
	}
	if expect(onDisk.Abs()) != onDisk.Path() {
		t.Fail()
	}
	if expect(onDisk.ExpandUser()) != onDisk.Path() {
		t.Fail()
	}
	if !onDisk.Eq(onDisk.Clean()) {
		t.Fatal(expect(onDisk.Abs()), expect(onDisk.Clean().Abs()))
	}
	if !strings.HasPrefix(onDisk.String(), "/") {
		t.Fail()
	}
	if expect(onDisk.Rel(onDisk.Parent())).String() != onDisk.BaseName() {
		t.Fail()
	}
	if _, err := onDisk.Localize(); err == nil { // localize fails when paths start with a /
		t.Fail()
	}
}

func TestOnDisk_manipulator(t *testing.T) {
	temp := pathlib.Dir(t.TempDir())
	dir := expect(temp.Join("foo").AsDir().Make(0777))
	onDisk := expect(dir.OnDisk())
	expect(onDisk.Chmod(0755))
	renamed := expect(onDisk.Rename(temp.Join("bar")))
	onDisk = expect(renamed.OnDisk())
	expect(onDisk.Remove())
}
