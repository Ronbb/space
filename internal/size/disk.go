package size

import "golang.org/x/sys/windows"

// DiskSpace .
type DiskSpace struct {
	freeBytesAvailableToCaller uint64
	TotalNumberOfBytes uint64
	TotalNumberOfFreeBytes uint64
}

// GetSpaceOfDisk .
func GetSpaceOfDisk(root string) (space DiskSpace, err error) {
	err = windows.GetDiskFreeSpaceEx(
		windows.StringToUTF16Ptr(root),
		&space.freeBytesAvailableToCaller,
		&space.TotalNumberOfBytes,
		&space.TotalNumberOfFreeBytes,
	)
	
	return
}