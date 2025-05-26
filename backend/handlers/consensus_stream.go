package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Znman/go-pow/backend/core"
)

type ConsensusProgress struct {
	Attempt    int    `json:"attempt"`
	Proof      int64  `json:"proof"`
	Hash       string `json:"hash"`
	Found      bool   `json:"found"`
	BlockIndex int    `json:"blockIndex"`
	Message    string `json:"message"`
}

// Factory function that returns an http.HandlerFunc bound to the given blockchain instance
func MineConsensusStreamHandler(bc *core.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		lastBlock := bc.GetLastBlock()
		blockIndex := lastBlock.Index + 1

		var proof int64 = 0
		var attempt int = 0
		var hash string

		for {
			hash = bc.HashProof(lastBlock.Proof, proof)
			found := bc.ValidProof(lastBlock.Proof, proof)
			progress := ConsensusProgress{
				Attempt:    attempt,
				Proof:      proof,
				Hash:       hash,
				Found:      found,
				BlockIndex: blockIndex,
				Message:    "",
			}
			if found {
				progress.Message = "Valid proof found!"
			} else {
				progress.Message = "Trying next proof..."
			}
			jsonData, _ := json.Marshal(progress)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()

			if found {
				// Actually mine the block (add to chain)
				bc.MineBlock()
				break
			}
			proof++
			attempt++
			time.Sleep(15 * time.Millisecond)
		}

		finalMsg := ConsensusProgress{
			Attempt:    attempt,
			Proof:      proof,
			Hash:       hash,
			Found:      true,
			BlockIndex: blockIndex,
			Message:    "Block mined and added to chain.",
		}
		jsonData, _ := json.Marshal(finalMsg)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	}
}
