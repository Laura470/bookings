package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Laura470/bookings/internal/config"
	"github.com/Laura470/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	//what i'm going to put in the session
	gob.Register(models.Reservation{})

	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true

	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false //beacuse we are not running an https server

	testApp.Session = session

	app = &testApp //app dichiarata in render.go

	os.Exit(m.Run())
}

//vedo che cosa ha il type ResponseWrite interface: un header, un write([]byte) e un writeheader

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil

}
