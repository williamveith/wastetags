package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/williamveith/wastetags/pkg/database"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
//go:exclude templates/.* templates/.*/**
var embeddedTemplatesFS embed.FS

//go:embed assets/*
//go:exclude assets/.* assets/.*/**
var embeddedStylesFS embed.FS

//go:embed query/*
var sqlFS embed.FS

var db *database.Database

func homePage(c *gin.Context) (string, gin.H) {
	return "home.html", nil
}

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

func pageHandler(handler func(c *gin.Context) (string, gin.H)) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.Request.URL.Path
		templateName, data := handler(c)

		if data == nil {
			data = gin.H{}
		}

		data["CurrentPath"] = currentPath
		c.HTML(http.StatusOK, templateName, data)
	}
}

func addCurrentPathMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("CurrentPath", c.Request.URL.Path)
		c.Next()
	}
}

func runLabelMaker() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.StaticFS("/static", http.FS(embeddedStylesFS))

	tmpl, _ := template.ParseFS(embeddedTemplatesFS, "templates/*")
	r.SetHTMLTemplate(tmpl)

	r.Use(addCurrentPathMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})

	r.POST("/api/generate-qr", func(c *gin.Context) {
		var requestData map[string]interface{}

		// Parse the JSON request body
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Call the Go function to generate QR code
		qrCodeData, err := convertMapToQRCodeData(requestData)
		dataURI, jsonContent, wasteTag := qrCodeData.makeCopy()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Respond with the result
		c.JSON(http.StatusOK, gin.H{
			"dataURI":     dataURI,
			"jsonContent": jsonContent,
			"wasteTag":    wasteTag,
		})
	})

	r.GET("/home", pageHandler(homePage))
	r.GET("/create-tag", pageHandler(makeWasteTagForm))
	r.POST("/wastetag", pageHandler(makeWasteTag))
	r.GET("/addchemical", pageHandler(addChemical))
	r.POST("/addchemical", pageHandler(addChemical))
	r.GET("/add-mixture", pageHandler(addMixture))
	r.POST("/add-mixture", pageHandler(addMixture))
	r.Run(":8080")
}

func main() {
	runLabelMaker()
}
