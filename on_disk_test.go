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
		MakeAll(0755, 0755).
		Unwrap()
	defer func() { dir.RemoveAll().Unwrap() }()

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
