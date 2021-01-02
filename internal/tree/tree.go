package tree

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/ronbb/space/internal/model"
	"github.com/ronbb/space/internal/usn"
	"github.com/tidwall/buntdb"
)

const (
	dbName          = "tree.db"
	keyLastEnumTime = "last:enum:time"
	// keyJournalID    = "journal:id"
	keyLastUSN      = "journal:last:usn"
	keyTamplateFile = "file:%s"
	patternFile     = "file:*"
	indexFile       = "file"
)

// Tree .
type Tree struct {
	origin *buntdb.DB
	handle syscall.Handle
}

// Open .
func Open() (*Tree, error) {
	origin, err := buntdb.Open(dbName)
	if err != nil {
		return nil, err
	}

	new := Tree{
		origin: origin,
	}

	volume := usn.Volume("C:\\")
	isNTFS, err := usn.IsNTFS(volume)
	if err != nil {
		return nil, err
	}
	if !isNTFS {
		return nil, errors.New("is not ntfs")
	}

	handle, err := usn.NewHandle(volume)
	if err != nil {
		return nil, err
	}

	new.handle = handle
	new.origin.CreateIndex(indexFile, patternFile, buntdb.IndexJSON("referenceNumber"))

	return &new, new.init()
}

func (t *Tree) init() error {
	var err error
	err = t.origin.View(func(tx *buntdb.Tx) error {
		_, err = tx.Get(keyLastUSN)
		return err
	})

	if err != nil {
		err = t.createJournal()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Tree) createJournal() error {
	return usn.CreateJournal(t.handle)
}

func (t *Tree) queryJournal() (usn.QueryUSNJournalData, error) {
	return usn.QueryJournal(t.handle)
}

func (t *Tree) enumNext(onEnum func(*model.Node)) error {
	var err error
	lastUSN := ""
	err = t.origin.View(func(tx *buntdb.Tx) error {
		lastUSN, err = tx.Get(keyLastUSN)
		return err
	})
	if lastUSN == "" {
		lastUSN = "0"
	}
	i, err := strconv.ParseInt(lastUSN, 10, 64)
	if err != nil {
		return err
	}
	q, err := t.queryJournal()
	if err != nil {
		return err
	}

	err = usn.EnumJournal(t.handle, i+1, q.NextUsn, onEnum)
	if err != nil {
		return err
	}

	err = t.origin.Update(func(tx *buntdb.Tx) error {
		_, _, err = tx.Set(keyLastUSN, fmt.Sprintf("%d", q.NextUsn), nil)
		return err
	})

	return err
}

// Save .
func (t *Tree) Save() error {
	return t.enumNext(func(n *model.Node) {
		b16, err := base64.StdEncoding.DecodeString(n.ReferenceNumber)
		size, err := usn.GetSizeByID(t.handle, b16)
		if err != nil {
			return
		}
		n.Size = uint64(size)
		b, err := json.Marshal(n)
		if err != nil {
			return
		}
		t.origin.Update(func(tx *buntdb.Tx) error {
			_, _, err = tx.Set(fmt.Sprintf(keyTamplateFile, n.ReferenceNumber), string(b), nil)
			return err
		})
	})
}

// GetSize .
func (t *Tree) GetSize(path string) (uint64, error) {
	sp := strings.TrimLeft(path, filepath.VolumeName(path))
	sp = strings.Trim(sp, "\\")
	ps := strings.Split(sp, "\\")
	pre := model.Node{
		ReferenceNumber: usn.RootBase64,
	}

	for _, p := range ps {
		tmp := model.Node{}
		found := false
		err := t.origin.View(func(tx *buntdb.Tx) error {
			err := tx.AscendEqual(indexFile, fmt.Sprintf(`{"parentReferenceNumber":"%s"}`, pre.ReferenceNumber), func(key, value string) bool {
				err := json.Unmarshal([]byte(value), &tmp)
				if err != nil {
					return true
				}
				if tmp.FileName == p {
					println(tmp.FileName)
					found = true
					pre = tmp
					return false
				}

				return true
			})

			if !found {
				return err
			}
			return nil
		})
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}
