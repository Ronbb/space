package model

// DirectoryHash .
type DirectoryHash struct {
	Directory string `json:"directory"`
	Hash      string `json:"hash"`
}

// VolumeHash .
type VolumeHash struct {
	Volume string `json:"volume"`
	Hash   string `json:"hash"`
}
