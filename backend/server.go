package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Znman/go-pow/backend/core"
)

var blockchain = core.NewBlockChain()

type TransactionRequest struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

type SearchResponse struct {
	Results []core.Block `json:"results"`
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain)
}

func searchBlocks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")

	var results []core.Block

	if searchType == "block" {
		if blockNum, err := strconv.Atoi(query); err == nil {
			if block := blockchain.GetBlockByIndex(blockNum); block != nil {
				results = append(results, *block)
			}
		}
	} else if searchType == "transaction" {
		results = blockchain.SearchTransactions(query)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SearchResponse{Results: results})
}

func addTransaction(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	index := blockchain.AddTransaction(req.Sender, req.Recipient, req.Amount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Transaction will be added to Block " + strconv.Itoa(index),
	})
}

func mineBlock(w http.ResponseWriter, r *http.Request) {
	block := blockchain.MineBlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/chain", getChain)
	mux.HandleFunc("/search", searchBlocks)
	mux.HandleFunc("/transactions/new", addTransaction)
	mux.HandleFunc("/mine", mineBlock)

	handler := enableCORS(mux)

	log.Println("Server running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}