package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"week-6-assignment-habibullahturkmen/config"
)

var (
	minimumBalanceAmount = config.Configuration.MinimumBalanceAmount
	initialBalanceAmount = config.Configuration.InitialBalanceAmount
)

type WalletStore interface {
	GetUsers(username string) map[string]int
	CreateUser(username string, initialBalance int) string
	UpdateUsers(username string, balance int, minimumValue int) (string, error)
}

type WalletServer struct {
	store WalletStore
}

func (p *WalletServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/")

	if r.Method == http.MethodPut {
		p.createUserWallet(w, username, initialBalanceAmount)
	} else if r.Method == http.MethodGet {
		p.showAllWallets(w, username)
	} else if r.Method == http.MethodPost {
		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		c := make(map[string]int)
		err = json.Unmarshal(b, &c)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		balance := c["balance"]
		p.updateUserWallet(w, username, balance)
	}


}

func (p *WalletServer) showAllWallets(w http.ResponseWriter, username string) {
	balance := p.store.GetUsers(username)

	if len(balance) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	m := map[string]int{}
	for key, value := range balance {
		m[key] = value
	}

	Json, _ := json.Marshal(m)
	w.Header().Add("content-type", "application/json")
	w.Write(Json)
}

func (p *WalletServer) createUserWallet(w http.ResponseWriter, username string, initialBalance int) {
	response := p.store.CreateUser(username, initialBalance)
	fmt.Println(response)
	w.WriteHeader(http.StatusAccepted)
}

func (p *WalletServer) updateUserWallet(w http.ResponseWriter, username string, balance int) {
	response, err := p.store.UpdateUsers(username, balance, minimumBalanceAmount)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	w.WriteHeader(http.StatusAccepted)
}