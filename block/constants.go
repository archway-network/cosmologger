package block

const (
	EVENT_ContractRewardCalculationEvent = "archway.gastracker.v1.ContractRewardCalculationEvent"
	EVENT_RewardDistributionEvent        = "archway.gastracker.v1.RewardDistributionEvent"

	EVENT_ContractRewardCalculationEvent_CONTRACT_ADDRESS  = EVENT_ContractRewardCalculationEvent + ".contract_address"
	EVENT_ContractRewardCalculationEvent_CONTRACT_REWARDS  = EVENT_ContractRewardCalculationEvent + ".contract_rewards"
	EVENT_ContractRewardCalculationEvent_GAS_CONSUMED      = EVENT_ContractRewardCalculationEvent + ".gas_consumed"
	EVENT_ContractRewardCalculationEvent_INFLATION_REWARDS = EVENT_ContractRewardCalculationEvent + ".inflation_rewards"
	EVENT_ContractRewardCalculationEvent_METADATA          = EVENT_ContractRewardCalculationEvent + ".metadata"
	EVENT_RewardDistributionEvent_CONTRACT_REWARDS         = EVENT_RewardDistributionEvent + ".contract_rewards"
	EVENT_RewardDistributionEvent_LEFTOVER_REWARDS         = EVENT_RewardDistributionEvent + ".leftover_rewards"
	EVENT_RewardDistributionEvent_REWARD_ADDRESS           = EVENT_RewardDistributionEvent + ".reward_address"

	EVENT_FIELD_DENOM                      = "denom"
	EVENT_FIELD_AMOUNT                     = "amount"
	EVENT_FIELD_DEVELOPER_ADDRESS          = "developer_address"
	EVENT_FIELD_REWARD_ADDRESS             = "reward_address"
	EVENT_FIELD_GAS_REBATE_TO_USER         = "gas_rebate_to_user"
	EVENT_FIELD_COLLECT_PREMIUM            = "collect_premium"
	EVENT_FIELD_PREMIUM_PERCENTAGE_CHARGED = "premium_percentage_charged"
)
