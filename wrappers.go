package pathlib

import (
	"io/fs"
	"os"
	"path/filepath"
)

// adaptations stdlib packages, mostly private


func isAbsolute[P ~string](p P) bool {
	return filepath.IsAbs(string(p))
}

func isLocal[P ~string](p P) bool {
	return filepath.IsLocal(string(p))
}

func stat[P ~string](p P) (info os.FileInfo, err error) {
	info, err = os.Stat(string(p))
	return
}
func lstat[P ~string](p P) (info os.FileInfo, err error) {
	info, err = os.Lstat(string(p))
	return
}

func isSymLink(m os.FileMode) bool {
	return (m & fs.ModeSymlink) == fs.ModeSymlink
}
func isFifo(m os.FileMode) bool {
	return (m & fs.ModeNamedPipe) == fs.ModeNamedPipe
}
func isDevice(m os.FileMode) bool {
	return (m & fs.ModeDevice) == fs.ModeDevice
}
func isCharDevice(m os.FileMode) bool {
	return (m & fs.ModeCharDevice) == fs.ModeCharDevice
}
func isSocket(m os.FileMode) bool {
	return (m & fs.ModeSocket) == fs.ModeSocket
}
func isTemporary(m os.FileMode) bool {
	return (m & fs.ModeTemporary) == fs.ModeTemporary
}
