package database

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"strings"
	"unicode"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// protobufKeyToDbColumnNames converts Proto field names (e.g. CamelCase) to snake_case.
func protobufKeyToDbColumnNames(protobufKeys []string) []string {
	for i, protobufKey := range protobufKeys {
		var words []string
		wordStart := 0
		runes := []rune(protobufKey)
		for j := 1; j < len(runes); j++ {
			if unicode.IsUpper(runes[j]) && (j != 0) {
				words = append(words, string(runes[wordStart:j]))
				wordStart = j
			}
		}
		words = append(words, string(runes[wordStart:]))
		protobufKeys[i] = strings.ToLower(strings.Join(words, "_"))
	}
	return protobufKeys
}

func (cdb *Database) FromProtobuf(tableName string, protoItemType proto.Message, protoCollectionType proto.Message, data []byte) error {
	cdb.lock.Lock()
	defer cdb.lock.Unlock()

	// Unmarshal data into protoCollectionType
	if err := proto.Unmarshal(data, protoCollectionType); err != nil {
		return fmt.Errorf("failed to unmarshal Protobuf data: %w", err)
	}

	// Obtain message reflection of the collection message
	m := protoCollectionType.ProtoReflect()
	md := m.Descriptor()

	// Get the descriptor of the item type
	itemMsgDesc := protoItemType.ProtoReflect().Descriptor()

	// Identify the repeated field in protoCollectionType that holds the items
	var itemsField protoreflect.FieldDescriptor
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if fd.IsList() && fd.Kind() == protoreflect.MessageKind {
			// Check if this field's message type matches the item type
			if fd.Message() == itemMsgDesc {
				itemsField = fd
				break
			}
		}
	}

	if itemsField == nil {
		return fmt.Errorf("no repeated field with the desired item message type found in %s", md.FullName())
	}

	// Retrieve the list of items
	val := m.Get(itemsField)
	itemList := val.List()
	rowCount := itemList.Len()
	if rowCount == 0 {
		fmt.Printf("No data to insert into table %s\n", tableName)
		return nil
	}

	// Determine columns from the first item
	firstItemMsg := itemList.Get(0).Message()
	firstItemMd := firstItemMsg.Descriptor()

	columns := make([]string, 0, firstItemMd.Fields().Len())
	placeholders := make([]string, 0, firstItemMd.Fields().Len())
	colCount := firstItemMd.Fields().Len()

	for i := 0; i < colCount; i++ {
		fd := firstItemMd.Fields().Get(i)
		columnName := protobufKeyToDbColumnNames([]string{string(fd.Name())})[0]
		columns = append(columns, columnName)
		placeholders = append(placeholders, "?")
	}

	columnsStr := strings.Join(columns, ", ")
	rowPlaceholder := "(" + strings.Join(placeholders, ", ") + ")"
	batchSize := 1000 // Adjust as needed
	valuesPerRow := colCount
	queryPrefix := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, columnsStr)

	allValues := []interface{}{}
	numInserted := 0
	batchCount := 0

	flushBatch := func() error {
		if len(allValues) == 0 {
			return nil // Nothing to flush
		}

		// Calculate how many rows are in this batch
		rowsInBatch := len(allValues) / valuesPerRow
		rowsPlaceholders := strings.Repeat(rowPlaceholder+",", rowsInBatch)
		rowsPlaceholders = strings.TrimRight(rowsPlaceholders, ",")
		query := queryPrefix + rowsPlaceholders

		if _, err := cdb.db.Exec(query, allValues...); err != nil {
			return fmt.Errorf("failed to insert batch into table %s: %w", tableName, err)
		}

		allValues = allValues[:0] // reset slice
		batchCount++
		return nil
	}

	// Convert protoreflect values to Go types
	convertValue := func(fd protoreflect.FieldDescriptor, v protoreflect.Value) interface{} {
		switch fd.Kind() {
		case protoreflect.BoolKind:
			return v.Bool()
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return int32(v.Int())
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return v.Int()
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return uint32(v.Uint())
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return v.Uint()
		case protoreflect.FloatKind:
			return float32(v.Float())
		case protoreflect.DoubleKind:
			return v.Float()
		case protoreflect.StringKind:
			return v.String()
		case protoreflect.BytesKind:
			return v.Bytes()
		case protoreflect.EnumKind:
			return v.Enum()
		case protoreflect.MessageKind:
			// If nested message, handle accordingly (e.g., serialize to JSON)
			msg := v.Message()
			pm := msg.Interface()
			jsonBytes, err := protojson.Marshal(pm)
			if err != nil {
				return fmt.Sprintf("error marshaling nested message: %v", err)
			}
			return string(jsonBytes)
		default:
			return v.Interface()
		}
	}

	// Collect items and insert in batches
	for i := 0; i < rowCount; i++ {
		itemMsg := itemList.Get(i).Message()
		for j := 0; j < colCount; j++ {
			fd := firstItemMd.Fields().Get(j)
			v := itemMsg.Get(fd)
			goValue := convertValue(fd, v)
			allValues = append(allValues, goValue)
		}

		numInserted++
		if numInserted%batchSize == 0 {
			if err := flushBatch(); err != nil {
				return err
			}
		}
	}

	// Flush any remaining rows
	if err := flushBatch(); err != nil {
		return err
	}

	fmt.Printf("Protobuf data successfully imported into table %s, inserted %d rows in %d batches\n", tableName, rowCount, batchCount)
	return nil
}

func (cdb *Database) ImportFromProtobuff(embeddedData fs.FS) {
	// Define the tables and their corresponding Protobuf types
	// Replace Chemical, ChemicalList, etc. with your actual generated types.
	tableMappings := []struct {
		name        string
		message     proto.Message
		listMessage proto.Message
	}{
		{"alias", &Alias{}, &AliasList{}},
		{"chemicals", &Chemical{}, &ChemicalList{}},
		{"mixtures", &Mixture{}, &MixtureList{}},
		{"locations", &Location{}, &LocationList{}},
		{"containers", &Container{}, &ContainerList{}},
		{"units", &Unit{}, &UnitList{}},
		{"states", &State{}, &StateList{}},
	}

	// Loop through each table and import its data from Protobuf
	for _, mapping := range tableMappings {
		inputFile, err := embeddedData.Open(mapping.name + ".bin")
		if err != nil {
			log.Fatalf("Failed to open input file for table %s: %v", mapping.name, err)
		}

		data, err := io.ReadAll(inputFile)
		if err != nil {
			log.Fatalf("Failed to read input file for table %s: %v", mapping.name, err)
		}

		err = cdb.FromProtobuf(
			mapping.name,        // Table name
			mapping.message,     // Single-row Protobuf message
			mapping.listMessage, // List Protobuf message
			data,                // Input file data
		)
		if err != nil {
			log.Fatalf("Failed to import table %s: %v", mapping.name, err)
		} else {
			log.Printf("Imported table %s from Protobuf data", mapping.name)
		}
	}
}
