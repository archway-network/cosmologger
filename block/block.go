package block

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/archway-network/cosmologger/configs"
	"github.com/archway-network/cosmologger/database"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

func ProcessEvents(db *database.Database, evr *coretypes.ResultEvent) error {

	rec := getBlockRecordFromEvent(evr)
	fmt.Printf("Block: %s\n", rec.BlockHash)

	dbRow := rec.getBlockDBRow()
	db.InsertAsync(database.TABLE_BLOCKS, dbRow)
	// _, err := db.Insert(database.TABLE_BLOCKS, dbRow)
	// if err != nil {
	// 	return err
	// }

	for i := range rec.Signers {

		dbRow := rec.Signers[i].getBlockSignerDBRow()
		db.InsertAsync(database.TABLE_BLOCK_SIGNERS, dbRow)
		// _, err := db.Insert(database.TABLE_BLOCK_SIGNERS, dbRow)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func getBlockRecordFromEvent(evr *coretypes.ResultEvent) *BlockRecord {
	var br BlockRecord

	b := evr.Data.(tmTypes.EventDataNewBlock)
	br.BlockHash = b.Block.Hash().String()

	br.Height = uint64(b.Block.Height)
	br.NumOfTxs = uint64(len(b.Block.Txs))
	br.Time = b.Block.Time

	for i := range b.Block.LastCommit.Signatures {

		consAddr, err := sdk.ConsAddressFromHex(b.Block.LastCommit.Signatures[i].ValidatorAddress.String())
		if err != nil {
			continue // just ignore this signer as it might not be running and we face some strange error
		}

		br.Signers = append(br.Signers, BlockSignersRecord{
			BlockHeight: br.Height,
			ValConsAddr: consAddr.String(),
			Time:        b.Block.LastCommit.Signatures[i].Timestamp,
			Signature:   base64.StdEncoding.EncodeToString(b.Block.LastCommit.Signatures[i].Signature),
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

func Start(cli *tmClient.HTTP, db *database.Database) {

	go func() {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(configs.Configs.GRPC.CallTimeout))
		defer cancel()

		eventChan, err := cli.Subscribe(ctx, configs.Configs.SubscriberName, tmTypes.QueryForEvent(tmTypes.EventNewBlock).String())
		if err != nil {
			panic(err)
		}

		for {
			evRes := <-eventChan
			err := ProcessEvents(db, &evRes)
			if err != nil {
				//TODO: We need some customizable log level
				log.Printf("Error in processing block event: %v", err)
			}
		}
	}()
}
