package space_test

import (
	"os"
	"testing"

	"github.com/ronbb/space/internal/space"
)

func TestDirectory(t *testing.T) {
	t.Run("TestDirectory", func(t *testing.T) {
		dir := "C:\\old"
		s, err := space.GetDirectorySpace(dir)
		if err != nil {
			t.Error(err)
			return
		}
		if s.Percentage > 1 || s.Percentage < 0 {
			t.Error("percentage should be 0 ~ 1 but got ", s.Percentage)
			return
		}

		t.Log(s.Directory, s.UsedSpace, s.Percentage)
	})
}

func TestVolume(t *testing.T) {
	t.Run("TestVolume", func(t *testing.T) {
		vol := "C:\\"
		s, err := space.GetVolumeSpace(vol)
		if err != nil {
			t.Error(err)
			return
		}
		if s.TotalSpace == 0 {
			t.Error("TotalSpace should be not 0 but got ", s.TotalSpace)
			return
		}

		t.Log(s.Volume, s.AvailableSpace, s.TotalSpace, s.FreeSpace)
	})
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
