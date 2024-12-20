package database

import (
	"os"
	"testing"

	"github.com/williamveith/wastetags/pkg/database"
	"google.golang.org/protobuf/proto"
)

func TestToProtobuf(t *testing.T) {
	type Config struct {
		DatabasePath string `json:"database_path"`
	}
	cfg := &Config{
		DatabasePath: "/Users/main/Projects/Go/wastetags/data/chemicals.sqlite3",
	}

	sqlStatement, _ := os.ReadFile("/Users/main/Projects/Go/wastetags/cmd/wastetags/query/schema.sql")
	db := database.NewDatabase(cfg.DatabasePath, sqlStatement)

	defer db.Close()

	// Run the ToProtobuf function
	outputFile := "mixtures.bin"

	err := db.ToProtobuf(
		"mixtures",
		&database.Mixture{},
		&database.MixtureList{},
		outputFile,
	)

	if err != nil {
		t.Fatalf("ToProtobuf failed: %v", err)
	}

	// Read and unmarshal the Protobuf file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read Protobuf output file: %v", err)
	}

	var result database.ChemicalList
	err = proto.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal Protobuf data: %v", err)
	}
}
