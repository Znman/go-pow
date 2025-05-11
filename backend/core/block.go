package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	Index        int          `json:"index"`
	Timestamp    int64        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int64        `json:"proof"`
	PreviousHash string       `json:"previousHash"`
	Hash         string       `json:"hash"`
}

func (b *Block) CalculateHash() string {
	data := struct {
		Index        int
		Timestamp    int64
		Transactions []Transaction
		Proof        int64
		PreviousHash string
	}{
		Index:        b.Index,
		Timestamp:    b.Timestamp,
		Transactions: b.Transactions,
		Proof:        b.Proof,
		PreviousHash: b.PreviousHash,
	}

	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}