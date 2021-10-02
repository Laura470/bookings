package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Laura470/bookings/internal/models"
)

/* type postData struct {
	key   string
	value string
} */

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	//inserisco gli elementi della struct, che sono slice
	{"home", "/", "GET", http.StatusOK}, //constant è uguale a 200
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availibility", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	//non poso testare le funzioni che hanno bisogno diuna session
	/* 	{"post-search-availibility", "/search-availibility", "Post", []postData{
	   		{key: "start", value: "2020-01-01"},
	   		{key: "end", value: "2020-01-02"},
	   	}, http.StatusOK},
	   	{"post-search-availibility-json", "/search-availibility-json", "Post", []postData{
	   		{key: "start", value: "2020-01-01"},
	   		{key: "end", value: "2020-01-02"},
	   	}, http.StatusOK},
	   	{"make-reservation", "/make-reservation", "Post", []postData{
	   		{key: "first_name", value: "John"},
	   		{key: "last_name", value: "Smith"},
	   		{key: "email", value: "me@here.com"},
	   		{key: "phone", value: "555-555-5555"},
	   	}, http.StatusOK}, */
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//devo crere un sever e un client che chiama il server, ma in go è già tutto creato
	ts := httptest.NewTLSServer(routes) //ts è il mio testserver
	defer ts.Close()                    //sempre meglio chiudere

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url) //creo il client e gli aggiungo l'url che voglio testare
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	//faccio una richiesta con un empty body
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	//uso il package di go per creare una risposta simulata di qualcuno
	//che ha visitato il nostro sito
	rr := httptest.NewRecorder()
	//metto la mia reservation nella session
	session.Put(ctx, "reservation", reservation)

	//taking my reservation and casting it
	//in a function that I can call
	//cast my reservation handler in to a handler func
	//no routes necessary
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//test case where reservation is not in session (reset everything)
	//reinizializzo req
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test with not existing room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {

	// ---------------------- 1° TEST ----------------------------------------------
	//test iwth everything ok
	//now I build the body request
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=jj@jj.it")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//in questo caso non posso fare una richiesta con un empty body, è un post!
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	//setting the header for the request che avvisa che è un POST
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	//metto lo status che si aspetta l'handler (StatusSeeOther)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// ---------------------- 2° TEST ----------------------------------------------
	//test for missing request body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	//testo il messaggio di errore, c'è un errore e il messaggio di errore deve essere quello giusto (non testo se rileva o no l'errore, testo il risultato della rilevazione)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// ---------------------- 3° TEST ----------------------------------------------

	// test for invalid start date , it will pass the first two test, but fail date parsing
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// ---------------------- 4° TEST ----------------------------------------------
	//in alternativa:
	/*
			postedData := url.Values{}
		postedData.Add("start_date", "2050-01-01")
		postedData.Add("end_date", "2050-01-02")
		postedData.Add("first_name", "John")
		postedData.Add("last_name", "Smith")
		postedData.Add("email", "john@smith.com")
		postedData.Add("phone", "555-555-5555")
		postedData.Add("room_id", "1")

		req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	*/

	// test for invalid end date
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// ---------------------- 5° TEST ----------------------------------------------
	// test for invalid room id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid room id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// ---------------------- 6° TEST ----------------------------------------------
	// test for invalid data first_name<3
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// ---------------------- 7° TEST ----------------------------------------------
	// test for failure to insert reservation into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// ---------------------- 8° TEST ----------------------------------------------
	// test for failure to insert restriction into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostAvailibility(t *testing.T) {

	// non si riesce a fare parsing di start
	reqBody := "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")

	//in questo caso non posso fare una richiesta con un empty body, è un post!
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	//setting the header for the request che avvisa che è un POST
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailibility)

	handler.ServeHTTP(rr, req)
	//metto lo status che si aspetta l'handler (StatusSeeOther)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Availibility with invalid start date and wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// non si riesce a fare parsing  di end
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")

	//in questo caso non posso fare una richiesta con un empty body, è un post!
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//setting the header for the request che avvisa che è un POST
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailibility)

	handler.ServeHTTP(rr, req)
	//metto lo status che si aspetta l'handler (StatusSeeOther)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Availibility with invalid end date and  wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// non si riesce a fare parsing della form
	//in questo caso non posso fare una richiesta con un empty body, è un post!
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//setting the header for the request che avvisa che è un POST
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailibility)

	handler.ServeHTTP(rr, req)
	//metto lo status che si aspetta l'handler (StatusSeeOther)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Availibility with body form empty and  wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//errore di connessione con il db
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")

	//in questo caso non posso fare una richiesta con un empty body, è un post!
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//setting the header for the request che avvisa che è un POST
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailibility)

	handler.ServeHTTP(rr, req)
	//metto lo status che si aspetta l'handler (StatusSeeOther)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post Availibility with connection error with the database and wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//len == 0  - rooms are not available
	/*****************************************/
	// create our request body
	reqBody = "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")

	//in questo caso non posso fare una richiesta con un empty body, è un post!
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailibility)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	//metto lo status che si aspetta l'handler (StatusSeeOther)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Availibility says room available when no rooms are available and wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//len>0
	reqBody = "start=2040-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2040-01-02")

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.PostAvailibility)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	//ATTENZIONE!!!!!!!!!!!!!
	// since we have rooms available, we expect to get status http.StatusOK
	if rr.Code != http.StatusOK {
		t.Errorf("Post availability when rooms are available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {

	/*****************************************
	// first case -- rooms are not available
	*****************************************/
	// create our request body
	reqBody := "start=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create our request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.AvailibilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	//come faccio ad avere no room available?????
	// this time we want to parse JSON and get the expected response
	var j jsonResponse
	//err := json.Unmarshal([]byte(rr.Body.String()), &j)
	err := json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date > 2049-12-31, we expect no availability
	//ha modificato la funzione SearchAvailabilityByDatesByRoomID in testrepo
	if j.OK {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}

	/*****************************************
	// Second case -- can't parse form,
	*****************************************/

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailibilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	// this time we want to parse JSON and get the expected response
	//err := json.Unmarshal([]byte(rr.Body.String()), &j)
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse form!")
	}

	// since we specified a start date > 2049-12-31, we expect no availability
	if j.OK {
		t.Error("Got availability when none was expected in AvailabilityJSON")
	}
	/*****************************************
	// third case -- got a database error
	*****************************************/
	//modifico ulteriormente la funzione SearchAvailabilityByDatesByRoomID in test tepo
	// create our request body
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// create our request
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set the request header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()

	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.AvailibilityJSON)

	// make the request to our handler
	handler.ServeHTTP(rr, req)

	// since we have no rooms available, we expect to get status http.StatusSeeOther
	//come faccio ad avere no room available?????
	// this time we want to parse JSON and get the expected response

	//err := json.Unmarshal([]byte(rr.Body.String()), &j)
	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// since we specified a start date = 2060.01.01, we expect error
	//ha modificato la funzione SearchAvailabilityByDatesByRoomID in testrepo
	if j.OK {
		t.Error("Got availability when an error was expected in AvailabilityJSON")
	}

}

