package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	dbName string
	schema string
	db     *sql.DB
	lock   sync.Mutex
}

func NewDatabase(dbName string, schema []byte) *Database {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}

	// Execute the schema
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	return &Database{
		dbName: dbName,
		schema: string(schema),
		db:     db,
	}
}

func (cdb *Database) InsertData(tableName, sqlStatement []byte, datavalues map[string][]map[string]string) error {
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

	// Iterate over datavalues to insert rows
	for key, rows := range datavalues {
		for index, row := range rows {
			params := []interface{}{key}
			for _, value := range row {
				params = append(params, value)
			}
			params = append(params, index)

			if _, err := stmt.Exec(params...); err != nil {
				return fmt.Errorf("failed to execute statement: %w", err)
			}
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

func (cdb *Database) Close() error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	return cdb.db.Close()
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
			row[column] = values[i] // Ensure column names are strings
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
