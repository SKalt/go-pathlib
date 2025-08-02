package pathlib_test

import (
	"fmt"
	"strings"

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

func ExamplePathStr_Parts() {
	example := func(p pathlib.PathStr) {
		fmt.Printf("%q => %#v\n", p, p.Parts())
	}
	fmt.Println("On Unix:")
	example("/a/b")
	example("./a/b")
	example("a/b")
	example("a/../b")
	example("a//b")
	// output:
	// On Unix:
	// "/a/b" => []string{"/", "a", "b"}
	// "./a/b" => []string{".", "a", "b"}
	// "a/b" => []string{"a", "b"}
	// "a/../b" => []string{"a", "..", "b"}
	// "a//b" => []string{"a", "b"}
}

func ExamplePathStr_Ext() {
	example := func(p pathlib.PathStr) {
		fmt.Printf("%q => %q\n", p, p.Ext())
	}
	example("index")
	example("index.js")
	example("main.test.js")
	// Output:
	// "index" => ""
	// "index.js" => ".js"
	// "main.test.js" => ".js"
}

func ExamplePathStr_BaseName() {
	example := func(p pathlib.PathStr) {
		fmt.Printf("%q => %q\n", p, p.BaseName())
	}
	fmt.Println("On Unix:")
	example("/foo/bar/baz.js")
	example("/foo/bar/baz")
	example("/foo/bar/baz/")
	example("dev.txt")
	example("../todo.txt")
	example("..")
	example(".")
	example("/")
	example("")
	// Output:
	// On Unix:
	// "/foo/bar/baz.js" => "baz.js"
	// "/foo/bar/baz" => "baz"
	// "/foo/bar/baz/" => "baz"
	// "dev.txt" => "dev.txt"
	// "../todo.txt" => "todo.txt"
	// ".." => ".."
	// "." => "."
	// "/" => "/"
	// "" => "."
}

func ExamplePathStr_ExpandUser() {
	home, _ := pathlib.UserHomeDir()
	example := func(p pathlib.PathStr) {
		expanded, _ := p.ExpandUser()
		fmt.Printf(
			"%q => %q\n",
			p,
			strings.Replace(string(expanded), string(home), "$HOME", 1),
		)
	}
	fmt.Println("On Unix:")
	example("~")
	example("~/foo/bar.txt")
	example("foo/~/bar")
	// Output:
	// On Unix:
	// "~" => "$HOME"
	// "~/foo/bar.txt" => "$HOME/foo/bar.txt"
	// "foo/~/bar" => "foo/~/bar"
}
