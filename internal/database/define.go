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
	keyLastRecordTime = "global.last.record"
)

// template key
const (
	keyTemplateDirectoryHash  = "hash:dir:%s"     // hash
	keyTemplateDirectorySpace = "space:dir:%s:%s" // hash:time
	keyTemplateVolumeHash     = "hash:vol:%s"     // hash
	keyTemplateVolumeSpace    = "space:vol:%s:%s" // hash:time
)

// pattern
var (
	patternDirectoryHash  = pattern(keyTemplateDirectoryHash)
	patternDirectorySpace = pattern(keyTemplateDirectorySpace)
	patternVolumeHash     = pattern(keyTemplateVolumeHash)
	patternVolumeSpace    = pattern(keyTemplateVolumeSpace)
)

// index name
const (
	indexDirectoryHash  = "dir.hash"
	indexDirectorySpace = "dir.time"
	indexVolumeHash     = "vol.hash"
	indexVolumeSpace    = "vol.time"
)

var indexes = []index{
	{
		name:      indexDirectoryHash,
		pattern:   patternDirectoryHash,
		jsonKey:   "hash",
		decending: true,
	},
	{
		name:    indexDirectorySpace,
		pattern: patternDirectorySpace,
		jsonKey: "time",
	},
	{
		name:      indexVolumeHash,
		pattern:   patternVolumeHash,
		jsonKey:   "hash",
		decending: true,
	},
	{
		name:      indexVolumeSpace,
		pattern:   patternVolumeSpace,
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

func timeJSON(t interface{}) (string, error) {
	b, err := json.Marshal(map[string]interface{}{
		"time": t,
	})
	if err != nil {
		return "", err
	}

	return string(b), nil
}
