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
func SaveSpace(db database.DB, fn func(*model.SpaceRecord, error)) {
	for {
		last, _ := db.GetLastRecord()
		next := last.Time + int64(SaveDuration.Seconds())
		now := time.Now().Unix()
		duration := int64(1)
		if next > now {
			duration = next - now
		}
		timer := time.NewTimer(time.Second * time.Duration(duration))
		<-timer.C
		timer.Stop()
		record, err := save(db)
		fn(record, err)
	}
}

func save(db database.DB) (*model.SpaceRecord, error) {
	dhs, err := db.GetDirectories()
	if err != nil {
		return nil, err
	}
	vhs, err := db.GetVolumes()
	if err != nil {
		return nil, err
	}

	record := model.SpaceRecord{
		Time:             time.Now().Unix(),
		DirectoriesSpace: []model.DirectorySpace{},
		VolumesSpace:     []model.VolumeSpace{},
	}

	for _, dh := range dhs {
		s, err := space.GetDirectorySpace(dh.Directory)
		if err != nil {
			return nil, err
		}
		s.Limit = dh.Limit
		record.DirectoriesSpace = append(record.DirectoriesSpace, s)
	}

	for _, vh := range vhs {
		s, err := space.GetVolumeSpace(vh.Volume)
		if err != nil {
			return nil, err
		}
		s.Limit = vh.Limit
		s.LimitPercentage = vh.LimitPercentage
		record.VolumesSpace = append(record.VolumesSpace, s)
	}

	return &record, db.PutLastRecord(record)
}
