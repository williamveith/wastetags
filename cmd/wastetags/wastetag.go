package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/williamveith/wastetags/pkg/qrcodegen"
)

type TagTemplateData struct {
	RoomNumber    string
	TagNumber     string
	GeneratorName string
	QrCodeBase64  template.URL
	QRCodeValue   string
	Components    []Component
}

type Component struct {
	Chemical   string
	Percentage string
}

type QRCodeData struct {
	Version    string            `json:"version"`
	WasteTags  []string          `json:"name"`
	Values     map[string]string `json:"values"`
	Components []Component       `json:"components"`
}

func (qrCodeData *QRCodeData) generateDataUri() (string, string) {
	jsonContentBytes, _ := json.Marshal(qrCodeData)
	jsonContent := string(jsonContentBytes)
	dataURI := qrcodegen.DataURI(jsonContent, nil)
	return dataURI, jsonContent
}

func (qrCodeData *QRCodeData) makeCopy() (string, string, string) {
	qrCodeData.WasteTags[0] = generateShortUUID(len(qrCodeData.WasteTags[0]))
	dataURI, jsonContent := qrCodeData.generateDataUri()
	return dataURI, jsonContent, qrCodeData.WasteTags[0]
}

func convertMapToQRCodeData(dataMap map[string]interface{}) (*QRCodeData, error) {
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	var qrCodeData QRCodeData
	if err := json.Unmarshal(jsonData, &qrCodeData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to QRCodeData: %w", err)
	}

	return &qrCodeData, nil
}

func generateShortUUID(length int) string {
	rawUUID := uuid.New()
	base64Encoded := base64.StdEncoding.EncodeToString(rawUUID[:])
	base64Clean := strings.NewReplacer("+", "", "/", "", "=", "").Replace(base64Encoded)
	return base64Clean[:length]
}

func completedWasteTag(c *gin.Context) (string, gin.H) {
	genericErrorMessage := gin.H{"message": "Internal Server Error"}

	if c.Request.Method != http.MethodPost {
		return "index.html", nil
	}

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
		return "error.html", genericErrorMessage
	}

	componentData := make([]Component, len(components))
	for i, component := range components {
		componentData[i] = Component{
			Chemical:   fmt.Sprint(component["component_name"]),
			Percentage: fmt.Sprint(component["percent"]),
		}
	}

	qrCodeData := &QRCodeData{
		Version:    "v1",
		WasteTags:  []string{wasteTag},
		Values:     values,
		Components: componentData,
	}

	qrCodeBase64, jsonContent := qrCodeData.generateDataUri()

	templateData := map[string]interface{}{
		"RoomNumber":    values["location"],
		"TagNumber":     wasteTag,
		"GeneratorName": "William Veith",
		"QRCodeValue":   jsonContent,
		"QrCodeBase64":  template.URL(qrCodeBase64),
		"Components":    componentData,
	}

	return "tag.html", templateData
}
