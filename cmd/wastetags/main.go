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

var chemicals *database.ChemicalDatabase

func init() {
	var err error
	chemicals, err = database.NewChemicalDatabase("data/chemicals.sqlite3")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

func generateQRCodeBase64(dataDict map[string]interface{}) string {
	jsonContent, err := json.Marshal(dataDict)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return ""
	}

	dataURI, err := qrcodegen.DataURI(string(jsonContent))
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return ""
	}

	return dataURI
}

func generateShortUUID(length int) string {
	rawUUID := uuid.New()
	base64Encoded := base64.StdEncoding.EncodeToString(rawUUID[:])
	base64Clean := strings.NewReplacer("+", "", "/", "", "=", "").Replace(base64Encoded)
	return base64Clean[:length]
}

func wasteLabelForm(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
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
		components, err := chemicals.GetRowsByName(values["chemName"])
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

		return
	}

	c.HTML(http.StatusOK, "tag-form.html", nil)
}

func runLabelMaker() {
	gin.SetMode(gin.ReleaseMode)

	tmpl, err := template.ParseFS(embeddedTemplatesFS, "templates/*")
	if err != nil {
		log.Fatalf("Failed to parse embedded templates: %v", err)
	}

	r := gin.Default()
	r.SetHTMLTemplate(tmpl)

	r.GET("/", wasteLabelForm)
	r.POST("/", wasteLabelForm)
	r.Run(":8080")
}

func main() {
	runLabelMaker()
}
