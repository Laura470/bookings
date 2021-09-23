package forms

import (
	"net/http"
	"net/url"
)

//Form creates a custom form struct, embeds an url.values object
type Form struct {
	url.Values
	Errors errors
}

//Valid returns true if there aren't errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

//New initailises the form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

//Has checks id form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}
