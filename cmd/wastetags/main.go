package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

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

type Config struct {
	DatabasePath string `json:"database_path"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return &cfg, nil
}

func readSql(filePath string) []byte {
	schema, schemaerror := sqlFS.ReadFile(filePath)
	if schemaerror != nil {
		fmt.Println("Failed to read embedded schema:", schemaerror)
		return nil
	}
	return schema
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
	// Add a command-line flag for the config file.
	configPath := flag.String("config", "/etc/wastetags/config.json", "Path to the config file")
	flag.Parse()

	// Load configuration
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	sqlStatement := readSql("query/schema.sql")
	if sqlStatement == nil {
		log.Fatalf("Failed to read schema.sql, cannot initialize database")
	}

	db = database.NewDatabase(cfg.DatabasePath, sqlStatement)

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
