package space

import (
	"os"
	"path/filepath"

	"github.com/ronbb/space/internal/model"
)

// GetDirectorySpace .
func GetDirectorySpace(dir string) (space model.DirectorySpace, err error) {
	used := uint64(0)
	filepath.Walk(dir, func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil {
			err = errWalk
			return errWalk
		}
		if !info.IsDir() {
			used += uint64(info.Size())
		}

		return nil
	})

	volumeSpace, err := GetVolumeSpace(filepath.VolumeName(dir))
	if err != nil {
		return
	}

	space.Directory = dir
	space.UsedSpace = used
	space.Percentage = float64(used) / float64(volumeSpace.TotalSpace)

	return
}
