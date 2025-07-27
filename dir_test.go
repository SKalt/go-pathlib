package pathlib_test

import (
	"fmt"

	"github.com/skalt/go-pathlib"
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
	demoDir := pathlib.TempDir().Join("glob-example").AsDir().MustMake()
	defer demoDir.RemoveAll()

	for _, name := range []string{"x", "y", "z", "a", "b", "c"} {
		f := demoDir.Join(name + ".txt").AsFile()
		if err := f.Touch(); err != nil {
			panic(err)
		}
	}

	for _, match := range demoDir.MustGlob("*.txt") {
		fmt.Println(match)
		// fmt.Println(strings.Replace(string(match), string(demoDir), "$DEMO", 1))
	}

	// output:
	// /tmp/glob-example/a.txt
	// /tmp/glob-example/b.txt
	// /tmp/glob-example/c.txt
	// /tmp/glob-example/x.txt
	// /tmp/glob-example/y.txt
	// /tmp/glob-example/z.txt
}
