package model

// Node .
type Node struct {
	ReferenceNumber       string `json:"referenceNumber"`
	ParentReferenceNumber string `json:"parentReferenceNumber"`
	FileName              string `json:"fileName"`
	Size                  uint64 `json:"size"`
}
