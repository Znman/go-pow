package core

type Block struct {
	Index        int           // 区块编号
	TimeStamp    int64         // 区块时间戳
	Transactions []Transaction // 交易信息
	Proof        int64         // 工作量证明
	PreviousHash string        // 上一个区块的哈希
}