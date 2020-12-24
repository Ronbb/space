package utils_test

import (
	"testing"

	"github.com/ronbb/space/internal/utils"
)

func TestHashPath(t *testing.T) {
	t.Run("TestHashDirectory", func(t *testing.T) {
		h1, err := utils.HashPath("C:\\Project\\")
		if err != nil {
			t.Error(err)
			return
		}
		h2, err := utils.HashPath("C:\\Project")
		if err != nil {
			t.Error(err)
			return
		}
		h3, err := utils.HashPath("C:/Project/")
		if err != nil {
			t.Error(err)
			return
		}
		h4, err := utils.HashPath("C:/Project")
		if err != nil {
			t.Error(err)
			return
		}
		h5, err := utils.HashPath("C:/Project///")
		if err != nil {
			t.Error(err)
			return
		}
		hs := []string{
			h2, h3, h4, h5,
		}
		for _, h := range hs {
			if h1 != h {
				t.Error("hashes are not same")
				return
			}
		}
	})

	t.Run("TestHashVolume", func(t *testing.T) {
		h1, err := utils.HashPath("C:\\")
		if err != nil {
			t.Error(err)
			return
		}
		h2, err := utils.HashPath("C:")
		if err != nil {
			t.Error(err)
			return
		}
		h3, err := utils.HashPath("C:/")
		if err != nil {
			t.Error(err)
			return
		}
		h4, err := utils.HashPath("C://")
		if err != nil {
			t.Error(err)
			return
		}
		hs := []string{
			h2, h3, h4,
		}
		for _, h := range hs {
			if h1 != h {
				t.Error("hashes are not same")
				return
			}
		}
	})
}
