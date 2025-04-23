package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"blockchain-backend/core"
)

var blockchain = core.NewBlockChain()

func main() {
	http.HandleFunc("/chain", getChain)
	http.HandleFunc("/transactions/new", addTransaction)
	http.HandleFunc("/mine", mineBlock)

	log.Println("Running server at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getChain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain)
}

func addTransaction(w http.ResponseWriter, r *http.Request) {
	var t core.Transaction
	_ = json.NewDecoder(r.Body).Decode(&t)
	blockchain.AddTransaction(t.Sender, t.Recipient, t.Amount)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Transaction added!")
}

func mineBlock(w http.ResponseWriter, r *http.Request) {
	block := blockchain.MineBlock(100)
	json.NewEncoder(w).Encode(block)
}