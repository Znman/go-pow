package core

type Transaction struct {
    Sender    string  `json:"sender"`    // 发送者
    Recipient string  `json:"recipient"` // 接收者
    Amount    float64 `json:"amount"`    // 交易金额
}