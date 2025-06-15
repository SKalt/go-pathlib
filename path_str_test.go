package pathlib_test

import (
	"fmt"

	"github.com/skalt/go-pathlib"
)

func ExamplePathStr_IsLocal() {
	p := pathlib.PathStr("/absolute/path/to/file.txt")
	fmt.Println(p.IsLocal()) // false

	p = pathlib.PathStr("relative/path/to/file.txt")
	fmt.Println(p.IsLocal()) // true

	p = pathlib.PathStr("./local/path/to/file.txt")
	fmt.Println(p.IsLocal()) // true
	// Output:
	// false
	// true
	// true
}
func ExamplePathStr_Join() {
	fmt.Println("On Unix:")
	fmt.Println(pathlib.PathStr("a").Join("b", "c"))
	fmt.Println(pathlib.PathStr("a").Join("b/c"))
	fmt.Println(pathlib.PathStr("a/b").Join("c"))
	fmt.Println(pathlib.PathStr("a/b").Join("/c"))

	fmt.Println(pathlib.PathStr("a/b").Join("../../../xyz"))
	// Output:
	// On Unix:
	// a/b/c
	// a/b/c
	// a/b/c
	// a/b/c
	// ../xyz
}

func ExamplePathStr_IsAbsolute() {
	fmt.Println("On Unix:")
	fmt.Println(pathlib.PathStr("/home/gopher").IsAbsolute())
	fmt.Println(pathlib.PathStr(".bashrc").IsAbsolute())
	fmt.Println(pathlib.PathStr("..").IsAbsolute())
	fmt.Println(pathlib.PathStr(".").IsAbsolute())
	fmt.Println(pathlib.PathStr("/").IsAbsolute())
	fmt.Println(pathlib.PathStr("").IsAbsolute())

	// Output:
	// On Unix:
	// true
	// false
	// false
	// false
	// true
	// false
}

func ExamplePathStr_Ancestors_absolute() {
	for d := range pathlib.PathStr("/foo/bar/baz").Ancestors() {
		fmt.Println(d)
	}
	// output:
	// /foo/bar
	// /foo
	// /
}

func ExamplePathStr_Ancestors_relative() {
	for d := range pathlib.PathStr("./foo/bar/baz").Ancestors() {
		fmt.Println(d)
	}
	// this is the same as:
	for d := range pathlib.PathStr("foo/bar/baz").Ancestors() {
		fmt.Println(d)
	}
	// output:
	// foo/bar
	// foo
	// .
	// foo/bar
	// foo
	// .
}

// func ExamplePathStr_Ancestors_unc() {
// 	// an UNC path:
// 	for d := range pathlib.PathStr("//foo/bar/baz").Ancestors() {
// 		fmt.Println(d)
// 	}
// 	// output:
// 	// /foo/bar
// 	// /foo
// 	// /
// }
