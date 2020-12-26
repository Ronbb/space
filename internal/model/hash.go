package model

// DirectoryHash .
type DirectoryHash struct {
	Directory string `json:"directory"`
	Hash      string `json:"hash"`
	Limit     int64  `json:"limit"` // used space should less than limit
}

// VolumeHash .
type VolumeHash struct {
	Volume string `json:"volume"`
	Hash   string `json:"hash"`
	Limit  int64  `json:"limit"` // used space should less than limit
}
