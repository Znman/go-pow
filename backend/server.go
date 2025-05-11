package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Znman/go-pow/backend/core"
)

var blockchain = core.NewBlockChain()

// Define the response structures
type TransactionRequest struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

type MiningResponse struct {
	Message string     `json:"message"`
	Block   core.Block `json:"block"`
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

func addTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	blockchain.AddTransaction(req.Sender, req.Recipient, req.Amount)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction added successfully",
	})
}

func mineBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	block := blockchain.MineBlock()

	response := MiningResponse{
		Message: "New Block Mined!",
		Block:   block,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/chain", getChain)
	mux.HandleFunc("/transactions/new", addTransaction)
	mux.HandleFunc("/mine", mineBlock)

	handler := enableCORS(mux)

	log.Println("Running server at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}