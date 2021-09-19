package models

//nel caso io debba mandare dei dati attraverso i miei handler delle pagine
//quanti dati, che dipo di dati?
// TemplateData holds data sent from handlers to template
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string //per cross site request forgery token
	Flash     string //flash message
	Warning   string
	Error     string
}
