package block

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/archway-network/cosmologger/configs"
	"github.com/archway-network/cosmologger/database"
	"github.com/archway-network/cosmologger/validators"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
)

var genesisValidatorsDone bool

func ProcessEvents(grpcCnn *grpc.ClientConn, evr *coretypes.ResultEvent, db *database.Database, insertQueue *database.InsertQueue) error {
	rec := getBlockRecordFromEvent(evr)
	fmt.Printf("Block: %s\tH: %d\tTxs: %d\n", rec.BlockHash, rec.Height, rec.NumOfTxs)

	dbRow := rec.getBlockDBRow()
	insertQueue.AddToInsertQueue(database.TABLE_BLOCKS, dbRow)

	dbRows := make([]database.RowType, len(rec.LastBlockSigners))
	for i := range rec.LastBlockSigners {
		dbRows[i] = rec.LastBlockSigners[i].getBlockSignerDBRow()
	}
	insertQueue.AddToInsertQueue(database.TABLE_BLOCK_SIGNERS, dbRows...)

	// Let's add genesis validator's info
	if !genesisValidatorsDone && rec.Height > 20 {
		// Let's do it only once
		genesisValidatorsDone = true

		// Just to make things non-blocking
		go func() {

			valList, err := validators.QueryValidatorsList(grpcCnn)
			if err != nil {
				log.Printf("Err in `validators.QueryValidatorsList`: %v", err)
				// return err
			}

			for i := range valList {
				err := validators.AddNewValidator(db, grpcCnn, valList[i])
				if err != nil {
					log.Printf("Err in `AddNewValidator`: %v", err)
					// return err
				}
			}

		}()
	}

	return ProcessContractEvents(grpcCnn, evr, db, insertQueue)
}

func getBlockRecordFromEvent(evr *coretypes.ResultEvent) *BlockRecord {

	b := evr.Data.(tmTypes.EventDataNewBlock)
	return getBlockRecordFromTmBlock(b.Block)
}

func getBlockRecordFromTmBlock(b *tmTypes.Block) *BlockRecord {
	var br BlockRecord

	br.BlockHash = b.Hash().String()

	br.Height = uint64(b.Height)
	br.NumOfTxs = uint64(len(b.Txs))
	br.Time = b.Time

	for i := range b.LastCommit.Signatures {

		consAddr, err := sdk.ConsAddressFromHex(b.LastCommit.Signatures[i].ValidatorAddress.String())
		if err != nil {
			continue // just ignore this signer as it might not be running and we face some strange error
		}

		br.LastBlockSigners = append(br.LastBlockSigners, BlockSignersRecord{
			BlockHeight: br.Height - 1, // Because the signers are for the previous block
			ValConsAddr: consAddr.String(),
			Time:        b.LastCommit.Signatures[i].Timestamp,
			Signature:   base64.StdEncoding.EncodeToString(b.LastCommit.Signatures[i].Signature),
		})
	}

	return &br
}

func (b *BlockRecord) getBlockDBRow() database.RowType {
	return database.RowType{
		database.FIELD_BLOCKS_BLOCK_HASH: b.BlockHash,
		database.FIELD_BLOCKS_HEIGHT:     b.Height,
		database.FIELD_BLOCKS_NUM_OF_TXS: b.NumOfTxs,
		database.FIELD_BLOCKS_TIME:       b.Time,
	}
}

func (s *BlockSignersRecord) getBlockSignerDBRow() database.RowType {
	return database.RowType{
		database.FIELD_BLOCK_SIGNERS_BLOCK_HEIGHT:  s.BlockHeight,
		database.FIELD_BLOCK_SIGNERS_VAL_CONS_ADDR: s.ValConsAddr,
		database.FIELD_BLOCK_SIGNERS_TIME:          s.Time,
		database.FIELD_BLOCK_SIGNERS_SIGNATURE:     s.Signature,
	}
}

func Start(cli *tmClient.HTTP, grpcCnn *grpc.ClientConn, db *database.Database, insertQueue *database.InsertQueue) {

	go func() {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(configs.Configs.GRPC.CallTimeout))
		defer cancel()

		eventChan, err := cli.Subscribe(
			ctx,
			configs.Configs.TendermintClient.SubscriberName,
			tmTypes.QueryForEvent(tmTypes.EventNewBlock).String(),
		)
		if err != nil {
			panic(err)
		}

		for {
			evRes := <-eventChan
			err := ProcessEvents(grpcCnn, &evRes, db, insertQueue)
			if err != nil {
				//TODO: We need some customizable log level
				log.Printf("Error in processing block event: %v", err)
			}
		}
	}()

	fixMissingBlocks(cli, db, insertQueue)
}

