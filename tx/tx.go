package tx

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/archway-network/cosmologger/configs"
	"github.com/archway-network/cosmologger/database"
	"github.com/archway-network/cosmologger/validators"

	// sdkClient "github.com/cosmos/cosmos-sdk/client"
	// authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
)

func ProcessEvents(grpcCnn *grpc.ClientConn, evr coretypes.ResultEvent, db *database.Database, insertQueue *database.InsertQueue) error {

	rec := getTxRecordFromEvent(evr)
	rec.LogTime = time.Now()

	dbRow := rec.getDBRow()

	qRes, _ := db.Load(database.TABLE_TX_EVENTS, database.RowType{database.FIELD_TX_EVENTS_TX_HASH: rec.TxHash})
	if len(qRes) > 0 && rec.Module != "" {
		// This tx is already in the DB, let's update it
		go func() {
			_, err := db.Update(database.TABLE_TX_EVENTS, dbRow, database.RowType{database.FIELD_TX_EVENTS_TX_HASH: rec.TxHash})
			if err != nil {
				log.Printf("Err in `Update TX`: %v", err)
			}
		}()

	} else {

		insertQueue.AddToInsertQueue(database.TABLE_TX_EVENTS, dbRow)
	}

	// Let's add validator's info
	if rec.Validator != "" ||
		rec.Action == ACTION_UNJAIL {
		// Just to make things non-blocking
		go func() {

			// When `unjail` actions is invoked, the validator address is in the `sender` filed (well mostly :D)
			if rec.Action == ACTION_UNJAIL &&
				strings.HasPrefix(rec.Sender, configs.Configs.Bech32Prefix.Validator.Address) {

				rec.Validator = rec.Sender
			}

			if rec.Validator != "" {

				err := validators.AddNewValidator(db, grpcCnn, rec.Validator)
				if err != nil {
					log.Printf("Err in `AddNewValidator`: %v", err)
					// return err
				}
			}
		}()
	}

	return nil
}

func getTxRecordFromEvent(evr coretypes.ResultEvent) TxRecord {
	var txRecord TxRecord

	if evr.Events["tx.height"] != nil && len(evr.Events["tx.height"]) > 0 {
		txRecord.Height, _ = strconv.ParseUint(evr.Events["tx.height"][0], 10, 64)
	}

	if evr.Events["tx.hash"] != nil && len(evr.Events["tx.hash"]) > 0 {
		txRecord.TxHash = evr.Events["tx.hash"][0]
	}

	if evr.Events["message.module"] != nil && len(evr.Events["message.module"]) > 0 {
		txRecord.Module = evr.Events["message.module"][0]
	}

	if evr.Events["message.sender"] != nil && len(evr.Events["message.sender"]) > 0 {
		txRecord.Sender = evr.Events["message.sender"][0]

	} else if evr.Events["transfer.sender"] != nil && len(evr.Events["transfer.sender"]) > 0 {

		txRecord.Sender = evr.Events["transfer.sender"][0]
	}

	if evr.Events["transfer.recipient"] != nil && len(evr.Events["transfer.recipient"]) > 0 {
		txRecord.Receiver = evr.Events["transfer.recipient"][0]
	}

	if evr.Events["delegate.validator"] != nil && len(evr.Events["delegate.validator"]) > 0 {
		txRecord.Validator = evr.Events["delegate.validator"][0]

	} else if evr.Events["create_validator.validator"] != nil && len(evr.Events["create_validator.validator"]) > 0 {

		txRecord.Validator = evr.Events["create_validator.validator"][0]
	}

	if evr.Events["message.action"] != nil && len(evr.Events["message.action"]) > 0 {
		txRecord.Action = evr.Events["message.action"][0]
	}

	if evr.Events["delegate.amount"] != nil && len(evr.Events["delegate.amount"]) > 0 {
		txRecord.Amount = evr.Events["delegate.amount"][0]

	} else if evr.Events["transfer.amount"] != nil && len(evr.Events["transfer.amount"]) > 0 {

		txRecord.Amount = evr.Events["transfer.amount"][0]
	}

	if evr.Events["tx.acc_seq"] != nil && len(evr.Events["tx.acc_seq"]) > 0 {
		txRecord.TxAccSeq = evr.Events["tx.acc_seq"][0]
	}

	if evr.Events["tx.signature"] != nil && len(evr.Events["tx.signature"]) > 0 {
		txRecord.TxSignature = evr.Events["tx.signature"][0]
	}

	if evr.Events["proposal_vote.proposal_id"] != nil && len(evr.Events["proposal_vote.proposal_id"]) > 0 {
		txRecord.ProposalId, _ = strconv.ParseUint(evr.Events["proposal_vote.proposal_id"][0], 10, 64)

	} else if evr.Events["proposal_deposit.proposal_id"] != nil && len(evr.Events["proposal_deposit.proposal_id"]) > 0 {

		txRecord.ProposalId, _ = strconv.ParseUint(evr.Events["proposal_deposit.proposal_id"][0], 10, 64)
	}

	// Memo cannot be retrieved through tx events, we may fill it up with another way later
	// txRecord.TxMemo =

	jsonBytes, err := json.Marshal(evr.Events)
	if err == nil {
		txRecord.Json = string(jsonBytes)
	}

	// LogTime: is recorded by the DBMS itself

	return txRecord
}

