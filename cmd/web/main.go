package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Laura470/bookings/internal/config"
	"github.com/Laura470/bookings/internal/handlers"
	"github.com/Laura470/bookings/internal/models"
	"github.com/Laura470/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() error {

	//what i'm going to put in the session
	gob.Register(models.Reservation{})

	//change this to true when in production
	//in here so it is available outside the main for the main package (middleware is in the main package)
	app.InProduction = false

	session = scs.New() //tolto il due punti
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //anche dopo il browser è chiuso
	// abbiamo fatto la stessa cosa in middleware con NoSurf package
	session.Cookie.SameSite = http.SameSiteLaxMode //  quanto tight ??
	session.Cookie.Secure = app.InProduction       //in production non encripted

	//inizializzo la session in config
	app.Session = session

	//chiamo la funzione CreateTemplateCache dal package render
	tc, err := render.CreateTemplateCache()
	if err != nil {
		//log.Fatal("cannot create template cache")
		log.Fatal(err)
		return err
	}

	//prende il suo valore da render, fare attenzione all'import
	app.TemplateCache = tc
	//setto la variabile a false
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	return nil
}
