package size

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

var (
	kernel32              = windows.NewLazySystemDLL("kernel32.dll")
	procGetDiskFreeSpaceW = kernel32.NewProc("GetDiskFreeSpaceW")
)

// GetSpaceOfDir .
func GetSpaceOfDir(dir string) (size int64, err error) {
	filepath.Walk(dir, func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil {
			err = errWalk
			return errWalk
		}
		if !info.IsDir() {
			size += info.Size()
		}

		return nil
	})

	return
}
