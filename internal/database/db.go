package database

import "github.com/tidwall/buntdb"

const (
	dbName = "data.db"
)

// DB .
type DB interface {
	createIndexes()
	Test()
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

	new.createIndexes()

	return &new, nil
}

type index struct {
	name      string
	pattern   string
	jsonKey   string
	decending bool
}

var (
	indexes = []index{
		{
			name:      "global.record",
			pattern:   "global:record",
			jsonKey:   "time",
			decending: true,
		},
		{
			name:    "dir.name",
			pattern: "name:dir:*",
			jsonKey: "name",
		},
		{
			name:    "dir.time",
			pattern: "space:dir:*:*",
			jsonKey: "time",
		},
		{
			name:    "vol.time",
			pattern: "space:vol:*:*",
			jsonKey: "time",
		},
	}
)

func (db *db) createIndexes() {
	for _, index := range indexes {
		less := buntdb.IndexJSON(index.jsonKey)
		if index.decending {
			less = buntdb.Desc(less)
		}
		db.origin.CreateIndex(index.name, index.pattern, less)
	}
}

func (db *db) Test() {
	db.origin.Update(func(tx *buntdb.Tx) error {
		tx.Set("space:aaa:123", `{"time": 0}`, nil)
		tx.Set("space:aab:124", `{"time": 1}`, nil)
		tx.Set("space:aca:223", `{"time": 2}`, nil)
		tx.Set("space:daa:133", `{"time": 7}`, nil)
		return nil
	})

	db.origin.View(func(tx *buntdb.Tx) error {
		tx.Ascend("time", func(key, value string) bool {
			println(key, value)
			return true
		})

		return nil
	})

	h1, _ := hash("C:\\1a")
	h2, _ := hash("C:\\1b")
	
	println("hash", h1, h2)

}
