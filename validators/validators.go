package validators

import (
	"context"
	"fmt"
	"time"

	"github.com/archway-network/cosmologger/configs"
	"github.com/archway-network/cosmologger/database"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
)

func queryValidatorInfoByValAddr(grpcCnn *grpc.ClientConn, valAddr string) (staking.Validator, error) {

	var err error
	var response *staking.QueryValidatorResponse

	for retry := 0; retry < configs.Configs.GRPC.APICallRetry; retry++ {

		c := staking.NewQueryClient(grpcCnn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(configs.Configs.GRPC.CallTimeout))
		defer cancel()

		response, err = c.Validator(ctx,
			&staking.QueryValidatorRequest{
				ValidatorAddr: valAddr,
			})
		if err != nil {
			fmt.Printf("\r[%d", retry+1)
			// fmt.Printf("\r\tRetrying [ %d ]...", retry+1)
			// fmt.Printf("\tErr: %s", err)

			// Ideally we want to retry after getting 502 http error, because sometimes server returns it
			// but we cannot have it as the protobuf Invoke does not return the status code
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return response.Validator, nil

	}
	return staking.Validator{}, err
}

// This function retrieves the consensus address from the consensus public key
func GetConsAddressFromConsPubKey(inKey []byte) string {

	// For some unknown reasons there are two extra bytes in the begining of the key
	// which cause the size error, so we remove them
	//TODO: has to be fixed with legacy unmarshal

	pubkey := &ed25519.PubKey{Key: inKey[2:]}
	return sdk.ConsAddress(pubkey.Address().Bytes()).String()
}

func DoesConsAddrExistInDB(db *database.Database, valAddr string) (bool, error) {

	SQL := fmt.Sprintf(`
			SELECT 
				COUNT(*) as "total"
			FROM "%s" 
			WHERE "%s" = $1`,
		database.TABLE_VALIDATORS,
		database.FIELD_VALIDATORS_OPR_ADDR,
	)

	rows, err := db.Query(SQL, database.QueryParams{valAddr})
	if err != nil {
		return false, err
	}

	return rows[0]["total"].(int64) > 0, nil
}

func AddNewValidator(db *database.Database, grpcCnn *grpc.ClientConn, valAddr string) error {

	exist, err := DoesConsAddrExistInDB(db, valAddr)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	vInfo, err := queryValidatorInfoByValAddr(grpcCnn, valAddr)
	if err != nil {
		return err
	}
	consAddr := GetConsAddressFromConsPubKey(vInfo.ConsensusPubkey.Value)
	moniker := vInfo.Description.Moniker

	sdkValAddr, err := sdk.ValAddressFromBech32(valAddr)
	if err != nil {
		return err
	}
	accountAddr := sdk.AccAddress(sdkValAddr.Bytes()).String()

	rec := ValidatorRecord{
		ConsAddr:    consAddr,
		OprAddr:     valAddr,
		AccountAddr: accountAddr,
		Moniker:     moniker,
	}

	dbRow := rec.getDBRow()
	_, err = db.Insert(database.TABLE_VALIDATORS, dbRow)
	return err
}

func (v ValidatorRecord) getDBRow() database.RowType {
	return database.RowType{

		database.FIELD_VALIDATORS_CONS_ADDR:    v.ConsAddr,
		database.FIELD_VALIDATORS_OPR_ADDR:     v.OprAddr,
		database.FIELD_VALIDATORS_ACCOUNT_ADDR: v.AccountAddr,
		database.FIELD_VALIDATORS_MONIKER:      v.Moniker,
	}
}

func queryValidatorsSetByOffset(conn *grpc.ClientConn, offset int) (response *staking.QueryValidatorsResponse, err error) {

	for retry := 0; retry < configs.Configs.GRPC.APICallRetry; retry++ {

		c := staking.NewQueryClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(configs.Configs.GRPC.CallTimeout))
		defer cancel()

		response, err := c.Validators(ctx,
			&staking.QueryValidatorsRequest{
				// Status: status,
				Pagination: &query.PageRequest{
					// Key:    nextKey,
					// Limit:  limit,
					Offset: uint64(offset),
					// Reverse: false,
				},
			})
		if err != nil {
			// fmt.Printf("\r[%d", retry+1)
			fmt.Printf("\n[%d", retry+1)
			// fmt.Printf("\r\tRetrying [ %d ]...", retry+1)
			fmt.Printf("\tErr: %s", err)

			// Ideally we want to retry after getting 502 http error, because sometimes server returns it
			// but we cannot have it as the protobuf Invoke does not return the status code
			time.Sleep(50 * time.Millisecond)
			continue
		}

		return response, nil
	}

	return nil, err
}

func QueryValidatorsList(grpcCnn *grpc.ClientConn) ([]string, error) {

	var validatorsList []string

	offset := 0
	for {
		response, err := queryValidatorsSetByOffset(grpcCnn, offset)
		if err != nil {
			return validatorsList, err
		}

		if response == nil || len(response.Validators) == 0 {
			break
		}
		offset += len(response.Validators)

		for i := range response.Validators {
			validatorsList = append(validatorsList, response.Validators[i].OperatorAddress)
		}
	}

	return validatorsList, nil
}