// Sometimes some blocks get missed, so this function attempts to find them and fix them
func fixMissingBlocks(cli *tmClient.HTTP, db *database.Database, insertQueue *database.InsertQueue) {

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

				// Get all missing blocks
				latestBlockHeight, err := GetLatestBlockHeight(db)
				if err != nil {
					log.Printf("Error in finding LatestBlockHeight: %v", err)
				}
				missingBlocks, err := findMissingBlocks(1, latestBlockHeight, db)
				if err != nil {
					log.Printf("Error in finding missing blocks: %v", err)
				}

				for _, bh := range missingBlocks {

					// Querying the Block from the Node...
					rec, txs, err := queryBlock(bh)
					if err != nil {
						log.Printf("Error in querying Block: %d\t %v", bh, err)
						continue
					}
					fmt.Printf("Block: %s\tH: %d\tTxs: %d\t[Add missing]\n", rec.BlockHash, rec.Height, rec.NumOfTxs)

					dbRow := rec.getBlockDBRow()
					insertQueue.AddToInsertQueue(database.TABLE_BLOCKS, dbRow)

					// Adding the signers of the previous block
					for i := range rec.LastBlockSigners {
						// We insert them one by one, in case one fails due to e.g. duplication, the others will go through
						insertQueue.AddToInsertQueue(database.TABLE_BLOCK_SIGNERS, rec.LastBlockSigners[i].getBlockSignerDBRow())
					}

					// Insert TX hashes into the `tx_events`,
					// so the `tx.fixEmptyEvents` can pick them up, query them and fix them
					for i := range *txs {
						txHash := strings.ToUpper(hex.EncodeToString((*txs)[i].Hash()))
						// We insert them one by one, in case one fails due to e.g. duplication, the others will go through
						insertQueue.AddToInsertQueue(database.TABLE_TX_EVENTS, database.RowType{
							database.FIELD_TX_EVENTS_TX_HASH: txHash,
						})
					}

				}

				time.Sleep(time.Second)
			}
		}
	}()

}

func queryBlock(height uint64) (*BlockRecord, *tmTypes.Txs, error) {

	wsURI := os.Getenv("RPC_ADDRESS")

	heightStr := fmt.Sprintf("%d", height)

	// A dirty hack to get the things done
	cmd := exec.Command("archwayd", "query", "block", heightStr, "--node", wsURI) //, "--output", "json")
	stdout, err := cmd.Output()

	if err != nil {
		return nil, nil, err
	}

	// Since unmarshaling fails to the tmType.Block (Consensus.block.header.version.block)
	// We need to modify the JSON output to make it work
	stdoutStr := string(stdout)
	stdoutStr = regexp.MustCompile(`("block":)"([0-9]*?)"`).ReplaceAllString(stdoutStr, `$1$2`)
	stdoutStr = regexp.MustCompile(`("height":)"([0-9]*?)"`).ReplaceAllString(stdoutStr, `$1$2`)

	type tmBlock struct {
		tmTypes.BlockID `json:"block_id"`
		tmTypes.Block   `json:"block"`
	}
	var b tmBlock

	err = json.Unmarshal([]byte(stdoutStr), &b)
	if err != nil {
		return nil, nil, err
	}

	rec := getBlockRecordFromTmBlock(&b.Block)
	return rec, &b.Block.Txs, nil

}

func findMissingBlocks(start, end uint64, db *database.Database) ([]uint64, error) {
	var missingBlocks []uint64

	totalBlocks, err := GetTotalBlocksByRange(start, end, db)
	if err != nil {
		return missingBlocks, err
	}
	expectedBlocks := end - start + 1

	if totalBlocks != expectedBlocks {
		if start == end {
			missingBlocks = append(missingBlocks, start)
		} else {
			middle := (start + end) / 2
			mb1, err := findMissingBlocks(start, middle, db)
			if err != nil {
				return missingBlocks, err
			}
			missingBlocks = append(missingBlocks, mb1...)

			mb2, err := findMissingBlocks(middle+1, end, db)
			if err != nil {
				return missingBlocks, err
			}
			missingBlocks = append(missingBlocks, mb2...)
		}
	}

	return missingBlocks, nil
}

func GetTotalBlocksByRange(start, end uint64, db *database.Database) (uint64, error) {

	SQL := fmt.Sprintf(`
		SELECT 
			COUNT(*) AS "total"
		FROM "%s"
		WHERE 
			"height" >= $1 AND 
			"height" <= $2`,
		database.TABLE_BLOCKS,
	)

	rows, err := db.Query(SQL, database.QueryParams{start, end})
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 ||
		rows[0] == nil ||
		rows[0]["total"] == nil {
		return 0, nil
	}

	return uint64(rows[0]["total"].(int64)), nil
}

func GetLatestBlockHeight(db *database.Database) (uint64, error) {

	SQL := fmt.Sprintf(
		`SELECT MAX("%s") AS "result" FROM "%s"`,

		database.FIELD_BLOCKS_HEIGHT,
		database.TABLE_BLOCKS,
	)

	rows, err := db.Query(SQL, database.QueryParams{})
	if err != nil {
		return 0, err
	}

	if len(rows) == 0 ||
		rows[0] == nil ||
		rows[0]["result"] == nil {
		return 0, nil
	}

	return uint64(rows[0]["result"].(int64)), nil
}
