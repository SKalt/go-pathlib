package pathlib

import (
	"io/fs"
	"os"
)

// adaptations stdlib packages, mostly private

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
