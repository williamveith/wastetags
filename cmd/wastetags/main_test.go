package main

import (
	"os"
	"testing"

	"github.com/williamveith/wastetags/pkg/database"
)

func TestMainFun(t *testing.T) {
	cfg = &Config{
		DatabasePath: "/Users/main/Projects/Go/wastetags/build/wastetags.sqlite3",
	}

	sqlStatement, _ := os.ReadFile("/Users/main/Projects/Go/wastetags/cmd/wastetags/query/schema.sql")
	db = database.NewDatabase(cfg.DatabasePath, sqlStatement)
	main()
}
