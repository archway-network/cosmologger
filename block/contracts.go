package block

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/archway-network/cosmologger/database"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"google.golang.org/grpc"
)

func ProcessContractEvents(grpcCnn *grpc.ClientConn, evr *coretypes.ResultEvent, db *database.Database, insertQueue *database.InsertQueue) error {

	rec, err := getContractRecordFromEvent(evr)
	if err != nil || rec == nil {
		// rec&err==nil: Nothing to process
		return err
	}

	dbRow := rec.getDBRow()
	insertQueue.AddToInsertQueue(database.TABLE_CONTRACTS, dbRow)
	return nil
}

func getContractRecordFromEvent(evr *coretypes.ResultEvent) (*ContractRecord, error) {
	var cr ContractRecord

	if _, ok := evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS]; !ok {
		// Nothing to process
		return nil, nil
	}

	b := evr.Data.(tmTypes.EventDataNewBlock)
	cr.BlockHeight = uint64(b.Block.Height) - 1 // The gastracking is processed in the next beginBlock

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS]) > 0 {
		cr.ContractAddress = strings.Trim(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS][0], "\"")
	}

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_METADATA]) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(evr.Events[EVENT_ContractRewardCalculationEvent_METADATA][0]), &metadata); err != nil {
			return nil, err
		}

		cr.RewardAddress = metadata[EVENT_FIELD_REWARD_ADDRESS].(string)
		cr.DeveloperAddress = metadata[EVENT_FIELD_DEVELOPER_ADDRESS].(string)
		cr.GasRebateToUser = metadata[EVENT_FIELD_GAS_REBATE_TO_USER].(bool)
		cr.CollectPremium = metadata[EVENT_FIELD_COLLECT_PREMIUM].(bool)
		cr.MetadataJson = evr.Events[EVENT_ContractRewardCalculationEvent_METADATA][0]

		intValue, err := strconv.ParseUint(metadata[EVENT_FIELD_PREMIUM_PERCENTAGE_CHARGED].(string), 10, 64)
		if err != nil {
			return nil, err
		}
		cr.PremiumPercentageCharged = intValue
	}

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_GAS_CONSUMED]) > 0 {

		intValue, err := strconv.ParseUint(strings.Trim(evr.Events[EVENT_ContractRewardCalculationEvent_GAS_CONSUMED][0], "\""), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error in Unmarshaling '%s': %v", EVENT_ContractRewardCalculationEvent_GAS_CONSUMED, err)
		}
		cr.GasConsumed = intValue
	}

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_REWARDS]) > 0 {
		var err error
		cr.ContractRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_ContractRewardCalculationEvent_CONTRACT_REWARDS][0])
		if err != nil {
			return nil, err
		}
	}

	if len(evr.Events[EVENT_ContractRewardCalculationEvent_INFLATION_REWARDS]) > 0 {
		var err error
		cr.InflationRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_ContractRewardCalculationEvent_INFLATION_REWARDS][0])
		if err != nil {
			return nil, err
		}
	}

	if len(evr.Events[EVENT_RewardDistributionEvent_LEFTOVER_REWARDS]) > 0 {
		var err error
		cr.LeftoverRewards, err = getGasTrackerRewardFromString(evr.Events[EVENT_RewardDistributionEvent_LEFTOVER_REWARDS][0])
		if err != nil {
			return nil, err
		}
	}

	return &cr, nil
}

func (c *ContractRecord) getDBRow() database.RowType {
	return database.RowType{
		database.FIELD_CONTRACTS_CONTRACT_ADDRESS:           c.ContractAddress,
		database.FIELD_CONTRACTS_REWARD_ADDRESS:             c.RewardAddress,
		database.FIELD_CONTRACTS_DEVELOPER_ADDRESS:          c.DeveloperAddress,
		database.FIELD_CONTRACTS_BLOCK_HEIGHT:               c.BlockHeight,
		database.FIELD_CONTRACTS_GAS_CONSUMED:               c.GasConsumed,
		database.FIELD_CONTRACTS_REWARDS_DENOM:              c.ContractRewards.Denom,
		database.FIELD_CONTRACTS_CONTRACT_REWARDS_AMOUNT:    c.ContractRewards.Amount,
		database.FIELD_CONTRACTS_INFLATION_REWARDS_AMOUNT:   c.InflationRewards.Amount,
		database.FIELD_CONTRACTS_LEFTOVER_REWARDS_AMOUNT:    c.LeftoverRewards.Amount,
		database.FIELD_CONTRACTS_COLLECT_PREMIUM:            c.CollectPremium,
		database.FIELD_CONTRACTS_GAS_REBATE_TO_USER:         c.GasRebateToUser,
		database.FIELD_CONTRACTS_PREMIUM_PERCENTAGE_CHARGED: c.PremiumPercentageCharged,
		database.FIELD_CONTRACTS_METADATA_JSON:              c.MetadataJson,
	}
}

func getGasTrackerRewardFromString(str string) (GasTrackerReward, error) {

	// Let's make it an array if not, to keep compatibility
	if !strings.HasPrefix(str, "[") {
		str = "[" + str + "]"
	}

	var tmpMapArr []map[string]interface{}
	if err := json.Unmarshal([]byte(str), &tmpMapArr); err != nil {
		return GasTrackerReward{}, err
	}

	if len(tmpMapArr) == 0 {
		return GasTrackerReward{}, fmt.Errorf("no GasTrackerReward found")
	}
	tmpMap := tmpMapArr[0]

	numValue, err := strconv.ParseFloat(tmpMap[EVENT_FIELD_AMOUNT].(string), 64)
	if err != nil {
		return GasTrackerReward{}, err
	}

	return GasTrackerReward{
		Denom:  tmpMap[EVENT_FIELD_DENOM].(string),
		Amount: numValue,
	}, nil
}
