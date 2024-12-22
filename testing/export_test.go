package database

import (
	"os"
	"testing"

	"github.com/williamveith/wastetags/pkg/database"
)

func TestToProtobuf(t *testing.T) {
	type Config struct {
		DatabasePath string `json:"database_path"`
	}
	cfg := &Config{
		DatabasePath: "/Users/main/Projects/Go/wastetags/build/wastetags.sqlite3",
	}

	sqlStatement, _ := os.ReadFile("/Users/main/Projects/Go/wastetags/cmd/wastetags/query/schema.sql")
	db := database.NewDatabase(cfg.DatabasePath, sqlStatement)
	defer db.Close()
	db.ExportToProtobuf()
}
