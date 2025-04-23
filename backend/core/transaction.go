package core

type Transaction struct {
	Sender    string  // 发送者
	Recipient string  // 接收者
	Amount    float64 // 交易金额
}