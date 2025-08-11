# `pathlib.go`

String types for path manipulation.

## Benefits of using string types to represent paths

  1. you can treat paths exactly like strings
  2. String subtypes prevent common mix-ups: the go compiler will warn you about passing a string that could represent a file to a function that accepts a `pathlib.Dir` argument.
  3. Each string subtype can have associated methods, consolidating functionality from `"path/filepath"`, `"os"`, and `"io/fs"`


## Prior art

Inspired by Python's [`pathlib`][python-pathlib] and Go's filesystem methods scattered throughout the <abbr title="Standard library">stdlib</abbr>.

There are plenty of [other][other] ports of `pathlib` to go, but only one other package [one other][forerunner] uses string types to represent paths.

## Trade-offs

Using this library adds some overhead compared to using the stdlib functions directly.
Rough measurements on my x86_64 machine indicated that linking `pathlib` adds around a kilobyte to the size of a binary that does nothing except links `path/filepath`.

This library also returns a `Result[T]` type instead of a more-idiomatic `(*T, error)` tuple.
While this makes tab-completing scripts easier, it also side-steps `errcheck` lints when neither method is called.
Not calling those methods is effectively equivalent to `_, _ = someFallibleOperation()`.
Since the result's fields are only accessible through the `Result.Unwrap() T` or `Result.Unpack() (T, error)` methods, the result type doesn't introduce any particularly new ways to shoot yourself in the foot.



[python-pathlib]: https://docs.python.org/3/library/pathlib.html
[other]: https://pkg.go.dev/search?q=pathlib
[forerunner]: https://pkg.go.dev/github.com/gershwinlabs/pathlib
