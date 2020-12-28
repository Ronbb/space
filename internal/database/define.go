package database

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ronbb/space/internal/utils"
)

type index struct {
	name      string
	pattern   string
	jsonKey   string
	decending bool
}

// TTL 过期时间
const TTL = time.Hour * 24 * 30

// key
const (
	keyLastRecord = "global.last.record"
)

// template key
const (
	keyTemplateDirectoryHash  = "hash:dir:%s"     // hash
	keyTemplateDirectorySpace = "space:dir:%s:%s" // hash:time
	keyTemplateVolumeHash     = "hash:vol:%s"     // hash
	keyTemplateVolumeSpace    = "space:vol:%s:%s" // hash:time
	keyTemplateRecord         = "record:%s"       // time
)

// pattern
var (
	patternDirectoryHash = pattern(keyTemplateDirectoryHash)
	patternVolumeHash    = pattern(keyTemplateVolumeHash)
	patternRecord        = pattern(keyTemplateRecord)
)

func patternVolumeSpace(hash string) string {
	return fmt.Sprintf(keyTemplateVolumeSpace, hash, "*")
}

func patternDirectorySpace(hash string) string {
	return fmt.Sprintf(keyTemplateDirectorySpace, hash, "*")
}

// index name
const (
	indexDirectoryHash = "dir.hash"
	indexVolumeHash    = "vol.hash"
	indexRecord        = "record"
)

const (
	indexTemplateDirectorySpace = "dir.time.%s"
	indexTemplateVolumeSpace    = "vol.time.%s"
)

var indexes = []index{
	{
		name:      indexDirectoryHash,
		pattern:   patternDirectoryHash,
		jsonKey:   "hash",
		decending: true,
	},
	{
		name:      indexVolumeHash,
		pattern:   patternVolumeHash,
		jsonKey:   "hash",
		decending: true,
	},
	{
		name:      indexRecord,
		pattern:   patternRecord,
		jsonKey:   "time",
		decending: true,
	},
}

func pattern(key string) string {
	return strings.ReplaceAll(key, "%s", "*")
}

func keyDirectoryHash(dir string) (key, hash string, err error) {
	hash, err = utils.HashPath(dir)
	if err != nil {
		return
	}
	key = fmt.Sprintf(keyTemplateDirectoryHash, hash)
	return
}

func keyDirectorySpace(dir string, time string) (key string, err error) {
	hash, err := utils.HashPath(dir)
	if err != nil {
		return
	}
	key = fmt.Sprintf(keyTemplateDirectorySpace, hash, time)
	return
}

func keyVolumeHash(dir string) (key, hash string, err error) {
	hash, err = utils.HashPath(dir)
	if err != nil {
		return
	}
	key = fmt.Sprintf(keyTemplateVolumeHash, hash)
	return
}

func keyVolumeSpace(dir string, time string) (key string, err error) {
	h, err := utils.HashPath(dir)
	if err != nil {
		return
	}
	key = fmt.Sprintf(keyTemplateVolumeSpace, h, time)
	return
}

func keyRecord(time string) (key string) {
	key = fmt.Sprintf(keyTemplateRecord, time)
	return
}

func indexDirectorySpace(hash string) string {
	return fmt.Sprintf(indexTemplateDirectorySpace, hash)
}

func indexVolSpace(hash string) string {
	return fmt.Sprintf(indexTemplateVolumeSpace, hash)
}

func timeJSON(t interface{}) (string, error) {
	b, err := json.Marshal(map[string]interface{}{
		"time": t,
	})
	if err != nil {
		return "", err
	}

	return string(b), nil
}