func (t TxRecord) getDBRow() database.RowType {
	return database.RowType{

		database.FIELD_TX_EVENTS_TX_HASH:      t.TxHash,
		database.FIELD_TX_EVENTS_HEIGHT:       t.Height,
		database.FIELD_TX_EVENTS_MODULE:       t.Module,
		database.FIELD_TX_EVENTS_SENDER:       t.Sender,
		database.FIELD_TX_EVENTS_RECEIVER:     t.Receiver,
		database.FIELD_TX_EVENTS_VALIDATOR:    t.Validator,
		database.FIELD_TX_EVENTS_ACTION:       t.Action,
		database.FIELD_TX_EVENTS_AMOUNT:       t.Amount,
		database.FIELD_TX_EVENTS_TX_ACCSEQ:    t.TxAccSeq,
		database.FIELD_TX_EVENTS_TX_SIGNATURE: t.TxSignature,
		database.FIELD_TX_EVENTS_PROPOSAL_ID:  t.ProposalId,
		database.FIELD_TX_EVENTS_TX_MEMO:      t.TxMemo,
		database.FIELD_TX_EVENTS_JSON:         t.Json,
		database.FIELD_TX_EVENTS_LOG_TIME:     t.LogTime,
	}
}

func Start(cli *tmClient.HTTP, grpcCnn *grpc.ClientConn, db *database.Database, insertQueue *database.InsertQueue) {

	go func() {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(configs.Configs.GRPC.CallTimeout))
		defer cancel()

		eventChan, err := cli.Subscribe(ctx, configs.Configs.TendermintClient.SubscriberName, tmTypes.QueryForEvent(tmTypes.EventTx).String())
		if err != nil {
			panic(err)
		}

		for {
			evRes := <-eventChan
			err := ProcessEvents(grpcCnn, evRes, db, insertQueue)
			if err != nil {
				log.Printf("Error in processing TX event: %v", err)
			}
		}
	}()

	fixEmptyEvents(cli, db)
}

// Since some TX events are delayed and we catch them empty, we need to query them later to get them fixed
func fixEmptyEvents(cli *tmClient.HTTP, db *database.Database) {

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Kill, //nolint
		os.Interrupt)

	go func() {

		for {
			select {
			case <-quitChannel:
				return
			default:

				// Get all TX events that are empty
				rows, err := db.Load(database.TABLE_TX_EVENTS, database.RowType{database.FIELD_TX_EVENTS_MODULE: ""})
				if err != nil {
					log.Printf("Error in loading empty TX events: %v", err)
				}

				for _, row := range rows {
					txHash := string(row[database.FIELD_TX_EVENTS_TX_HASH].([]uint8))
					// Quering the TX from the Node...
					rec, err := queryTx(cli, txHash)
					if err != nil {
						log.Printf("Error in querying TX: %s\t %v", txHash, err)
						continue
					}
					rec.LogTime = time.Now()
					dbRow := rec.getDBRow()

					_, err = db.Update(database.TABLE_TX_EVENTS, dbRow, database.RowType{database.FIELD_TX_EVENTS_TX_HASH: rec.TxHash})
					if err != nil {
						log.Printf("[FixEmptyEvents] Err in `Update TX`: %s\t %v", txHash, err)
					}
				}

				time.Sleep(time.Second)
			}
		}
	}()

}

