package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Laura470/bookings/internal/config"
	"github.com/Laura470/bookings/internal/driver"
	"github.com/Laura470/bookings/internal/handlers"
	"github.com/Laura470/bookings/internal/helpers"
	"github.com/Laura470/bookings/internal/models"
	"github.com/Laura470/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	/*
		//setting an email to be send whne the app startswith standar library
		from := "me@here.com"
		auth := smtp.PlainAuth("", from, "", "localhost")
		err = smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("Hello world"))
		if err != nil {
			log.Println(err)
		}*/

	defer close(app.MailChan)
	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() (*driver.DB, error) {

	//what i'm going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//change this to true when in production
	//in here so it is available outside the main for the main package (middleware is in the main package)
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New() //tolto il due punti
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //anche dopo il browser è chiuso
	// abbiamo fatto la stessa cosa in middleware con NoSurf package
	session.Cookie.SameSite = http.SameSiteLaxMode //  quanto tight ??
	session.Cookie.Secure = app.InProduction       //in production non encripted

	//inizializzo la session in config
	app.Session = session

	//inizializzo il db
	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=laura")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}

	log.Println("Connected to database!")

	//chiamo la funzione CreateTemplateCache dal package render
	tc, err := render.CreateTemplateCache()
	if err != nil {
		//log.Fatal("cannot create template cache")
		log.Fatal(err)
		return nil, err
	}

	//prende il suo valore da render, fare attenzione all'import
	app.TemplateCache = tc
	//setto la variabile a false
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
