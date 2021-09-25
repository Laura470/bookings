package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/Laura470/bookings/internal/config"
	"github.com/Laura470/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

var pathToTemplates = "./templates"

// NewTemplates set the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

//AddDefaultData aggiunge data a ogni pagina, lo uso per il tocken csrf
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

//REnderTEmplate renders templates using html/templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	//posso scegliere se usare la cache o no (intanto che sviluppo non la uso, così vedo subito le modifiche)
	if app.UseCache {
		// get the template cach from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		//log.Fatal("could not get template from template cache")
		return errors.New("could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}
	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	//creo una cache dove trovare i miei template pronti per l'uso
	//la key è il nome della pagina, il value è un puntatore alla stessa
	myCache := map[string]*template.Template{}

	//vado a costruire quello che la cache deve contenere

	// *. qualsiasi cosa ci sia prima del .page.tmpl
	// creo una tabella di path alle pagine ?
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	//faccio scorrere la tabella dei path delle pagine
	for _, page := range pages {
		name := filepath.Base(page)

		//ora creo un template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
