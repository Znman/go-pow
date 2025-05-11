package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/Znman/go-pow/backend/core"
)

var blockchain = core.NewBlockChain()

type MineResponse struct {
    Message string     `json:"message"`
    Block   core.Block `json:"block"`
}

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.Header().Set("Content-Type", "application/json")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
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
    
    // Register routes
    mux.HandleFunc("/chain", getChain)
    mux.HandleFunc("/mine", mineBlock)

    // Wrap with CORS middleware
    handler := enableCORS(mux)

    log.Println("Server running on http://localhost:8000")
    log.Fatal(http.ListenAndServe(":8000", handler))
}