package models

import "github.com/Laura470/bookings/internal/forms"

//nel caso io debba mandare dei dati attraverso i miei handler delle pagine
//quanti dati, che dipo di dati?
// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
