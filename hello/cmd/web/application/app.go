package application

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"grusp.io/hello/internal/models"
)

type Application struct {
	TemplateCache map[string]*template.Template
}

func NewApplication(templateCache map[string]*template.Template) *Application {
	return &Application{
		TemplateCache: templateCache,
	}
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	fmt.Println(err)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) Render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	ts, ok := app.TemplateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
