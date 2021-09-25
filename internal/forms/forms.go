package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
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

//variatic function ??? I may have how many fields I want ???
//Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is cannot be blank")
		}
	}
}

//Has checks id form field is in post and not empty
//has guarda solo se esiste la key, on oil value a quanto ho capito
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

//Min length checks for strin g minimumlength
func (f *Form) MinLength(field string, length int) bool {
	//x := r.Form.Get(field)
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("THis field must be al least %d characters long", length))
		return false
	}
	return true
}

//IsEmail checks for vaild email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
