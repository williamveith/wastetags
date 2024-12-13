package main

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/williamveith/wastetags/pkg/database"
	"github.com/williamveith/wastetags/pkg/qrcodegen"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:embed templates/*
var embeddedTemplatesFS embed.FS

//go:embed query/*
var sqlFS embed.FS

var db *database.Database

func readSql(filePath string) []byte {
	schema, schemaerror := sqlFS.ReadFile(filePath)
	if schemaerror != nil {
		fmt.Println("Failed to read embedded schema:", schemaerror)
		return nil
	}
	return schema
}

func init() {
	dbName := "data/chemicals.sqlite3"
	sqlStatement := readSql("query/schema.sql")
	if sqlStatement == nil {
		log.Fatalf("Failed to read schema.sql, cannot initialize database")
	}
	db = database.NewDatabase(dbName, sqlStatement)
}

func generateQRCodeBase64(dataDict map[string]interface{}) string {
	jsonContent, _ := json.Marshal(dataDict)
	dataURI := qrcodegen.DataURI(string(jsonContent), nil)
	return dataURI
}

func generateShortUUID(length int) string {
	rawUUID := uuid.New()
	base64Encoded := base64.StdEncoding.EncodeToString(rawUUID[:])
	base64Clean := strings.NewReplacer("+", "", "/", "", "=", "").Replace(base64Encoded)
	return base64Clean[:length]
}

func completedWasteTag(c *gin.Context) {
	values := map[string]string{
		"chemName":  c.PostForm("chemName"),
		"location":  c.PostForm("location"),
		"contCount": c.PostForm("contCount"),
		"contSize":  c.PostForm("contSize"),
		"sizeUnit":  c.PostForm("sizeUnit"),
		"contType":  c.PostForm("contType"),
		"quantity":  c.PostForm("quantity"),
		"unit":      c.PostForm("unit"),
		"physState": c.PostForm("physState"),
	}

	wasteTag := generateShortUUID(20)
	sqlStatement := readSql("query/get_rows_by_col_value.sql")
	components, err := db.GetRowsByColumnValue("mixtures", sqlStatement, "chem_name", values["chemName"])
	if err != nil {
		fmt.Println("Error retrieving components:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Internal Server Error"})
		return
	}

	qrCodeBase64 := generateQRCodeBase64(map[string]interface{}{
		"wasteTags":  []string{wasteTag},
		"values":     values,
		"components": components,
	})

	componentData := make([]map[string]string, len(components))
	for i, component := range components {
		componentData[i] = map[string]string{
			"Chemical":   fmt.Sprint(component["component_name"]),
			"Percentage": fmt.Sprint(component["percent"]),
		}
	}

	templateData := map[string]interface{}{
		"RoomNumber":    values["location"],
		"TagNumber":     wasteTag,
		"GeneratorName": "William Veith",
		"QrCodeBase64":  template.URL(qrCodeBase64),
		"components":    componentData,
	}

	c.HTML(http.StatusOK, "tag.html", templateData)
}

func wasteLabelForm(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		completedWasteTag(c)
		return
	}

	sqlStatement := readSql("query/get_distinct_col_values.sql")
	components, err := db.GetColumnValues("mixtures", sqlStatement, "chem_name")
	if err != nil {
		fmt.Println("Error retrieving components:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Internal Server Error"})
		return
	}
	cleanComponents := make([]map[string]string, len(components))
	for i, comp := range components {
		if chemName, ok := comp["chem_name"].(string); ok {
			cleanComponents[i] = map[string]string{"chem_name": chemName}
		}
	}

	locations, err := db.GetColumnValues("locations", sqlStatement, "location")
	if err != nil {
		fmt.Println("Error retrieving locations:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Internal Server Error"})
		return
	}

	cleanLocations := make([]map[string]string, len(locations))
	for i, comp := range locations {
		if location, ok := comp["location"].(string); ok {
			cleanLocations[i] = map[string]string{"location": location}
		}
	}

	units, err := db.GetColumnValues("units", sqlStatement, "full_name")
	if err != nil {
		fmt.Println("Error retrieving locations:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Internal Server Error"})
		return
	}

	cleanUnits := make([]map[string]string, len(units))
	for i, comp := range units {
		if unit, ok := comp["full_name"].(string); ok {
			cleanUnits[i] = map[string]string{"full_name": unit}
		}
	}

	containers, err := db.GetColumnValues("containers", sqlStatement, "full_name")
	if err != nil {
		fmt.Println("Error retrieving locations:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Internal Server Error"})
		return
	}

	cleanContainers := make([]map[string]string, len(containers))
	for i, comp := range containers {
		if container, ok := comp["full_name"].(string); ok {
			cleanContainers[i] = map[string]string{"full_name": container}
		}
	}

	states, err := db.GetColumnValues("states", sqlStatement, "state")
	if err != nil {
		fmt.Println("Error retrieving locations:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Internal Server Error"})
		return
	}

	cleanStates := make([]map[string]string, len(states))
	for i, comp := range states {
		if state, ok := comp["state"].(string); ok {
			cleanStates[i] = map[string]string{"state": state}
		}
	}

	c.HTML(http.StatusOK, "tag-form.html", gin.H{
		"Components": cleanComponents,
		"Locations":  cleanLocations,
		"Units":      cleanUnits,
		"Containers": cleanContainers,
		"States":     cleanStates,
	})

}

func runLabelMaker() {
	gin.SetMode(gin.ReleaseMode)

	tmpl, _ := template.ParseFS(embeddedTemplatesFS, "templates/*")

	r := gin.Default()
	r.SetHTMLTemplate(tmpl)

	r.GET("/", wasteLabelForm)
	r.POST("/", wasteLabelForm)
	r.Run(":8080")
}

func main() {
	runLabelMaker()
}
