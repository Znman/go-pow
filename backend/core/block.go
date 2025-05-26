package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	Index        int           `json:"index"`
	TimeStamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int64         `json:"proof"`
	PreviousHash string        `json:"previousHash"`
	Hash         string        `json:"hash"`
	MiningStats  MiningStats   `json:"miningStats"`
}

// Add CalculateHash method to Block
func (b *Block) CalculateHash() string {
	// Create a data structure for hashing that includes all fields except Hash
	data := struct {
		Index        int           `json:"index"`
		TimeStamp    int64         `json:"timestamp"`
		Transactions []Transaction `json:"transactions"`
		Proof        int64         `json:"proof"`
		PreviousHash string        `json:"previousHash"`
	}{
		Index:        b.Index,
		TimeStamp:    b.TimeStamp,
		Transactions: b.Transactions,
		Proof:        b.Proof,
		PreviousHash: b.PreviousHash,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}
