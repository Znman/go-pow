package core

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

const MINING_DIFFICULTY = 4

type Blockchain struct {
	Chain               []Block       `json:"chain"`
	CurrentTransactions []Transaction `json:"currentTransactions"`
}

func NewBlockChain() *Blockchain {
	bc := &Blockchain{
		Chain:               make([]Block, 0),
		CurrentTransactions: make([]Transaction, 0),
	}
	bc.CreateGenesisBlock()
	return bc
}

func (bc *Blockchain) GetLastBlock() Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) CreateGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: make([]Transaction, 0),
		Proof:        100,
		PreviousHash: "0",
	}
	genesisBlock.Hash = genesisBlock.CalculateHash()
	bc.Chain = append(bc.Chain, genesisBlock)
}

func (bc *Blockchain) AddTransaction(sender, recipient string, amount float64) {
	transaction := Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	bc.CurrentTransactions = append(bc.CurrentTransactions, transaction)
}

func (bc *Blockchain) MineBlock() Block {
	lastBlock := bc.GetLastBlock()
	
	newBlock := Block{
		Index:        len(bc.Chain),
		Timestamp:    time.Now().Unix(),
		Transactions: bc.CurrentTransactions,
		PreviousHash: lastBlock.Hash,
	}

	proof := bc.ProofOfWork(lastBlock.Proof)
	newBlock.Proof = proof
	newBlock.Hash = newBlock.CalculateHash()

	// Add the block to the chain
	bc.Chain = append(bc.Chain, newBlock)
	
	// Clear the current transactions
	bc.CurrentTransactions = make([]Transaction, 0)

	return newBlock
}

func (bc *Blockchain) ProofOfWork(lastProof int64) int64 {
	proof := int64(0)
	for !bc.ValidProof(lastProof, proof) {
		proof++
	}
	return proof
}

func (bc *Blockchain) ValidProof(lastProof, proof int64) bool {
	guess := []byte(string(lastProof) + string(proof))
	hash := sha256.Sum256(guess)
	hexHash := hex.EncodeToString(hash[:])
	return strings.HasPrefix(hexHash, strings.Repeat("0", MINING_DIFFICULTY))
}