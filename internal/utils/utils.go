package utils

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"strings"
)

// Hash .
func Hash(in string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(in))
	if err != nil {
		return "", err
	}
	out := h.Sum(nil)
	return fmt.Sprintf("%x", out), nil
}

// HashPath .
func HashPath(in string) (out string, err error) {
	p := in
	if p == filepath.VolumeName(p) {
		p += string(filepath.Separator)
	}
	p = filepath.Clean(p)
	if strings.HasSuffix(p, ".") {
		err = fmt.Errorf("%s is not a path", p)
	}
	p = strings.TrimRight(p, string(filepath.Separator))

	return Hash(p)
}
