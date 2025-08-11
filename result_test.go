package pathlib_test

import (
	"fmt"

	"github.com/skalt/pathlib.go"
)

func ExampleResult_Unwrap() {
	var failure = pathlib.PathStr("does/not/exist").OnDisk()
	fmt.Println(".IsOk():", failure.IsOk())
	val, err := failure.Unpack()
	fmt.Printf(".Unpack(): %#v %#v\n", val, err)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(".Unwrap() panicked")
		}
	}()
	val = failure.Unwrap()
	fmt.Printf(".Unwrap(): %#v", val)
	// Output:
	// .IsOk(): false
	// .Unpack(): <nil> &fs.PathError{Op:"lstat", Path:"does/not/exist", Err:0x2}
	// .Unwrap() panicked
}
