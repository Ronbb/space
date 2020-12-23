package database

import (
	"crypto/sha1"
	"fmt"
)

func hash(in string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(in))

	if err != nil {
		return "", err
	}
	out := h.Sum(nil)
	return fmt.Sprintf("%x", out), nil
}
