package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/Znman/go-pow/backend/core"
)

var blockchain = core.NewBlockChain()

type TransactionRequest struct {
    Sender    string  `json:"sender"`
    Recipient string  `json:"recipient"`
    Amount    float64 `json:"amount"`
}

type TransactionResponse struct {
    Message     string           `json:"message"`
    Transaction core.Transaction `json:"transaction"`
}

type MineResponse struct {
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

    transaction := blockchain.AddTransaction(req.Sender, req.Recipient, req.Amount)

    response := TransactionResponse{
        Message:     "Transaction added successfully",
        Transaction: transaction,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func mineBlock(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    newBlock := blockchain.MineBlock()
    
    response := MineResponse{
        Message: "New Block Mined Successfully!",
        Block:   newBlock,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func getChain(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(blockchain)
}

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/chain", getChain)
    mux.HandleFunc("/transactions/new", addTransaction)
    mux.HandleFunc("/mine", mineBlock)

    handler := enableCORS(mux)

    log.Println("Server running on http://localhost:8000")
    log.Fatal(http.ListenAndServe(":8000", handler))
}