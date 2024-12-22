package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	dbName              string
	schema              string
	db                  *sql.DB
	lock                sync.Mutex
	NeedsInitialization bool
}

func NewDatabase(dbName string, schema []byte) *Database {
	database := &Database{
		dbName:              dbName,
		schema:              string(schema),
		NeedsInitialization: false,
	}

	_, err := os.Stat(dbName)
	if err != nil {
		database.NeedsInitialization = true
	}

	database.db, err = sql.Open("sqlite3", database.dbName)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	_, err = database.db.Exec(database.schema)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	return database
}

func (cdb *Database) Close() error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	return cdb.db.Close()
}

func (cdb *Database) InsertData(tableName string, sqlStatement []byte, datavalues [][]string) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	queryTemplate := string(sqlStatement)

	query := fmt.Sprintf(queryTemplate, tableName)

	tx, err := cdb.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, row := range datavalues {
		valuesInterface := make([]interface{}, len(row))
		for i, v := range row {
			valuesInterface[i] = v
		}

		_, err = stmt.Exec(valuesInterface...)
		if err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}

	return tx.Commit()
}

func (cdb *Database) GetColumnValues(tableName string, sqlStatement []byte, columnName string) ([]map[string]interface{}, error) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	queryTemplate := string(sqlStatement)

	query := fmt.Sprintf(queryTemplate, columnName, tableName)

	rows, err := cdb.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := SqlRowsToArray(rows)

	return result, err
}

func (cdb *Database) GetRowsByColumnValue(tableName string, sqlStatement []byte, columnName string, searchValue string) ([]map[string]interface{}, error) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	queryTemplate := string(sqlStatement)

	query := fmt.Sprintf(queryTemplate, tableName, columnName)

	rows, err := cdb.db.Query(query, searchValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := SqlRowsToArray(rows)

	return result, err
}

func SqlRowsToArray(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		err := rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, column := range columns {
			row[column] = values[i]
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
