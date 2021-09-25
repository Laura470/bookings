package main

import (
	"net/http"
	"os"
	"testing"
)

//dice prima di fare il test fai qualcosa epoi esegui il test epoi esci
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
