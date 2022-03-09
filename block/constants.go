package block

const (
	EVENTS_ContractRewardCalculationEvent = "archway.gastracker.v1.ContractRewardCalculationEvent"
	EVENTS_RewardDistributionEvent        = "archway.gastracker.v1.RewardDistributionEvent"
)

const (
	ACTION_CREATE_VALIDATOR          = "/cosmos.staking.v1beta1.MsgCreateValidator"              // "create_validator"
	ACTION_SEND                      = "/cosmos.bank.v1beta1.MsgSend"                            // "send"
	ACTION_DELEGATE                  = "/cosmos.staking.v1beta1.MsgDelegate"                     // "delegate"
	ACTION_BEGIN_REDELEGATE          = "/cosmos.staking.v1beta1.MsgBeginRedelegate"              // "begin_redelegate"
	ACTION_BEGIN_UNBONDING           = "/cosmos.staking.v1beta1.MsgUndelegate"                   // "begin_unbonding"
	ACTION_WITHDRAW_DELEGATOR_REWARD = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward" // "withdraw_delegator_reward"
	ACTION_SUBMIT_PROPOSAL           = "/cosmos.gov.v1beta1.MsgSubmitProposal"                   // "submit_proposal"
	ACTION_VOTE                      = "/cosmos.gov.v1beta1.MsgVote"                             // "vote"
	ACTION_UNJAIL                    = "/cosmos.slashing.v1beta1.MsgUnjail"                      // "unjail"
)
