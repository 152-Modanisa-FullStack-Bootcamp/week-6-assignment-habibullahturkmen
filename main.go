package main

import (
	"log"
	"net/http"
)

func main() {
	server := &WalletServer{NewInMemoryWallet()}
	log.Fatal(http.ListenAndServe(":5000", server))
}