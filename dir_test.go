package pathlib_test

import (
	"fmt"

	"github.com/skalt/pathlib.go"
)

func ExampleDir_Eq() {
	cwd, err := pathlib.Cwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(cwd.Eq("."))
	fmt.Println(pathlib.Dir("/foo").Eq("/foo"))
	fmt.Println(cwd.Join("./relative").Eq("relative"))
	fmt.Println(pathlib.Dir("/a/b").Eq("b"))
	// output:
	// true
	// true
	// true
	// false
}

func ExampleDir_Join() {
	result := pathlib.Dir("/tmp").Join("a/b")
	fmt.Printf("%s is a %T", result, result)
	// output:
	// /tmp/a/b is a pathlib.PathStr
}

func ExampleDir_Glob() {
	demoDir := pathlib.
		TempDir().
		Join("glob-example").
		AsDir().
		MustMakeAll(0o777, 0o777)

	defer demoDir.MustRemoveAll()

	for _, name := range []string{"x", "y", "z", "a", "b", "c"} {
		demoDir.Join(name + ".txt").AsFile().MustMake(0o777)
	}

	for _, match := range demoDir.MustGlob("*.txt") {
		fmt.Println(match)
	}

	// output:
	// /tmp/glob-example/a.txt
	// /tmp/glob-example/b.txt
	// /tmp/glob-example/c.txt
	// /tmp/glob-example/x.txt
	// /tmp/glob-example/y.txt
	// /tmp/glob-example/z.txt
}
