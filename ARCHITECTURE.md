# Architecture

## Goals
- move as much work to the type system as possible without allocating new memory on the heap
- keep as much as possible in string subtypes


## Non-goals
- validating paths. Given different OS path-naming requirements, validation would be hard to do (see https://stackoverflow.com/a/31976060/6571327), and


## Ontology

- `PathStr` represents any kind of path that may/not exist
- it can be cast to several specialized subtypes without checks
  - `Dir` (for globbing, listing subdirectories)
  - `Symlink`
  - `Fifo`
  - `Device`
- Any subtype can be converted to an `OnDisk[T]`, which is a wrapper around a `fs.FileInfo` pointer.
- `OnDisk[T]` grants access to file-manipulation methods, such as `Chmod()`, `Chown()`. The benefit of knowing that the file existed on-disk is obviated by possible concurrent writes to disk: a file can be changed/deleted once it has been `stat`ed.
