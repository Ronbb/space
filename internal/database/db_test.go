package database_test

import (
	"os"
	"testing"

	"github.com/ronbb/space/internal/database"
)

func TestOpen(t *testing.T) {
	t.Run("Open DB", func(t *testing.T) {
		db, err := database.Open()
		if err != nil {
			t.Error(err)
			return
		}

		defer db.Close()
	})
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
