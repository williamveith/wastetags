package main

import (
	"embed"
	"encoding/json"
	"flag"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/williamveith/wastetags/pkg/database"
	"github.com/williamveith/wastetags/pkg/errors"
	"github.com/williamveith/wastetags/pkg/idgen"

	"github.com/gin-gonic/gin"
)

//go:embed assets/css/*.css assets/js/*.js assets/images/*.png
//go:embed configs/*.json
//go:embed data/*.bin
//go:embed query/*.sql
//go:embed templates/*.html
var embeddedFS embed.FS

type Config struct {
	DatabasePath string `json:"database_path"`
}

var (
	db          *database.Database
	cfg         *Config
	idGenerator = idgen.GenerateID()
)

func loadConfig() *Config {
	configPath := flag.String("config", "", "Path to the config file")
	devMode := flag.Bool("dev", false, "Run in dev mode")
	flag.Parse()

	var configs []byte
	if *devMode {
		configs = readEmbeddedFile(filepath.Join("configs", "dev.json"))
	} else if *configPath == "" {
		configs = readEmbeddedFile(filepath.Join("configs", runtime.GOOS) + ".json")
	} else {
		configs = errors.Must(os.ReadFile(*configPath))
	}

	var cfg Config
	if err := json.Unmarshal(configs, &cfg); err != nil {
		log.Fatalf("could not parse config file: %v", err)
	}

	return &cfg
}

func readEmbeddedFile(filePath string) []byte {
	return errors.Must(embeddedFS.ReadFile(filePath))
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

func init() {
	// Load configuration
	cfg = loadConfig()

	// Initialize database
	sqlStatement := readEmbeddedFile("query/schema.sql")
	db = database.NewDatabase(cfg.DatabasePath, sqlStatement)
	if db.NeedsInitialization {
		dataFS := errors.Must(fs.Sub(embeddedFS, "data"))
		db.ImportFromProtobuff(dataFS)
	}
}

func main() {
	// Engine
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Embedded Assets
	tmpl, _ := template.ParseFS(embeddedFS, "templates/*.html")
	r.SetHTMLTemplate(tmpl)
	staticFS := errors.Must(fs.Sub(embeddedFS, "assets"))
	r.StaticFS("/static", http.FS(staticFS))

	// Middleware
	r.Use(addCurrentPathMiddleware())

	// Redirects
	r.GET("/", redirectHandler("/home"))

	// Page Routes
	r.GET("/home", pageHandler(HomePage))
	r.GET("/waste-tag-form", pageHandler(MakeWasteTagForm))
	r.POST("/waste-tag", pageHandler(MakeWasteTag))
	r.GET("/add-chemical", pageHandler(AddChemical))
	r.POST("/add-chemical", pageHandler(AddChemical))
	r.GET("/add-mixture", pageHandler(AddMixture))
	r.POST("/add-mixture", pageHandler(AddMixture))

	// API Routes
	r.POST("/api/generate-qr-code", MakeNewQRCode)

	// Start HTTP Server
	r.Run(":8080")
}
