package core

type Transaction struct {
    Sender    string  `json:"sender"`    
    Recipient string  `json:"recipient"` 
    Amount    float64 `json:"amount"`    
    Message   string  `json:"message"`   // Added message field
}