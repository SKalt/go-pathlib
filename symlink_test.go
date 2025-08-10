package pathlib_test

import (
	"fmt"

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
