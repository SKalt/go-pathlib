package pathlib_test

import (
	"strings"
	"testing"

	"github.com/skalt/pathlib.go"
)

func TestOnDisk_ensemble(t *testing.T) {
	dir := pathlib.TempDir().
		Join("test/on-disk/ensemble.dir").
		AsDir().
		MakeAll(0777, 0777).
		Unwrap()

	assertEq := func(a, b any) {
		if a != b {
			t.Fatalf("- %#v\n+ %#v\n", a, b)
		}
	}
	info := dir.OnDisk().Unwrap()
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
	onDisk := pathlib.Dir(t.TempDir()).OnDisk().Unwrap()

	if !onDisk.IsAbsolute() {
		t.Fail()
	}
	if onDisk.Abs().Unwrap() != onDisk.Path() {
		t.Fail()
	}
	if onDisk.ExpandUser().Unwrap() != onDisk.Path() {
		t.Fail()
	}
	if !onDisk.Eq(onDisk.Clean()) {
		t.Fatal(onDisk.Abs().Unwrap(), onDisk.Clean().Abs().Unwrap())
	}
	if !strings.HasPrefix(onDisk.String(), "/") {
		t.Fail()
	}
	if onDisk.Rel(onDisk.Parent()).Unwrap().String() != onDisk.BaseName() {
		t.Fail()
	}
	if onDisk.Localize().IsOk() { // localize fails when paths start with a /
		t.Fail()
	}
}

func TestOnDisk_manipulator(t *testing.T) {
	dir := pathlib.Dir(t.TempDir())
	onDisk := dir.Join(t.Name()).
		AsDir().
		Make(0777).
		Unwrap().OnDisk().Unwrap()
	onDisk.Chmod(0755).Unwrap()
	onDisk.Rename(dir.Join(t.Name() + "_foo")).Unwrap().
		OnDisk().Unwrap().
		Remove().Unwrap()
}
