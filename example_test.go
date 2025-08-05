package pathlib_test

import (
	"fmt"

	"github.com/skalt/pathlib.go"
)

func Example() {
	dir := pathlib.TempDir().Join("pathlib-example").AsDir()
	defer dir.MustRemoveAll()
	if dir.Exists() {
		dir.MustRemoveAll()
	}

	onDisk := dir.MustMake(0o777).MustBeOnDisk()
	fmt.Printf("Created %s with mode %s\n", dir, onDisk.Mode())

	for i, subPath := range []string{"a.txt", "b.txt", "c/d.txt"} {
		file := dir.Join(subPath).AsFile()
		handle := file.MustMakeAll(0o666, 0o777)
		_, err := fmt.Fprintf(handle, "%d", i)
		if err != nil {
			panic(err)
		}
		if err = handle.Close(); err != nil {
			panic(err)
		}

		fmt.Printf("contents of %s: %q\n", file, string(file.MustRead()))
	}

	fmt.Printf("contents of %s:\n", dir)
	for _, entry := range dir.MustRead() {
		fmt.Println("  - " + entry.Name())
	}
	// Output:
	// Created /tmp/pathlib-example with mode drwxr-xr-x
	// contents of /tmp/pathlib-example/a.txt: "0"
	// contents of /tmp/pathlib-example/b.txt: "1"
	// contents of /tmp/pathlib-example/c/d.txt: "2"
	// contents of /tmp/pathlib-example:
	//   - a.txt
	//   - b.txt
	//   - c
}
