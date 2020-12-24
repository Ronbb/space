package model

// VolumeSpace .
type VolumeSpace struct {
	Time           int64  `json:"time"` // unix second
	Volume         string `json:"volume"`
	AvailableSpace uint64 `json:"availableSpace"`
	TotalSpace     uint64 `json:"totalSpace"`
	FreeSpace      uint64 `json:"freeSpace"`
}

// DirectorySpace .
type DirectorySpace struct {
	Time       int64   `json:"time"` // unix second
	Directory  string  `json:"directory"`
	UsedSpace  uint64  `json:"usedSpace"`
	Percentage float64 `json:"percentage"`
}

// SpaceInfo .
type SpaceInfo struct {
	Time             int64            `json:"time"` // unix second
	DirectoriesSpace []DirectorySpace `json:"directoriesSpace"`
	VolumesSpace     []VolumeSpace    `json:"volumesSpace"`
}
