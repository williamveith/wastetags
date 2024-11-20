package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type ChemicalDatabase struct {
	dbName string
	db     *sql.DB
	lock   sync.Mutex
}

func NewChemicalDatabase(dbName string) *ChemicalDatabase {
	db, _ := sql.Open("sqlite3", dbName)

	cdb := &ChemicalDatabase{
		dbName: dbName,
		db:     db,
	}

	cdb.createTable()

	return cdb
}

func (cdb *ChemicalDatabase) createTable() error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	_, err := cdb.db.Exec(`
		CREATE TABLE IF NOT EXISTS chemicals (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chem_name TEXT,
			component_name TEXT,
			cas TEXT,
			percent TEXT,
			component_order INTEGER
		)
	`)
	return err
}

func (cdb *ChemicalDatabase) InsertData(datavalues map[string][]map[string]string) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	tx, err := cdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO chemicals (chem_name, component_name, cas, percent, component_order)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for chemName, components := range datavalues {
		for order, component := range components {
			_, err := stmt.Exec(chemName, component["Name"], component["CAS"], component["Percent"], order)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (cdb *ChemicalDatabase) GetRowsByName(chemName string) ([]map[string]interface{}, error) {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	rows, err := cdb.db.Query("SELECT * FROM chemicals WHERE chem_name = ?", chemName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return result, nil
}

func (cdb *ChemicalDatabase) Close() error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	return cdb.db.Close()
}
