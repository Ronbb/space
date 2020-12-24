package database_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ronbb/space/internal/database"
	"github.com/ronbb/space/internal/model"
)

var (
	db   database.DB
	dirs []string
	vols []string

	randNum = rand.Uint32() / 2
	randAlp = rand.Intn(12)
)

func init() {
	dirs = []string{
		randDir(),
		randDir(),
		randDir(),
		randDir(),
		randDir(),
	}

	vols = []string{
		randVol(),
		randVol(),
		randVol(),
		randVol(),
		randVol(),
	}
}

func randDir() string {
	randNum += 1
	return fmt.Sprintf("C:\\%d\\", randNum)
}

func randVol() string {
	randAlp += 1
	return fmt.Sprintf("%s:\\", string(rune(int('A')+randAlp)))
}

func TestDirectory(t *testing.T) {
	t.Run("Put Directory", func(t *testing.T) {
		for _, dir := range dirs {
			err := db.PutDirectory(dir)
			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Get Directory", func(t *testing.T) {
		dhs, err := db.GetDirectories()
		if err != nil {
			t.Error(err)
		}
		if len(dhs) != len(dirs) {
			t.Error("len(dhs) != len(dirs)")
		}
		for _, dir := range dirs {
			found := false
			for _, dh := range dhs {
				if dh.Directory == dir {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("can not find %s", dir)
			}
		}
	})

	t.Run("Remove Directory", func(t *testing.T) {
		dir := dirs[rand.Intn(len(dirs))]
		err := db.RemoveDirectory(dir)
		if err != nil {
			t.Error(err)
		}

		dhs, err := db.GetDirectories()
		if err != nil {
			t.Error(err)
		}

		if len(dhs) != len(dirs)-1 {
			t.Error("len(dhs) != len(dirs) - 1")
		}
	})
}

func TestVolume(t *testing.T) {
	t.Run("Put Volume", func(t *testing.T) {
		for _, vol := range vols {
			err := db.PutVolume(vol)
			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Get Volume", func(t *testing.T) {
		vhs, err := db.GetVolumes()
		if err != nil {
			t.Error(err)
		}
		if len(vhs) != len(vols) {
			t.Error("len(vhs) != len(vols)", len(vhs), len(vols))
		}
		for _, vol := range vols {
			found := false
			for _, vh := range vhs {
				if vh.Volume == vol {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("can not find %s", vol)
			}
		}
	})

	t.Run("Remove Volume", func(t *testing.T) {
		vol := vols[rand.Intn(len(vols))]
		err := db.RemoveVolume(vol)
		if err != nil {
			t.Error(err)
		}

		vhs, err := db.GetVolumes()
		if err != nil {
			t.Error(err)
		}

		if len(vhs) != len(vols)-1 {
			t.Error("len(vhs) != len(vols) - 1")
		}
	})
}

func TestSpace(t *testing.T) {
	t.Run("Put Space Info", func(t *testing.T) {
		dhs, err := db.GetDirectories()
		if err != nil {
			t.Error(err)
		}
		vhs, err := db.GetVolumes()
		if err != nil {
			t.Error(err)
		}

		info := model.SpaceInfo{
			Time:             time.Now().Unix(),
			DirectoriesSpace: []model.DirectorySpace{},
			VolumesSpace:     []model.VolumeSpace{},
		}

		for _, dh := range dhs {
			info.DirectoriesSpace = append(info.DirectoriesSpace, model.DirectorySpace{
				Directory:  dh.Directory,
				UsedSpace:  rand.Uint64(),
				Percentage: rand.Float64(),
			})
		}

		for _, vh := range vhs {
			info.VolumesSpace = append(info.VolumesSpace, model.VolumeSpace{
				Volume:         vh.Volume,
				AvailableSpace: rand.Uint64(),
				TotalSpace:     rand.Uint64(),
				FreeSpace:      rand.Uint64(),
			})
		}

		db.PutSpaceInfo(info)
	})

	t.Run("Get Space", func(t *testing.T) {
		dhs, err := db.GetDirectories()
		if err != nil {
			t.Error(err)
		}
		vhs, err := db.GetVolumes()
		if err != nil {
			t.Error(err)
		}

		if len(dhs) == 0 || len(vhs) == 0 {
			t.Error(len(dhs) == 0 || len(vhs) == 0)
			return
		}

		dir := dhs[rand.Intn(len(dhs)-1)].Directory
		sps, err := db.GetDirectorySpace(dir, time.Now().Unix() - 1e5, time.Now().Unix())
		if err != nil {
			t.Error(err)
		}

		if len(sps) == 0 {
			t.Error("empty sps")
			return
		}

		sp := sps[0]
		t.Log(sp.Time, sp.Directory, sp.UsedSpace, sp.Percentage)
	})
}

func TestClose(t *testing.T) {
	err := db.Reset()
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = db.Close()
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestMain(m *testing.M) {
	var err error
	db, err = database.Open()
	if err != nil {
		println(err.Error())
		return
	}

	code := m.Run()
	os.Exit(code)
}
