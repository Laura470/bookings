package handlers

import (
	"net/http"

	"github.com/Laura470/bookings/pkg/config"
	"github.com/Laura470/bookings/pkg/models"
	"github.com/Laura470/bookings/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
// aggiunto a receiver alla funzione
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	//Grab the ip address of the person visiting my site and store it in homepage
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
// aggiunto a receiver alla funzione
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some business logic
	//devo mandare dei dati, mi serve un terzo argomento, ma di che tipo? lo creo sopra
	stringMap := make(map[string]string)
	stringMap["test"] = "hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
