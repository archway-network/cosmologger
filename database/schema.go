package database

const TABLE_TX_EVENTS = "tx_events"

const (
	FIELD_TX_EVENTS_TX_HASH      = "txHash"
	FIELD_TX_EVENTS_HEIGHT       = "height"
	FIELD_TX_EVENTS_MODULE       = "module"
	FIELD_TX_EVENTS_SENDER       = "sender"
	FIELD_TX_EVENTS_RECEIVER     = "receiver"
	FIELD_TX_EVENTS_VALIDATOR    = "validator"
	FIELD_TX_EVENTS_ACTION       = "action"
	FIELD_TX_EVENTS_AMOUNT       = "amount"
	FIELD_TX_EVENTS_TX_ACCSEQ    = "txAccSeq"
	FIELD_TX_EVENTS_TX_SIGNATURE = "txSignature"
	FIELD_TX_EVENTS_PROPOSAL_ID  = "proposalId"
	FIELD_TX_EVENTS_TX_MEMO      = "txMemo"
	FIELD_TX_EVENTS_JSON         = "json"
	FIELD_TX_EVENTS_LOG_TIME     = "logTime"
)

/*-------------------*/

const TABLE_BLOCKS = "blocks"

const (
	FIELD_BLOCKS_BLOCK_HASH = "blockHash"
	FIELD_BLOCKS_HEIGHT     = "height"
	FIELD_BLOCKS_NUM_OF_TXS = "numOfTxs"
	FIELD_BLOCKS_TIME       = "time"
)

/*-------------------*/

const TABLE_BLOCK_SIGNERS = "block_signers"

const (
	FIELD_BLOCK_SIGNERS_BLOCK_HEIGHT  = "blockHeight"
	FIELD_BLOCK_SIGNERS_VAL_CONS_ADDR = "valConsAddr"
	FIELD_BLOCK_SIGNERS_TIME          = "time"
	FIELD_BLOCK_SIGNERS_SIGNATURE     = "signature"
)

/*-------------------*/

const TABLE_VALIDATORS = "validators"

const (
	FIELD_VALIDATORS_CONS_ADDR    = "consAddr"
	FIELD_VALIDATORS_OPR_ADDR     = "oprAddr"
	FIELD_VALIDATORS_ACCOUNT_ADDR = "accountAddr"
	FIELD_VALIDATORS_MONIKER      = "moniker"
)

/*-------------------*/

const TABLE_PARTICIPANTS = "participants"

const (
	FIELD_PARTICIPANTS_EMAIL_ADDRESS   = "emailAddress"
	FIELD_PARTICIPANTS_FULL_LEGAL_NAME = "fullLegalName"
	FIELD_PARTICIPANTS_COUNTRY         = "country"
	FIELD_PARTICIPANTS_ACCOUNT_ADDRESS = "accountAddress"
	FIELD_PARTICIPANTS_GITHUB_HANDLE   = "githubHandle"
	FIELD_PARTICIPANTS_PUBKEY          = "pubkey"
	FIELD_PARTICIPANTS_KYC_SESSION_ID  = "kycSessionId"
	FIELD_PARTICIPANTS_KYC_VERIFIED    = "kycVerified"
)

/*-------------------*/

const TABLE_CONTRACTS = "contracts"

const (
	FIELD_CONTRACTS_CONTRACT_ADDRESS           = "contractAddress"
	FIELD_CONTRACTS_REWARD_ADDRESS             = "rewardAddress"
	FIELD_CONTRACTS_DEVELOPER_ADDRESS          = "developerAddress"
	FIELD_CONTRACTS_BLOCK_HEIGHT               = "blockHeight"
	FIELD_CONTRACTS_GAS_CONSUMED_STR           = "gasConsumedStr" //Let's keep it for the moment, just to verify if the number is big enough
	FIELD_CONTRACTS_GAS_CONSUMED               = "gasConsumed"
	FIELD_CONTRACTS_REWARDS_DENOM              = "rewardsDenom"
	FIELD_CONTRACTS_CONTRACT_REWARDS_AMOUNT    = "contractRewardsAmount"
	FIELD_CONTRACTS_INFLATION_REWARDS_AMOUNT   = "inflationRewardsAmount"
	FIELD_CONTRACTS_LEFTOVER_REWARDS_AMOUNT    = "leftoverRewardsAmount"
	FIELD_CONTRACTS_COLLECT_PREMIUM            = "collectPremium"
	FIELD_CONTRACTS_GAS_REBATE_TO_USER         = "gasRebateToUser"
	FIELD_CONTRACTS_PREMIUM_PERCENTAGE_CHARGED = "premiumPercentageCharged"
	FIELD_CONTRACTS_METADATA_JSON              = "metadataJson"
	FIELD_CONTRACTS_INC_ID                     = "incId"
)

/*-------------------*/
