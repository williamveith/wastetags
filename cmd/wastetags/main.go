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

//go:embed templates/*.html
var embeddedTemplatesFS embed.FS

//go:embed assets/*
//go:exclude assets/.* assets/.*/**
var embeddedStylesFS embed.FS

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

func redirectHandler(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, path)
	}
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

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	tmpl, _ := template.ParseFS(embeddedTemplatesFS, "templates/*")
	r.SetHTMLTemplate(tmpl)
	r.StaticFS("/static", http.FS(embeddedStylesFS))

	r.Use(addCurrentPathMiddleware())

	r.GET("/", redirectHandler("/home"))
	r.GET("/home", pageHandler(HomePage))
	r.GET("/waste-tag-form", pageHandler(MakeWasteTagForm))
	r.POST("/waste-tag", pageHandler(MakeWasteTag))
	r.GET("/add-chemical", pageHandler(AddChemical))
	r.POST("/add-chemical", pageHandler(AddChemical))
	r.GET("/add-mixture", pageHandler(AddMixture))
	r.POST("/add-mixture", pageHandler(AddMixture))
	r.POST("/api/generate-qrcode", MakeNewQRCode)
	r.Run(":8080")
}
