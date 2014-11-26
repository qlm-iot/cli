package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHttpServerConnection(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "REQUEST", r.FormValue("msg"))
		fmt.Fprint(w, "RESPONSE")
	}))
	defer ts.Close()

	send := make(chan []byte)
	recv := make(chan []byte)

	if assert.Nil(t, createServerConnection(ts.URL, &send, &recv)) {
		send <- []byte("REQUEST")
		response := <-recv
		assert.Equal(t, "RESPONSE", string(response))
	}
}
