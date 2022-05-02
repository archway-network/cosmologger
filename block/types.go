package block

import "time"

type BlockSignersRecord struct {
	BlockHeight uint64
	ValConsAddr string
	Time        time.Time
	Signature   string
}

type BlockRecord struct {
	BlockHash        string
	Height           uint64
	NumOfTxs         uint64
	Time             time.Time
	LastBlockSigners []BlockSignersRecord
}

type ContractRecord struct {
	ContractAddress  string
	RewardAddress    string
	DeveloperAddress string
	BlockHeight      uint64

	GasConsumed      uint64           // in DB we changed it to varchar(50) as postgresql does not support uint64
	ContractRewards  GasTrackerReward // For sake of simplicity, we consider only one denom per record
	InflationRewards GasTrackerReward
	LeftoverRewards  GasTrackerReward

	CollectPremium           bool
	GasRebateToUser          bool
	PremiumPercentageCharged uint64

	MetadataJson string
}

type GasTrackerReward struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}