func queryTx(cli *tmClient.HTTP, txHash string) (TxRecord, error) {

	wsURI := os.Getenv("RPC_ADDRESS")

	// A dirty hack to get the things done
	cmd := exec.Command("archwayd", "query", "tx", txHash, "--node", wsURI, "--output", "json")
	stdout, err := cmd.Output()

	if err != nil {
		return TxRecord{}, err
	}

	rec := getTxRecordFromJson(stdout)

	return rec, nil

}

func getTxRecordFromJson(jsonByte []byte) TxRecord {
	var txRecord TxRecord
	// jsonStr = strings.Trim(jsonStr, " \r\n\t")

	jsonVar := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonByte), &jsonVar)
	if err != nil {
		fmt.Printf("Unmarshaling JSON Err: %v\n", err.Error())
		return txRecord
	}

	if jsonVar["height"] != nil && len(jsonVar["height"].(string)) > 0 {
		txRecord.Height, _ = strconv.ParseUint(jsonVar["height"].(string), 10, 64)
	}

	if jsonVar["txhash"] != nil && len(jsonVar["txhash"].(string)) > 0 {
		txRecord.TxHash = jsonVar["txhash"].(string)
	}

	if jsonVar["codespace"] != nil && len(jsonVar["codespace"].(string)) > 0 {
		txRecord.Module = jsonVar["codespace"].(string)
	}

	messages := []interface{}{}
	if txJson, ok := jsonVar["tx"].(map[string]interface{}); ok {
		if body, ok := txJson["body"].(map[string]interface{}); ok {
			if msgs, ok := body["messages"].([]interface{}); ok {
				messages = msgs
			}
		}

		if val, ok := txJson["signatures"].([]interface{}); ok {
			txRecord.TxSignature = val[0].(string)
		}
	}

	for i := range messages {
		msg := messages[i].(map[string]interface{})
		if val, ok := msg["@type"].(string); ok {
			if val == "" {
				val = "NA"
			}
			txRecord.Action = val
		}

		if val, ok := msg["sender"].(string); ok {
			txRecord.Sender = val
		} else if val, ok := msg["delegator_address"].(string); ok {
			txRecord.Sender = val
		} else if val, ok := msg["inputs"].([]interface{}); ok {

			if addr, ok := val[0].(map[string]interface{})["address"].(string); ok {
				txRecord.Sender = addr
			}
		}

		if val, ok := msg["validator_address"].(string); ok {
			txRecord.Receiver = val
			txRecord.Validator = val
		} else if val, ok := msg["recipient"].(string); ok {
			txRecord.Receiver = val
		} else if val, ok := msg["outputs"].([]interface{}); ok {
			if addr, ok := val[0].(map[string]interface{})["address"].(string); ok {
				txRecord.Receiver = addr
			}
		}

		if val, ok := msg["value"].(map[string]interface{}); ok {
			txRecord.Amount = val["amount"].(string) + val["denom"].(string)
		} else if val, ok := msg["amount"].(map[string]interface{}); ok {
			txRecord.Amount = val["amount"].(string) + val["denom"].(string)
		}

	}

	// if jsonVar["proposal_vote.proposal_id"] != nil && len(jsonVar["proposal_vote.proposal_id"]) > 0 {
	// 	txRecord.ProposalId, _ = strconv.ParseUint(jsonVar["proposal_vote.proposal_id"][0], 10, 64)

	// } else if jsonVar["proposal_deposit.proposal_id"] != nil && len(jsonVar["proposal_deposit.proposal_id"]) > 0 {

	// 	txRecord.ProposalId, _ = strconv.ParseUint(jsonVar["proposal_deposit.proposal_id"][0], 10, 64)
	// }

	if txRecord.Module == "" {
		txRecord.Module = "NA"
	}

	txRecord.Json = string(jsonByte)

	// LogTime: is recorded by the DBMS itself

	return txRecord
}
