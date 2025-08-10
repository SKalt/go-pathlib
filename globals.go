package pathlib

import "os"

// Gets the Current Working Directory. See [os.Getwd].
func Cwd() Result[Dir] {
	dir, err := os.Getwd()
	return Result[Dir]{Dir(dir), err}
}

// See [os.UserHomeDir].
func UserHomeDir() Result[Dir] {
	dir, err := os.UserHomeDir()
	return Result[Dir]{Dir(dir), err}
}

// See [os.UserCacheDir].
func UserCacheDir() Result[Dir] {
	dir, err := os.UserCacheDir()
	return Result[Dir]{Dir(dir), err}
}

// See [os.UserConfigDir].
func UserConfigDir() Result[Dir] {
	dir, err := os.UserConfigDir()
	return Result[Dir]{Dir(dir), err}
}

// returns the process/os-wide temporary directory. See [os.TempDir].
func TempDir() Dir {
	return Dir(os.TempDir())
}
