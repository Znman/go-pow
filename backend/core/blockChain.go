package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type BlockChain struct {
	Chain               []Block       // 区块链
	CurrentTransactions []Transaction // 当前交易池
}

func NewBlockChain() *BlockChain {
	bc := &BlockChain{}
	bc.CreateGenesisBlock()
	return bc
}

func (bc *BlockChain) CreateGenesisBlock() {
	genesisBlock := Block{
		Index:        0,
		TimeStamp:    time.Now().Unix(),
		Transactions: nil,
		Proof:        100,
		PreviousHash: "0",
	}
	bc.Chain = append(bc.Chain, genesisBlock)
}

func (bc *BlockChain) AddTransaction(sender, recipient string, amount float64) {
	transaction := Transaction{Sender: sender, Recipient: recipient, Amount: amount}
	bc.CurrentTransactions = append(bc.CurrentTransactions, transaction)
}

func (bc *BlockChain) MineBlock(proof int64) *Block {
	lastBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := Block{
		Index:        len(bc.Chain),
		TimeStamp:    time.Now().Unix(),
		Transactions: bc.CurrentTransactions,
		Proof:        proof,
		PreviousHash: Hash(lastBlock),
	}
	bc.Chain = append(bc.Chain, newBlock)
	bc.CurrentTransactions = nil
	return &newBlock
}

func Hash(block Block) string {
	data, _ := json.Marshal(block)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}