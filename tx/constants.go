package tx

const (
	MODULE_BANK         = "bank"
	MODULE_STAKING      = "staking"
	MODULE_DISTRIBUTION = "distribution"
	MODULE_GOVERNANCE   = "governance"
	MODULE_SLASHING     = "slashing"
)

const (
	ACTION_CREATE_VALIDATOR          = "create_validator"
	ACTION_SEND                      = "send"
	ACTION_DELEGATE                  = "delegate"
	ACTION_BEGIN_REDELEGATE          = "begin_redelegate"
	ACTION_BEGIN_UNBONDING           = "begin_unbonding"
	ACTION_WITHDRAW_DELEGATOR_REWARD = "withdraw_delegator_reward"
	ACTION_SUBMIT_PROPOSAL           = "submit_proposal"
	ACTION_VOTE                      = "vote"
	ACTION_UNJAIL                    = "unjail"
)
