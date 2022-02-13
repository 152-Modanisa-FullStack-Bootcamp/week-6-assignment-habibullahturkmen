//package data
//
//type IInMemoryWallet interface {
//	Create(string, int)
//	Get(string, int) int
//}
//
//type InMemoryWallet struct {
//	Wallet map[string]int
//}
//
//func (i *InMemoryWallet) Creat(username string, initialBalance int) {
//	i.Wallet[username] = initialBalance
//}
//
//func (i *InMemoryWallet) Get(username string) int {
//
//	return len(i.Wallet)
//	//return i.Wallet[username]
//}
//
//func NewInMemoryWallet() *InMemoryWallet {
//	return &InMemoryWallet{map[string]int{}}
//}

package main

import (
	"errors"
	"strings"
)

func NewInMemoryWallet() *InMemoryWalletStore {
	return &InMemoryWalletStore{map[string]int{}}
}

type InMemoryWalletStore struct {
	store map[string]int
}

func (i *InMemoryWalletStore) CreateUser(username string, initialBalance int) string {
	for key, _ := range i.store {
		if strings.ToLower(key) == strings.ToLower(username) {
			return "User already have a wallet!"
		}
	}
	i.store[username] = initialBalance
	return "Wallet created for " + username
}

func (i *InMemoryWalletStore) GetUsers(username string) map[string]int {
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

func (i *InMemoryWalletStore) UpdateUsers(username string, balance int, minimumValue int) (string, error) {
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