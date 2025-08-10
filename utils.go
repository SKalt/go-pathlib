package pathlib

import (
	"io/fs"
	"os"
)

func isSymLink(m os.FileMode) bool {
	return (m & fs.ModeSymlink) == fs.ModeSymlink
}
