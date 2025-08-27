package pathlib_test

import (
	"fmt"

	"github.com/skalt/pathlib.go"
)

// Syntactic sugar.
func enforce(err error) {
	if err != nil {
		panic(err)
	}
}

// More syntactic sugar.
func expect[T any](val T, err error) T {
	enforce(err)
	return val
}

func Example() {
	dir := expect(pathlib.TempDir().Join("pathlib-example").AsDir().Make(0o777))
	defer func() { expect(dir.RemoveAll()) }()

	onDisk := expect(dir.Stat())
	fmt.Printf("Created %s with mode %s\n", dir, onDisk.Mode())

	for i, subPath := range []string{"a.txt", "b.txt", "c/d.txt"} {
		file := dir.Join(subPath).AsFile()
		handle := expect(file.MakeAll(0o666, 0o777))
		expect(fmt.Fprintf(handle, "%d", i))
		enforce(handle.Close())
		fmt.Printf("contents of %s: %q\n", file, string(expect(file.Read())))
	}

	fmt.Printf("contents of %s:\n", dir)
	for _, entry := range expect(dir.Read()) {
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
