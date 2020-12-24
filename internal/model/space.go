package model

// VolumeSpace .
type VolumeSpace struct {
	Time           int64 // unix second
	Volume         string
	AvailableSpace uint64
	TotalSpace     uint64
	FreeSpace      uint64
}

// DirectorySpace .
type DirectorySpace struct {
	Time       int64 // unix second
	Directory  string
	UsedSpace  uint64
	Percentage float64
}

// SpaceInfo .
type SpaceInfo struct {
	Time             int64 // unix second
	DirectorirsSpace []DirectorySpace
	VolumesSpace     []VolumeSpace
}
