package pathlib

import "os"

func Cwd() (Dir, error) {
	dir, err := os.Getwd()
	return Dir(dir), err
}

func UserHomeDir() (Dir, error) {
	dir, err := os.UserHomeDir()
	return Dir(dir), err
}

func UserCacheDir() (Dir, error) {
	dir, err := os.UserCacheDir()
	return Dir(dir), err
}

func UserConfigDir() (Dir, error) {
	dir, err := os.UserConfigDir()
	return Dir(dir), err
}

// returns the process/os-wide temporary directory
func TempDir() Dir {
	return Dir(os.TempDir())
}
