package pathlib

import "os"

// Gets the Current Working Directory. See [os.Getwd].
func Cwd() (Dir, error) {
	dir, err := os.Getwd()
	return Dir(dir), err
}

// See [os.UserHomeDir].
func UserHomeDir() (Dir, error) {
	dir, err := os.UserHomeDir()
	return Dir(dir), err
}

// See [os.UserCacheDir].
func UserCacheDir() (Dir, error) {
	dir, err := os.UserCacheDir()
	return Dir(dir), err
}

// See [os.UserConfigDir].
func UserConfigDir() (Dir, error) {
	dir, err := os.UserConfigDir()
	return Dir(dir), err
}

// returns the process/os-wide temporary directory. See [os.TempDir].
func TempDir() Dir {
	return Dir(os.TempDir())
}
