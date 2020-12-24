package database_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/ronbb/space/internal/database"
)

var (
	db   database.DB
	dirs []string
	vols []string
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
	return fmt.Sprintf("C:\\%d\\", rand.Int())
}

func randVol() string {
	return fmt.Sprintf("%s:\\", string(rune(int('A')+rand.Intn(25))))
}

func TestPutDirectory(t *testing.T) {
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
}

func TestMain(m *testing.M) {
	var err error
	db, err = database.Open()
	if err != nil {
		println(err.Error())
		return
	}

	defer func() {
		err = db.Reset()
		if err != nil {
			println(err.Error())
			return
		}
		err = db.Close()
		if err != nil {
			println(err.Error())
			return
		}
	}()

	code := m.Run()
	os.Exit(code)
}
