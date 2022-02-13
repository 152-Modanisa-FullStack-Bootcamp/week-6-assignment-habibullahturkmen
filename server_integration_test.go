package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test the whole application
func TestTheWholeApplication(t *testing.T) {
	// Created new store and service
	store := NewInMemoryWallet()
	server := WalletServer{store}

	// Created three users
	user1 := "Recep"
	user2 := "Mehmet"
	user3 := "Davut"

	// Creates
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user1))
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user2))

	// Tries to create the same user twice
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user3))
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user3))

	// Tries to create a user with empty string
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(""))

	// Calls PUT method three times
	server.ServeHTTP(httptest.NewRecorder(), newPostWalletRequest(user1, 420))
	server.ServeHTTP(httptest.NewRecorder(), newPostWalletRequest(user2, -200))
	server.ServeHTTP(httptest.NewRecorder(), newPostWalletRequest("Habib", 340))

	// Calls Get method two times
	server.ServeHTTP(httptest.NewRecorder(), newGetWalletsRequest(user1))
	server.ServeHTTP(httptest.NewRecorder(), newGetWalletsRequest("Unknown User")) // Calls GET with unknown user

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetWalletsRequest("")) // Calls Get method
	assertStatus(t, response.Code, http.StatusOK) // asserts status code

	// expected response
	want := fmt.Sprintf("{\"%s\":%v,\"%s\":%d,\"%s\":%d}",
		user3, store.store[user3], user2, store.store[user2], user1, store.store[user1])

	assertResponseBody(t, response.Body.String(), want) // asserts response body
}
