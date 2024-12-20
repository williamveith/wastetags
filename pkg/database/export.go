package database

import (
	"fmt"
	"log"
	"os"
	reflect "reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/proto"
)

func dbColumnNamesToProtobufKey(columnNames []string) []string {
	titleCase := cases.Title(language.English)
	for i, columnName := range columnNames {
		columnName = strings.ReplaceAll(columnName, "_", " ")
		columnNames[i] = strings.ReplaceAll(titleCase.String(columnName), " ", "")
	}
	return columnNames
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
