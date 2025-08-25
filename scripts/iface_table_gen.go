package main

import (
	"fmt"
	"iter"
	"os"
	"reflect"
	"strings"

	"github.com/skalt/pathlib.go"
)

func expect[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

type MethodResult struct {
	// types
	// interfaces
	methods map[string]int // note: this should be ordered
}
type InterfaceResult struct {
	types []assignability
}

func (InterfaceResult) interfaces() []string {
	return []string{
		"PurePath",
		"Transformer[P]",
		"Beholder[P]",
		"Maker[P]",
		"Changer[P]",
		"Mover[P]",
		"Destroyer[P]",
	}
}

func n[T any](val any) (name string) {
	var zero [0]T
	tt := reflect.TypeOf(zero).Elem()
	vt := reflect.TypeOf(val)
	if vt.AssignableTo(tt) {
		name = fmt.Sprintf("%v", tt)
	}
	return
}

func typeName(x any) (name string) {
	name = reflect.TypeOf(x).Name()
	if name == "" {
		name = n[pathlib.FileHandle](x)
		if name == "" {
			name = n[pathlib.Info[pathlib.PathStr]](x)
		}
	}
	if name == "" {
		panic(fmt.Sprintf("name for %T is %q", x, x))
	}
	name = strings.ReplaceAll(name, "pathlib.", "")
	name = strings.ReplaceAll(name, "github.com/skalt/pathlib%2ego.", "")

	return
}

type assignability struct {
	typeName   string
	interfaces map[string]bool
}

// func checks[T any]() {
// 	var zero [0]T
// }

func (r InterfaceResult) main() {
	for _, example := range types() {
		name := typeName(example)
		interfaces := r.interfaces()
		a := map[string]bool{}
		for _, i := range interfaces {
			a[i] = false
			// FIXME: loop over reflect interface values
			// _
		}
		r.types = append(r.types, assignability{name, a})
	}

	out := strings.Builder{}
	fmt.Println(out.String())
	// r.types = append(r.types, typeName())

}

var str pathlib.PathStr = "./README.md"
var dir pathlib.Dir = "."
var link pathlib.Symlink = "./link"
var file pathlib.File = str.AsFile()

func info() pathlib.Info[pathlib.PathStr] {
	return expect(str.OnDisk())
}
func handle() pathlib.FileHandle {
	return expect(str.AsFile().Open(os.O_RDONLY, 0666))
}

func types() []any {
	return []any{
		str,
		dir,
		link,
		file,
		info(),
		handle(),
	}
}

// [type][method]~[interface]

// these are pre-sorted alphabetically!
func methods(t reflect.Type) iter.Seq[reflect.Method] {
	return func(yield func(reflect.Method) bool) {
		n := t.NumMethod()
		for i := range n {
			if !yield(t.Method(i)) {
				return
			}
		}
	}
}

func args(t reflect.Type) iter.Seq[reflect.Type] {
	return func(yield func(reflect.Type) bool) {
		n := t.NumIn()
		for i := range n {
			if !yield(t.In(i)) {
				return
			}
		}
	}
}

func returnTypes(t reflect.Type) iter.Seq[reflect.Type] {
	return func(yield func(reflect.Type) bool) {
		n := t.NumOut()
		for i := range n {
			if !yield(t.Out(i)) {
				return
			}
		}
	}
}

type BitMap uint8

func assignableTo[Self pathlib.Kind]() (bitmap BitMap) {
	var s Self
	t := any(s)
	// name := typeName(s)
	i := 0
	if _, ok := t.(pathlib.PurePath); ok {
		bitmap |= 1 << i
	} else {
		bitmap |= 0 << i
	}
	i += 1
	if _, ok := t.(pathlib.Transformer[Self]); ok {
		bitmap |= 1 << i
	} else {
		bitmap |= 0 << i
	}
	i += 1
	if _, ok := t.(pathlib.Beholder[Self]); ok {
		bitmap |= 1 << i
	} else {
		bitmap |= 0 << i
	}
	i += 1
	if _, ok := t.(pathlib.Changer); ok {
		bitmap |= 1 << i
	} else {
		bitmap |= 0 << i
	}
	i += 1
	if _, ok := t.(pathlib.Remover[Self]); ok {
		bitmap |= 1 << i
	} else {
		bitmap |= 0 << i
	}
	i += 1
	if _, ok := t.(pathlib.Destroyer[Self]); ok {
		bitmap |= 1 << i
	} else {
		bitmap |= 0 << i
	}
	i += 1
	return
}

func main() {
	out := strings.Builder{}
	str := assignableTo[pathlib.PathStr]()
	dir := assignableTo[pathlib.Dir]()
	file := assignableTo[pathlib.File]()
	symlink := assignableTo[pathlib.Symlink]()
	table := []BitMap{str, dir, file, symlink}
	colHeaders := []string{
		"interface \\\\ `P`",
		"`PathStr`",
		"`Dir`",
		"`File`",
		"`Symlink`",
	}
	rowHeaders := []string{
		"PurePath",
		"Transformer[P]",
		"Beholder[P]",
		"Changer",
		"Reover[P]",
		"Destroyer[P]",
	}
	colWidths := make([]int, len(colHeaders))
	c0 := 0
	for _, rowHeader := range rowHeaders {
		c0 = max(c0, len(rowHeader))
	}
	colWidths[0] = c0 + 2
	for i, colHeader := range colHeaders[1:] {
		colWidths[i+1] = max(len("**false**"), len(colHeader)+2)
	}
	fmt.Println(colWidths)

	write := func(s string) {
		_, _ = out.WriteString(s)
	}
	for i, h := range colHeaders[:] {
		write("| ")
		write(h)
		write(" ")
		write(strings.Repeat(" ", max(0, colWidths[i]-len(h))))
	}
	write("|\n")
	for _, w := range colWidths {
		write("| ")
		write(strings.Repeat("-", w))
		write(" ")
	}
	write("|\n")

	for i, r := range rowHeaders {
		write("| ")
		if len(r) > 0 {
			write("`")
			write(r)
			write("`")
			write(strings.Repeat(" ", colWidths[0]-len(r)-2))
		} else {
			write(strings.Repeat(" ", colWidths[0]-2))
			write("  ")
		}
		write(" ")

		for j, rowVals := range table {
			write("| ")
			if rowVals&(1<<i) > 0 {
				write("true")
				write(strings.Repeat(" ", max(0, colWidths[j+1]-len("true"))))
			} else {
				write("**false**")
				write(strings.Repeat(" ", max(0, colWidths[j+1]-len("**false**"))))
			}
			write(" ")
		}
		write("|\n")
	}

	// // ifs := InterfaceResult{}
	// for _, example := range types() {
	// 	out.WriteString(typeName(example) + "\n")
	// 	t := reflect.TypeOf(example)
	// 	for m := range methods(t) {
	// 		out.WriteString("  ")
	// 		out.WriteString(m.Name)
	// 		out.WriteRune('(')
	// 		for arg := range args(m.Type) {
	// 			out.WriteString(arg.Name())
	// 			out.WriteRune(',')
	// 		}
	// 		out.WriteString( ") (")
	// 		for ret := range returnTypes(m.Type) {
	// 			out.WriteString(ret.Name())
	// 			out.WriteString(", ")
	// 		}
	// 		out.WriteString(")\n")
	// 	}
	// }
	fmt.Println(out.String())
	// 	t := reflect.TypeOf(example)
	// 	if t == nil {
	// 		panic("t is nil")
	// 	}
	// 	name := t.Name()
	// 	if name == "" { // t is an interface
	// 		name = t.Elem().Name()
	// 		//name = fmt.Sprintf("%T", example)
	// 	}

	// 	output.WriteString("name=" + name)
	// 	if _, ok := example.(pathlib.PurePath); ok {
	// 		output.WriteString(" PurePath")
	// 	}
	// 	// parametrized by self?
	// 	output.WriteString("\n")
	// }
	// fmt.Println(output.String())
	// str := pathlib.PathStr("")
	// t := reflect.TypeOf(str)

	// for m := range methods(t) {
	// 	fmt.Println(m.Name)
	// }
	// dir := pathlib.Dir(".")
	// link := pathlib.Symlink("link")
	// file := pathlib.File("example.txt")
	// handle := pathlib.Handle{}

	// using `reflect`, iterate over each type's methods

}
