package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Laura470/bookings/internal/config"
	"github.com/Laura470/bookings/internal/driver"
	"github.com/Laura470/bookings/internal/models"
	"github.com/Laura470/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}
var db *driver.DB

func getRoutes() http.Handler {

	//what i'm going to put in the session
	gob.Register(models.Reservation{})

	//change this to true when in production
	//in here so it is available outside the main for the main package (middleware is in the main package)
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New() //tolto il due punti
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //anche dopo il browser è chiuso
	// abbiamo fatto la stessa cosa in middleware con NoSurf package
	session.Cookie.SameSite = http.SameSiteLaxMode //  quanto tight ??
	session.Cookie.Secure = app.InProduction       //in production non encripted

	//inizializzo la session in config
	app.Session = session

	//chiamo la funzione CreateTemplateCache dal package render
	tc, err := CreateTestTemplateCache()
	if err != nil {
		//log.Fatal("cannot create template cache")
		log.Fatal(err)

	}

	//prende il suo valore da render, fare attenzione all'import
	app.TemplateCache = tc
	//setto la variabile a false
	app.UseCache = true

	repo := NewRepo(&app, db)
	NewHandlers(repo)
	render.NewRenderer(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf) non uso il token perchè l'ho già testato nel middleware
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availibility", Repo.Availibility)
	mux.Post("/search-availibility", Repo.PostAvailibility)
	mux.Post("/search-availibility-json", Repo.AvailibilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	//per potere visualizzare i file statici nelle mie pagine html
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

//NoSurf adds CSRF protection to all Post request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {

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
