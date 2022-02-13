package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Stub for wallet
type StubWalletStore struct {
	store map[string]int
}

// GetUsers Stub
func (i *StubWalletStore) GetUsers(username string) map[string]int {
	// returns user and balance
	m := map[string]int{}
	for key, value := range i.store {
		if strings.ToLower(key) == strings.ToLower(username) {
			m[key] = value
			return m
		}
	}

	// returns empty map
	if len(username) > 0 {
		return m
	}

	// returns all the wallets
	return i.store
}

// CreateUser Stub
func (i *StubWalletStore) CreateUser(username string, initialBalance int) string {
	for key, _ := range i.store {
		if strings.ToLower(key) == strings.ToLower(username) {
			return "User already have a wallet!"
		}
	}
	i.store[username] = initialBalance
	return "Wallet created for " + username
}

// UpdateUsers Stub
func (i *StubWalletStore) UpdateUsers(username string, balance int, minimumValue int) (string, error) {
	for key, value := range i.store {
		if strings.ToLower(key) == strings.ToLower(username) {
			if value + balance < minimumValue {
				return "", errors.New("should not be less than minimum balance amount")
			}
			i.store[key] = value + balance
			return "Wallet updated for " + username, nil
		}
	}
	return "User not found!", nil
}

// Test for GET request
func TestGetUsers(t *testing.T) {
	store := StubWalletStore{
		map[string]int{
			"Mehmet": 20,
			"Burcu":  10,
			"Davut":  130,
		},
	}
	server := &WalletServer{&store}

	// returns users wallets - test
	for username, balance := range store.store {
		t.Run(fmt.Sprintf("returns %s's wallet", username), func(t *testing.T) {
			request := newGetWalletsRequest(username)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			want := fmt.Sprintf("{\"%s\":%d}", username, balance)

			assertStatus(t, response.Code, http.StatusOK)
			assertResponseBody(t, response.Body.String(), want)
		})
	}

	// missing users - test
	t.Run("returns 404 on missing users", func(t *testing.T) {
		request := newGetWalletsRequest("Cansu")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

// Test for PUT request
func TestPutUsers(t *testing.T) {
	store := StubWalletStore{
		map[string]int{},
	}
	server := &WalletServer{&store}

	users := []string{"Habib", "Ahmet", "Gonca"}

	for _, username := range users {
		t.Run(fmt.Sprintf("creates wallet for %s", username), func(t *testing.T) {
			request := newPutWalletRequest(username)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertCreateWallet(t, store.store, username)
			assertStatus(t, response.Code, http.StatusAccepted)
			assertResponseBody(t, response.Body.String(), "")
		})
	}
}

// Test for POST request
func TestPostUser(t *testing.T) {
	store := StubWalletStore{
		map[string]int{
			"Mehmet": 20,
			"Burcu":  10,
			"Davut":  130,
		},
	}
	server := &WalletServer{&store}
	username := "Mehmet"
	balance := -500
	request := newPostWalletRequest(username, balance)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)
	fmt.Println(store)

	assertStatus(t, response.Code, http.StatusAccepted)
	assertResponseBody(t, response.Body.String(), "")

}

// GET all user wallets
func newGetWalletsRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", name), nil)
	return request
}

// Put/Create user wallet
func newPutWalletRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/%s", name), nil)
	return request
}

// POST/Update user wallet
func newPostWalletRequest(name string, balance int) *http.Request {
	value := map[string]int{"balance": balance}
	jsonValue, _ := json.Marshal(value)
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/%s", name), bytes.NewBuffer(jsonValue))
	return request
}

// assert method for response body
func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

// assert method for status request
func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

// assert method for creat wallet
func assertCreateWallet(t testing.TB, store map[string]int, username string) {
	t.Helper()

	if _, ok := store[username]; !ok {
		t.Errorf("store doesn't contain the user, got %v, want %v", ok, username)
	}
}