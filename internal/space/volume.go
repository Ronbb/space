package space

import (
	"github.com/ronbb/space/internal/model"
	"golang.org/x/sys/windows"
)

// GetVolumeSpace .
func GetVolumeSpace(volume string) (space model.VolumeSpace, err error) {
	space.Volume = volume
	err = windows.GetDiskFreeSpaceEx(
		windows.StringToUTF16Ptr(volume),
		&space.AvailableSpace,
		&space.TotalSpace,
		&space.FreeSpace,
	)

	return
}
