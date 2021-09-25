package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}

}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	//testo valori richiesti vuota mi deve dare non valida, se dà valida c'è un errore
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form valida ma i campi richiesti sono vuoti")
	}

	//testo valori richiesti con alcuni valori, se dà non valida c'è un errore
	//faccio riferimento alla struct values, come definita in forms
	//type Values map[string][]string
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	//non capisco questo passaggio provo a escluderlo
	//r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	//chiamo la funzione
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("dice che ci sono errori ma non è vero")
	}

}

func TestForm_Has(t *testing.T) {

	form := New(url.Values{})

	has := form.Has("whatever")

	if has {
		t.Error("dice che il campo whatever è compilato ma non è vero")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	// devo reinizializzare la form con le nuove request
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("dice che il campo a non è compilato ma non è vero")
	}

	//aggiunta mia
	postedData = url.Values{}
	postedData.Add("b", "")
	// devo reinizializzare la form con le nuove request
	form = New(postedData)
	has = form.Has("b")
	if has {
		t.Error("dice che il campo b è compilato ma non è vero")
	}

}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}
	// controllo la funzione get, per avere 100% di coverage
	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have error but did not get one")
	}

	//reinizializzo per togliere errore nello slice e aggiunere i dati
	postedData = url.Values{}
	postedData.Add("some_field", "some_value")

	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows min length of 100 met when data is shorter")
	}

	//reinizializzo per togliere errore nello slice e aggiunere i dati
	postedData = url.Values{}
	postedData.Add("another_field", "abc123")
	form = New(postedData)

	//controllo la funzione
	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows min length if 1 is not met when it is")
	}
	// controllo la funzione get, per avere 100% di coverage
	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have error but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "me@here.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email address")
	}
}

/*
Mie func


// test fatti da me

func TestForm_MinLength(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever", nil)
	var length = 3

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "bbbbbbbbbbb")

	r.PostForm = postedData
	form := New(r.PostForm)

	if form.MinLength("a", length) {
		t.Error(" lunghezza data per valida ma minore di 3")
	}

	if !form.MinLength("b", length) {
		t.Error(" lunghezza data per non valida ma maggiore di 3")
	}

}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.IsEmail("a")
	if form.Valid() {
		t.Error("form valida ma in realtà è vuota")
	}

	postedData := url.Values{}
	postedData.Add("b", "ggg@jouih.ue")

	// devo reinizializzare la form per togliere il messaggio di errore inserito con il test precedente
	form = New(postedData)

	form.IsEmail("b")
	if !form.Valid() {
		t.Error("form non valida ma la stringa  è un indirizzo")
	}

}



*/
