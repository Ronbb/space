package model

// VolumeSpace .
type VolumeSpace struct {
	Time           int64  `json:"time"` // unix second
	Volume         string `json:"volume"`
	AvailableSpace uint64 `json:"availableSpace"`
	TotalSpace     uint64 `json:"totalSpace"`
	FreeSpace      uint64 `json:"freeSpace"`
	Limit          int64  `json:"limit"`
}

// DirectorySpace .
type DirectorySpace struct {
	Time       int64   `json:"time"` // unix second
	Directory  string  `json:"directory"`
	UsedSpace  uint64  `json:"usedSpace"`
	Percentage float64 `json:"percentage"`
	Limit      int64   `json:"limit"`
}

// SpaceRecord .
type SpaceRecord struct {
	Time             int64            `json:"time"` // unix second
	DirectoriesSpace []DirectorySpace `json:"directoriesSpace"`
	VolumesSpace     []VolumeSpace    `json:"volumesSpace"`
}
