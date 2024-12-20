package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	reflect "reflect"
	"strings"
	"sync"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/proto"
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

func dbColumnNamesToProtobufKey(columnNames []string) []string {
	titleCase := cases.Title(language.English)
	for i, columnName := range columnNames {
		columnName = strings.ReplaceAll(columnName, "_", " ")
		columnNames[i] = strings.ReplaceAll(titleCase.String(columnName), " ", "")
	}
	return columnNames
}

func protobufKeyToDbColumnNames(protobufKeys []string) []string {
	for i, protobufKey := range protobufKeys {
		// Split the camel case into words
		var words []string
		wordStart := 0
		for j, r := range protobufKey {
			if j > 0 && unicode.IsUpper(r) {
				words = append(words, protobufKey[wordStart:j])
				wordStart = j
			}
		}
		words = append(words, protobufKey[wordStart:])

		// Convert words to lowercase and join with underscores
		protobufKeys[i] = strings.ToLower(strings.Join(words, "_"))
	}
	return protobufKeys
}

// ToProtobuf generalizes the conversion of a database table to a Protobuf message
func (cdb *Database) ToProtobuf(tableName string, protoItemType proto.Message, protoCollectionType proto.Message, outputFile string) error {
	// Lock the database
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	// Fetch all rows from the specified table
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := cdb.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query table %s: %w", tableName, err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to fetch column names: %w", err)
	}

	// Prepare reflection details
	itemType := reflect.TypeOf(protoItemType).Elem()
	collectionValue := reflect.ValueOf(protoCollectionType).Elem()
	repeatedField := collectionValue.FieldByName(reflect.TypeOf(protoItemType).Elem().Name())
	if !repeatedField.IsValid() || repeatedField.Kind() != reflect.Slice {
		return fmt.Errorf("protoCollectionType does not contain a valid repeated field")
	}
	columnNames, _ := rows.Columns()
	protobufKeys := dbColumnNamesToProtobufKey(columnNames)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		item := reflect.New(itemType).Interface().(proto.Message)

		// Populate fields dynamically
		itemValue := reflect.ValueOf(item).Elem()
		for i, protobufKey := range protobufKeys {
			field := itemValue.FieldByName(protobufKey)
			if field.IsValid() && field.CanSet() {
				// Convert database value to the correct type and set it
				dbValue := reflect.ValueOf(values[i])
				if dbValue.IsValid() && dbValue.Type().ConvertibleTo(field.Type()) {
					field.Set(dbValue.Convert(field.Type()))
				}
			}
		}

		// Append the item to the repeated field
		repeatedField.Set(reflect.Append(repeatedField, itemValue.Addr()))
	}

	// Serialize the Protobuf collection to binary format
	data, err := proto.Marshal(protoCollectionType)
	if err != nil {
		return fmt.Errorf("failed to marshal Protobuf data: %w", err)
	}

	// Save to a binary file
	err = os.WriteFile(outputFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write Protobuf file: %w", err)
	}

	fmt.Printf("Protobuf data saved to %s\n", outputFile)
	return nil
}

func (cdb *Database) FromProtobuf(tableName string, protoItemType proto.Message, protoCollectionType proto.Message, data []byte) error {
	// Lock the database
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	// Read the Protobuf binary file

	// Deserialize the Protobuf data
	if err := proto.Unmarshal(data, protoCollectionType); err != nil {
		return fmt.Errorf("failed to unmarshal Protobuf data: %w", err)
	}

	// Prepare reflection details
	collectionValue := reflect.ValueOf(protoCollectionType).Elem()
	repeatedField := collectionValue.FieldByName(reflect.TypeOf(protoItemType).Elem().Name())
	if !repeatedField.IsValid() || repeatedField.Kind() != reflect.Slice {
		return fmt.Errorf("protoCollectionType does not contain a valid repeated field")
	}

	// Iterate over the repeated field
	for i := 0; i < repeatedField.Len(); i++ {
		item := repeatedField.Index(i).Interface().(proto.Message)
		itemValue := reflect.ValueOf(item).Elem()

		// Extract fields dynamically
		columns := []string{}
		placeholders := []string{}
		values := []interface{}{}

		for j := 0; j < itemValue.NumField(); j++ {
			field := itemValue.Type().Field(j)
			fieldValue := itemValue.Field(j)

			if fieldValue.IsValid() && fieldValue.CanInterface() {
				columnName := protobufKeyToDbColumnNames([]string{field.Name})[0]
				columns = append(columns, columnName)
				placeholders = append(placeholders, "?")
				values = append(values, fieldValue.Interface())
			}
		}

		// Construct the INSERT statement
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

		// Execute the INSERT statement
		if _, err := cdb.db.Exec(query, values...); err != nil {
			return fmt.Errorf("failed to insert row into table %s: %w", tableName, err)
		}
	}

	fmt.Printf("Protobuf data successfully imported into table %s\n", tableName)
	return nil
}

func (cdb *Database) ExportToProtobuf() {
	// Define the tables and their corresponding Protobuf types
	tableMappings := []struct {
		name        string
		message     proto.Message
		listMessage proto.Message
	}{
		{"mixtures", &Mixture{}, &MixtureList{}},
		{"chemicals", &Chemical{}, &ChemicalList{}},
		{"locations", &Location{}, &LocationList{}},
		{"containers", &Container{}, &ContainerList{}},
		{"units", &Unit{}, &UnitList{}},
		{"states", &State{}, &StateList{}},
	}

	// Loop through each table and export its data to Protobuf
	for _, mapping := range tableMappings {
		outputFile := "/Users/main/Projects/Go/wastetags/cmd/wastetags/data/" + mapping.name + ".bin" // File name is based on the table name
		err := cdb.ToProtobuf(
			mapping.name,        // Table name
			mapping.message,     // Single-row Protobuf message
			mapping.listMessage, // List Protobuf message
			outputFile,          // Output file
		)
		if err != nil {
			log.Fatalf("Failed to export table %s: %v", mapping.name, err)
		} else {
			log.Printf("Exported table %s to %s", mapping.name, outputFile)
		}
	}
}

func (cdb *Database) ImportFromProtobuff(embeddedData embed.FS) {
	// Define the tables and their corresponding Protobuf types
	tableMappings := []struct {
		name        string
		message     proto.Message
		listMessage proto.Message
	}{
		{"chemicals", &Chemical{}, &ChemicalList{}},
		{"mixtures", &Mixture{}, &MixtureList{}},
		{"locations", &Location{}, &LocationList{}},
		{"containers", &Container{}, &ContainerList{}},
		{"units", &Unit{}, &UnitList{}},
		{"states", &State{}, &StateList{}},
	}

	// Loop through each table and export its data to Protobuf
	for _, mapping := range tableMappings {
		inputFile, _ := embeddedData.ReadFile("data/" + mapping.name + ".bin") // File name is based on the table name
		err := cdb.FromProtobuf(
			mapping.name,        // Table name
			mapping.message,     // Single-row Protobuf message
			mapping.listMessage, // List Protobuf message
			inputFile,           // Output file
		)
		if err != nil {
			log.Fatalf("Failed to export table %s: %v", mapping.name, err)
		} else {
			log.Printf("Exported table %s to %s", mapping.name, inputFile)
		}
	}
}
