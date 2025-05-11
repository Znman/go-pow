package core

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
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

func (bc *Blockchain) GetLastBlock() Block {
    return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) AddTransaction(sender, recipient string, amount float64) int {
    transaction := Transaction{
        Sender:    sender,
        Recipient: recipient,
        Amount:    amount,
    }
    bc.CurrentTransactions = append(bc.CurrentTransactions, transaction)
    return bc.GetLastBlock().Index + 1
}

func (bc *Blockchain) GetBlockByIndex(index int) *Block {
    for i := range bc.Chain {
        if bc.Chain[i].Index == index {
            return &bc.Chain[i]
        }
    }
    return nil
}

func (bc *Blockchain) SearchTransactions(query string) []Block {
    query = strings.ToLower(query)
    var results []Block

    for _, block := range bc.Chain {
        for _, tx := range block.Transactions {
            if strings.Contains(strings.ToLower(tx.Sender), query) ||
                strings.Contains(strings.ToLower(tx.Recipient), query) ||
                strings.Contains(strings.ToLower(fmt.Sprintf("%.2f", tx.Amount)), query) {
                results = append(results, block)
                break
            }
        }
    }

    return results
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

    bc.Chain = append(bc.Chain, newBlock)
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
    guess := []byte(fmt.Sprintf("%d%d", lastProof, proof))
    hash := sha256.Sum256(guess)
    hexHash := hex.EncodeToString(hash[:])
    return strings.HasPrefix(hexHash, strings.Repeat("0", MINING_DIFFICULTY))
}