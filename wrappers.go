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

func stat[P ~string](p P) (r result[os.FileInfo]) {
	*r.value, r.err = os.Stat(string(p))
	if r.err != nil {
		r.value = nil // for consistency
	}
	return
}
func lstat[P ~string](p P) (r result[os.FileInfo]) {
	*r.value, r.err = os.Lstat(string(p))
	if r.err != nil {
		r.value = nil // for consistency
	}
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
