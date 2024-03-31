package main

import (
	"esolang/lang-esolang/repl"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Esolang struct {
	Input      string
	SourceCode string
}

func newEsolangConstruct(input, sc string) Esolang {
	return Esolang{
		Input:      input,
		SourceCode: sc,
	}
}

type EsoInputs = []Esolang

type EsolangData struct {
	Esos EsoInputs
}

func newEsolangData() EsolangData {
	return EsolangData{
		Esos: []Esolang{},
	}
}

func main() {
	e := echo.New()
	e.Renderer = NewTemplates()
	e.Use(middleware.Logger())

	data := newEsolangData()
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", data)
	})

	e.POST("/playground", func(c echo.Context) error {
		//  curl -X POST -F 'sourceCode=println("hi mom")' http://localhost:8080/playground
		sc := c.FormValue("sourceCode")
		playGroundRes := repl.EvlauateFromPlayground(sc)
		clear(data.Esos)
		data.Esos = append(data.Esos, newEsolangConstruct(sc, playGroundRes))
		return c.Render(http.StatusOK, "evaluatedView", data)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
