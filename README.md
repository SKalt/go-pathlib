# `pathlib.go`

String types for path manipulation.


<!-- TODO: usage section -->

## Prior art

Inspired by Python's ergonomic [`pathlib`][python-pathlib] and Go's cornucopia of useful filesystem methods scattered throughout the <abbr title="Standard library">stdlib</abbr>.

There are plenty of [other][other] ports of `pathlib` to go, but this one does something only [one other][forerunner] did: used a string subtype to represent paths. This has several benefits:
  1. you can treat paths exactly like strings
      ```go
      const p pathlib.PathStr = "/works"
      var   q pathlib.PathStr = "./works.too"
      concat := p + q
      fmt.Printf("%T", concat) // => pathlib.PathStr
      ```
  2. String subtypes prevent common mix-ups: the go compiler will warn you about passing a string that could represent a file to a function that accepts a `pathlib.Dir` argument.
  3. Each string subtype can have associated methods, consolidating functionality from `"path/filepath"`, `"os"`, and `"io/fs"`

## Trade-offs: syntactic sugar for performance

I haven't checked how using this library compares to using the stdlib functions. I expect it adds some overhead to binary size and runtime.


[python-pathlib]: https://docs.python.org/3/library/pathlib.html
[other]: https://pkg.go.dev/search?q=pathlib
[forerunner]: https://pkg.go.dev/github.com/gershwinlabs/pathlib
