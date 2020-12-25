package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ronbb/space/internal/model"
	"github.com/ronbb/space/internal/utils"
	"github.com/tidwall/buntdb"
)

const (
	dbName = "data.db"
)

// DB .
type DB interface {
	Close() error
	Reset() error

	PutDirectory(dir string) error
	RemoveDirectory(dir string) error
	GetDirectories() ([]model.DirectoryHash, error)
	PutVolume(vol string) error
	RemoveVolume(vol string) error
	GetVolumes() ([]model.VolumeHash, error)

	PutLastRecord(info model.SpaceRecord) error

	GetDirectorySpace(dir string, start, end int64) ([]model.DirectorySpace, error)
	GetVolumeSpace(vol string, start, end int64) ([]model.VolumeSpace, error)

	GetLastRecord() (model.SpaceRecord, error)
}

type db struct {
	origin *buntdb.DB
}

// Open .
func Open() (DB, error) {
	origin, err := buntdb.Open(dbName)
	if err != nil {
		return nil, err
	}

	new := db{
		origin: origin,
	}

	err = new.createIndexes()
	if err != nil {
		return nil, err
	}

	return &new, nil
}

func (db *db) GetLastRecord() (model.SpaceRecord, error) {
	record := model.SpaceRecord{}
	err := db.origin.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(keyLastRecord)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(v), &record)
	})
	if err != nil {
		return record, err
	}

	return record, nil
}

// func (db *db) SetLastRecordTime(t int64) error {
// 	return db.origin.Update(func(tx *buntdb.Tx) error {
// 		_, _, err := tx.Set(keyLastRecordTime, strconv.FormatInt(t, 10), nil)
// 		return err
// 	})
// }

func (db *db) Close() error {
	return db.origin.Close()
}

func (db *db) Reset() error {
	err := db.origin.Update(func(tx *buntdb.Tx) error {
		return tx.DeleteAll()
	})
	if err != nil {
		return err
	}

	return db.origin.Shrink()
}

func createDirectoryIndex(tx *buntdb.Tx, hash string) error {
	return tx.CreateIndex(indexDirectorySpace(hash), patternDirectorySpace(hash), buntdb.IndexJSON("time"))
}

func createVolumeIndex(tx *buntdb.Tx, hash string) error {
	return tx.CreateIndex(indexVolSpace(hash), patternVolumeSpace(hash), buntdb.IndexJSON("time"))
}

func (db *db) createIndexes() error {
	for _, index := range indexes {
		less := buntdb.IndexJSON(index.jsonKey)
		if index.decending {
			less = buntdb.Desc(less)
		}
		err := db.origin.CreateIndex(index.name, index.pattern, less)
		if err != nil {
			return err
		}
	}

	dhs, err := db.GetDirectories()
	if err != nil {
		return err
	}

	vhs, err := db.GetVolumes()
	if err != nil {
		return err
	}

	db.origin.Update(func(tx *buntdb.Tx) error {
		for _, dh := range dhs {
			err := createDirectoryIndex(tx, dh.Hash)
			if err != nil {
				return err
			}
		}
		for _, vh := range vhs {
			err := createVolumeIndex(tx, vh.Hash)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func (db *db) PutDirectory(dir string) error {
	return db.origin.Update(func(tx *buntdb.Tx) error {
		key, hash, err := keyDirectoryHash(dir)
		if err != nil {
			return err
		}

		value, err := json.Marshal(&model.DirectoryHash{
			Directory: dir,
			Hash:      hash,
		})
		if err != nil {
			return err
		}

		_, _, err = tx.Set(key, string(value), nil)
		if err != nil {
			return err
		}
		err = createDirectoryIndex(tx, hash)
		return err
	})
}

func (db *db) RemoveDirectory(dir string) error {
	return db.origin.Update(func(tx *buntdb.Tx) error {
		key, _, err := keyDirectoryHash(dir)
		if err != nil {
			return err
		}

		_, err = tx.Delete(key)
		return err
	})
}

func (db *db) GetDirectories() ([]model.DirectoryHash, error) {
	dirs := []model.DirectoryHash{}

	err := db.origin.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(indexDirectoryHash, func(key, value string) bool {
			dir := model.DirectoryHash{}
			err := json.Unmarshal([]byte(value), &dir)
			if err != nil {
				return true
			}
			dirs = append(dirs, dir)
			return true
		})
	})

	return dirs, err
}

func (db *db) PutVolume(vol string) error {
	vol = filepath.VolumeName(vol) + "\\"
	return db.origin.Update(func(tx *buntdb.Tx) error {
		key, hash, err := keyVolumeHash(vol)
		if err != nil {
			return err
		}

		value, err := json.Marshal(&model.VolumeHash{
			Volume: vol,
			Hash:   hash,
		})
		if err != nil {
			return err
		}

		_, _, err = tx.Set(key, string(value), nil)
		if err != nil {
			return err
		}
		err = createVolumeIndex(tx, hash)
		return err
	})
}

func (db *db) RemoveVolume(vol string) error {
	return db.origin.Update(func(tx *buntdb.Tx) error {
		key, _, err := keyVolumeHash(vol)
		if err != nil {
			return err
		}

		_, err = tx.Delete(key)
		return err
	})
}

func (db *db) GetVolumes() ([]model.VolumeHash, error) {
	dirs := []model.VolumeHash{}

	err := db.origin.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(indexVolumeHash, func(key, value string) bool {
			dir := model.VolumeHash{}
			err := json.Unmarshal([]byte(value), &dir)
			if err != nil {
				return true
			}
			dirs = append(dirs, dir)
			return true
		})
	})

	return dirs, err
}

