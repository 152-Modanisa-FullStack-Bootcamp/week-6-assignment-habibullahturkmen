package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	// Created new store and service
	store := NewInMemoryWallet()
	server := WalletServer{store}

	// Created three users
	user1 := "Recep"
	user2 := "Mehmet"
	user3 := "Davut"

	// Calls PUT method three times
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user1))
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user2))
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user3))

	response := httptest.NewRecorder() // spy
	server.ServeHTTP(response, newGetWalletsRequest("")) // server
	assertStatus(t, response.Code, http.StatusOK) // asserts status code

	want := fmt.Sprintf("{\"%s\":%d,\"%s\":%d,\"%s\":%d}", user3, 0, user2, 0, user1, 0)

	assertResponseBody(t, response.Body.String(), want) // asserts response body
}
