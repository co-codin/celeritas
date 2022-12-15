package celeritas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/co-codin/celeritas/render"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Celeritas struct {
	AppName string
	Debug   bool
	Version string
	ErrorLog *log.Logger
	InfoLog *log.Logger
	RootPath string
	Routes *chi.Mux
	Render *render.Render
	config config
}

type config struct {
	port string
	renderer string
}

func (c *Celeritas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath: rootPath,
		folderNames: []string{
			"handlers",
			"migrations",
			"views",
			"data",
			"public",
			"tmp",
		},
	}
	err := c.Init(pathConfig)

	if err != nil {
		return err
	}

	err = c.checkDotEnv(rootPath)

	if err!= nil {
        return err
    }

	err = godotenv.Load(rootPath + "/.env")

	if err != nil {
		return err
	}

	infoLog, errorLog := c.startLoggers()
	
	c.InfoLog = infoLog
    c.ErrorLog = errorLog
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = version
	c.RootPath = rootPath
	c.Routes = c.routes().(*chi.Mux)

	c.config = config{
		port: os.Getenv("PORT"),
        renderer: os.Getenv("RENDERER"),
	}

	c.Render = c.createRenderer(c)

	return nil
}

func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it doesnt exist
		err := c.CreateDirIfNotExists(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Celeritas) ListenAndServe() {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog: c.ErrorLog,
		Handler: c.routes(),
		IdleTimeout: 30 * time.Second,
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	c.InfoLog.Println("Listening on port %s", os.Getenv("PORT"))

	err := srv.ListenAndServe()

	c.ErrorLog.Fatal(err)
}

func (c *Celeritas) checkDotEnv(path string) error {
	err := c.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))

	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (c *Celeritas) createRenderer(cel *Celeritas) *render.Render {
	myRenderer := render.Render{
		Renderer: cel.config.renderer,
		RootPath: cel.RootPath,
		Port: cel.config.port,
	}
	return &myRenderer
}