func TestRepository_ReservationSummary(t *testing.T) {

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	//faccio una richiesta E' una get, non  una post, quindi non mi serve il bogy!!!!
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	//Ricordati di cambiare l'handler!!!!!!!
	handler := http.HandlerFunc(Repo.ReservationSummary)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//reservation not in session
	//faccio una richiesta E' una get, non  una post, quindi non mi serve il bogy!!!!
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	//uso il package di go per creare una risposta simulata di qualcuno
	//che ha visitato il nostro sito
	rr = httptest.NewRecorder()

	//non metto la reservation nella session

	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	/*
		In your test for ChooseRoom, you will want to set the URL on your request as follows:
		req.RequestURI = "/choose-room/1"
	*/
	/* 	reservation := models.Reservation{
	   		RoomID: 1,
	   		Room: models.Room{
	   			ID:       1,
	   			RoomName: "General's Quarters",
	   		},
	   	}
	   	//faccio una richiesta E' una get, non  una post, quindi non mi serve il bogy!!!!
	   	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	   	ctx := getCtx(req)
	   	req = req.WithContext(ctx)

	   	// set the RequestURI on the request so that we can grab the ID
	   	// from the URL
	   	req.RequestURI = "/choose-room/1"

	   	rr := httptest.NewRecorder()
	   	session.Put(ctx, "reservation", reservation)

	   	//Ricordati di cambiare l'handler!!!!!!!
	   	handler := http.HandlerFunc(Repo.ChooseRoom)

	   	handler.ServeHTTP(rr, req)

	   	if rr.Code != http.StatusOK {
	   		t.Errorf("Choosen room handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	   	} */

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/choose-room/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//secondo test missing url parameter
	/* 	req, _ = http.NewRequest("GET", "/choose-room/cat", nil)
	   	ctx = getCtx(req)
	   	req = req.WithContext(ctx)
	   	req.RequestURI = "/choose-room/cat"

	   	rr = httptest.NewRecorder()

	   	//Ricordati di cambiare l'handler!!!!!!!
	   	handler = http.HandlerFunc(Repo.ChooseRoom)

	   	handler.ServeHTTP(rr, req)

	   	if rr.Code != http.StatusTemporaryRedirect {
	   		t.Errorf("Choosen romm handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	   	}
	*/
	req, _ = http.NewRequest("GET", "/choose-room/fish", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/fish"

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//terzo test Can't get reservation from session
	/* 	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	   	ctx = getCtx(req)
	   	req = req.WithContext(ctx)

	   	rr = httptest.NewRecorder()

	   	//Ricordati di cambiare l'handler!!!!!!!
	   	handler = http.HandlerFunc(Repo.ChooseRoom)

	   	handler.ServeHTTP(rr, req)

	   	if rr.Code != http.StatusTemporaryRedirect {
	   		t.Errorf("choosen room handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	   	} */

	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_BookRoom(t *testing.T) {

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	//testo se tutto va bene

	// create our request  E' una get!!!!!
	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	// get the context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	//nella get non serve la header?
	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	// make our handler a http.HandlerFunc
	handler := http.HandlerFunc(Repo.BookRoom)
	// make the request to our handler
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//testo errore connessione db con id>2

	// create our request
	req, _ = http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=5", nil)
	// get the context with session
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// create our response recorder, which satisfies the requirements
	// for http.ResponseWriter
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	// make our handler a http.HandlerFunc
	handler = http.HandlerFunc(Repo.BookRoom)
	// make the request to our handler
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

//
//
//la sessione deve avere un context
//x-session è la chiave per leggere la sessione
func getCtx(req *http.Request) context.Context {

	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
