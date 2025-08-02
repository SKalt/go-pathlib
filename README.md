# `pathlib.go`

String types for path manipulation.

## Benefits of using string types to represent paths

  1. you can treat paths exactly like strings
  2. String subtypes prevent common mix-ups: the go compiler will warn you about passing a string that could represent a file to a function that accepts a `pathlib.Dir` argument.
  3. Each string subtype can have associated methods, consolidating functionality from `"path/filepath"`, `"os"`, and `"io/fs"`


## Prior art

Inspired by Python's ergonomic [`pathlib`][python-pathlib] and Go useful filesystem methods scattered throughout the <abbr title="Standard library">stdlib</abbr>.

There are plenty of [other][other] ports of `pathlib` to go, but only one other package [one other][forerunner] uses string types to represent paths.

## Trade-offs

Using this library adds some overhead compared to using the stdlib functions directly.
Rough measurements on my x86_64 machine indicated that linking `pathlib` adds around a kilobyte to the size of a binary that does nothing except links `path/filepath`.


[python-pathlib]: https://docs.python.org/3/library/pathlib.html
[other]: https://pkg.go.dev/search?q=pathlib
[forerunner]: https://pkg.go.dev/github.com/gershwinlabs/pathlib
