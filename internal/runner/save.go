package runner

import (
	"time"

	"github.com/ronbb/space/internal/database"
	"github.com/ronbb/space/internal/model"
	"github.com/ronbb/space/internal/space"
)

// value 1
const (
	SaveDuration = 10 * time.Second
)

// SaveSpace .
func SaveSpace(db database.DB, fn func(error)) {
	for {
		last, _ := db.GetLastRecordTime()
		next := last + int64(SaveDuration.Seconds())
		now := time.Now().Unix()
		duration := int64(1)
		if next > now {
			duration = next - now
		}
		timer := time.NewTimer(time.Second * time.Duration(duration))
		<-timer.C
		timer.Stop()
		err := save(db)
		fn(err)
	}
}

func save(db database.DB) error {
	dhs, err := db.GetDirectories()
	if err != nil {
		return err
	}
	vhs, err := db.GetVolumes()
	if err != nil {
		return err
	}

	info := model.SpaceInfo{
		Time:             time.Now().Unix(),
		DirectoriesSpace: []model.DirectorySpace{},
		VolumesSpace:     []model.VolumeSpace{},
	}

	for _, dh := range dhs {
		s, err := space.GetDirectorySpace(dh.Directory)
		if err != nil {
			return err
		}
		info.DirectoriesSpace = append(info.DirectoriesSpace, s)
	}

	for _, vh := range vhs {
		s, err := space.GetVolumeSpace(vh.Volume)
		if err != nil {
			return err
		}
		info.VolumesSpace = append(info.VolumesSpace, s)
	}

	return db.PutSpaceInfo(info)
}