func (db *db) PutLastRecord(record model.SpaceRecord) error {
	if record.DirectoriesSpace == nil || record.VolumesSpace == nil {
		return errors.New("info.DirectorirsSpace or info.VolumesSpace is nil")
	}

	// aligned time
	t := record.Time
	tStr := fmt.Sprintf("%d", t)

	err := db.origin.Update(func(tx *buntdb.Tx) error {
		for _, dirSpace := range record.DirectoriesSpace {
			dirSpace.Time = t
			key, err := keyDirectorySpace(dirSpace.Directory, tStr)
			if err != nil {
				return err
			}

			value, err := json.Marshal(&dirSpace)
			if err != nil {
				return err
			}

			_, _, err = tx.Set(key, string(value), &buntdb.SetOptions{
				Expires: true,
				TTL:     TTL,
			})

			if err != nil {
				return err
			}
		}

		for _, volSpace := range record.VolumesSpace {
			volSpace.Time = t
			key, err := keyVolumeSpace(volSpace.Volume, tStr)
			if err != nil {
				return err
			}

			value, err := json.Marshal(&volSpace)
			if err != nil {
				return err
			}

			_, _, err = tx.Set(key, string(value), &buntdb.SetOptions{
				Expires: true,
				TTL:     TTL,
			})

			if err != nil {
				return err
			}
		}

		b, err := json.Marshal(&record)
		if err != nil {
			return err
		}
		tx.Set(keyLastRecord, string(b), &buntdb.SetOptions{
			Expires: true,
			TTL:     TTL,
		})

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (db *db) GetDirectorySpace(dir string, start, end int64) ([]model.DirectorySpace, error) {
	spaces := []model.DirectorySpace{}
	hash, err := utils.HashPath(dir)
	index := indexDirectorySpace(hash)

	startJSON, err := timeJSON(start)
	if err != nil {
		return nil, err
	}

	endJSON, err := timeJSON(end)
	if err != nil {
		return nil, err
	}

	db.origin.View(func(tx *buntdb.Tx) error {
		return tx.AscendRange(index, startJSON, endJSON, func(key, value string) bool {
			space := model.DirectorySpace{}
			err := json.Unmarshal([]byte(value), &space)
			if err != nil {
				return true
			}

			spaces = append(spaces, space)
			return true
		})
	})

	return spaces, err
}

func (db *db) GetVolumeSpace(vol string, start, end int64) ([]model.VolumeSpace, error) {
	spaces := []model.VolumeSpace{}
	hash, err := utils.HashPath(vol)
	index := indexVolSpace(hash)
	if err != nil {
		return nil, err
	}

	startJSON, err := timeJSON(start)
	if err != nil {
		return nil, err
	}

	endJSON, err := timeJSON(end)
	if err != nil {
		return nil, err
	}

	db.origin.View(func(tx *buntdb.Tx) error {
		return tx.AscendRange(index, startJSON, endJSON, func(key, value string) bool {
			space := model.VolumeSpace{}
			err := json.Unmarshal([]byte(value), &space)
			if err != nil {
				return true
			}

			spaces = append(spaces, space)
			return true
		})
	})

	return spaces, err
}
