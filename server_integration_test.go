package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryWallet()
	server := WalletServer{store}

	user1 := "Recep"
	user2 := "Mehmet"
	user3 := "Davut"

	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user1))
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user2))
	server.ServeHTTP(httptest.NewRecorder(), newPutWalletRequest(user3))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetWalletsRequest(""))
	assertStatus(t, response.Code, http.StatusOK)

	want := fmt.Sprintf("{\"%s\":%d,\"%s\":%d,\"%s\":%d}", user3, 0, user2, 0, user1, 0)

	assertResponseBody(t, response.Body.String(), want)
}
