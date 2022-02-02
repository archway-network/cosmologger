package block

import "time"

type BlockSignersRecord struct {
	BlockHeight uint64
	ValConsAddr string
	Time        time.Time
	Signature   string
}

type BlockRecord struct {
	BlockHash string
	Height    uint64
	NumOfTxs  uint64
	Time      time.Time
	Signers   []BlockSignersRecord
}
