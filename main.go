package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// "github.com/archway-network/archway/app/params"
	// "github.com/archway-network/archway/app"
	"github.com/archway-network/cosmologger/block"
	"github.com/archway-network/cosmologger/configs"
	"github.com/archway-network/cosmologger/database"
	"github.com/archway-network/cosmologger/dbinit"
	"github.com/archway-network/cosmologger/tx"

	// "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	tmClient "github.com/tendermint/tendermint/rpc/client/http"
	// tmClient "github.com/tendermint/tendermint/rpc/client"
)

/*--------------*/
func main() {

	// str := "{\"developer_address\":\"archway1p69nfnrn9ju8ghqn4pcv7zqk2tpcsdwnganvnn\",\"reward_address\":\"archway1p69nfnrn9ju8ghqn4pcv7zqk2tpcsdwnganvnn\",\"gas_rebate_to_user\":false,\"collect_premium\":true,\"premium_percentage_charged\":\"70\"}"
	// str := "[{\"denom\":\"stake\",\"amount\":\"0.001988500000000000\"}]"
	// var metadataX []map[string]interface{}
	// if err := json.Unmarshal([]byte(str), &metadataX); err != nil {
	// 	panic(err)
	// }
	// metadata := metadataX[0]

	// fmt.Printf("metadata[\"developer_address\"]: %s\n", metadata["developer_address"])
	// fmt.Printf("metadata[\"developer_address\"]: %s\n", metadata["Gooz"])

	// for i := range metadata {

	// 	fmt.Printf("\n%#v ==> %#v\n", i, metadata[i])
	// }

	// numValue, err := strconv.ParseFloat(metadata["amount"].(string), 64)
	// if err != nil {
	// 	fmt.Printf("\nError in Unmarshaling '%s': %v\n", "amount", err)
	// }

	// fmt.Printf("numValue: %#v\n", numValue)

	// return

	/*-------------*/

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	fmt.Printf("\nConnecting to the Database... ")

	db := database.New(database.Postgres, psqlconn)
	defer db.Close()

	// Check if we need to create tables and stuff on the DB
	dbinit.DatabaseInit(db)

	fmt.Printf("\nDone")

	insertQueue := database.NewInsertQueue(db)
	if err := insertQueue.Start(); err != nil {
		fmt.Printf("error in starting insert queue: %v\n", err)
		return
	}
	defer insertQueue.Stop()

	/*-------------*/

	SetBech32Prefixes()

	/*-------------*/

	wsURI := os.Getenv("RPC_ADDRESS")

	fmt.Printf("\nConnecting to the RPC [%s]... ", wsURI)

	//TODO: There is a known issue with the TM client when we use TLS
	// cli, err := tmClient.NewWithClient(wsURI, "/websocket", client)
	cli, err := tmClient.New(wsURI, "/websocket")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Done")

	/*------------------*/

	fmt.Printf("\nStarting the client...\n")

	var cliErr error
	for i := 1; i <= configs.Configs.TendermintClient.ConnectRetry; i++ {

		fmt.Printf("\tTrial #%d\n", i)
		cliErr = cli.Start()
		if cliErr == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if cliErr != nil {
		panic(cliErr)
	}

	fmt.Println("\nDone")

	/*------------------*/

	// txHash := "A7E403D4B07A1C0D969DDE2560D306FC161650FF129B86382E213313F5757818"
	// query := fmt.Sprintf("tx.hash='%s'", txHash)

	// cliCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(configs.Configs.GRPC.CallTimeout))
	// defer cancel()

	// // hashByte, err := hex.DecodeString(txHash)
	// // cli.Tx(hashByte, false)

	// // res, err := cli.TxSearch(cliCtx, query, true, nil, nil, "")
	// // fmt.Printf("res: %#v\n", res.Txs[0])
	// // tx := *sdk.TxDecoder(res.Txs[0].TxResult.Data)

	// // tx, err := sdk.TxDecoder(res.Txs[0].Tx.)
	// // fmt.Printf("\nTX: %+v\n", tx.TxBody)

	// // qClient := authx.NewQueryClient(&grpc.ClientConn{})
	// // fmt.Printf("qClient: %v\n", qClient)

	// // cliCtx

	// // cliCtx := sdkClient.conte

	// // res, err := authx.QueryTx(cliCtx, txHash)

	// tcli, err := sdkClient.NewClientFromNode(wsURI)
	// fmt.Printf("err: %v\n", err)

	// res, err := tcli.TxSearch(cliCtx, query, true, nil, nil, "")
	// fmt.Printf("err: %v\n", err)
	// // fmt.Printf("res: %#v\n", res.Txs[0].Tx)

	// /*------------*/

	// // encodingConfig := params.MakeEncodingConfig()
	// encodingConfig := MakeEncodingConfig()
	// // encodingConfig.TxConfig.TxJSONDecoder()

	// txb, err := encodingConfig.TxConfig.TxDecoder()(res.Txs[0].Tx)
	// // txb, err := sdk.TxDecoder(res.Txs[0].Tx)

	// // var cdc *codec.LegacyAmino

	// // txb := legacytx.StdTx{}
	// // err = cdc.Unmarshal(res.Txs[0].Tx, &txb)

	// fmt.Printf("\n========\ntxb: %v\n\n========\n", txb)

	// clientCtx := sdkClient.Context{
	// 	// NodeURI: wsURI,
	// 	// ChainID: "torii-1",
	// 	Client:   cli,
	// 	TxConfig: encodingConfig.TxConfig,
	// }

	// // // fmt.Printf("\n========\nclientCtx: %+v\n", clientCtx)
	// output, err := authtx.QueryTx(clientCtx, txHash)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("\n========\noutput: %+v\n\n========\n", output)

	// if output.Empty() {
	// 	panic(fmt.Errorf("no transaction found with hash %s", txHash))
	// }

	// fmt.Printf("\n========\n%+v\n", clientCtx.PrintProto(output))

	// // A dirty hack to get the things done
	// cmd := exec.Command("archwayd", "query", "tx", txHash, "--node", wsURI, "--output", "json")
	// stdout, err := cmd.Output()

	// rec := getTxRecordFromJson(string(stdout))

	// js, _ := json.MarshalIndent(rec, "", "  ")

	// fmt.Printf("\n-------------------------\n\nREC: %s\n", js)

	// panic(err)

	/*------------------*/

	// Due to some limitations of the RPC APIs we need to call GRPC ones as well
	grpcCnn, err := GrpcConnect()
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
		return
	}
	defer grpcCnn.Close()

	/*------------------*/

	fmt.Println("\nListening...")
	// Running the listeners
	tx.Start(cli, grpcCnn, db, insertQueue)
	// tx.FixEmptyEvents(cli, grpcCnn, db)
	block.Start(cli, grpcCnn, db, insertQueue)

	/*------------------*/

	// Exit gracefully
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Kill, //nolint
		os.Interrupt)
	<-quitChannel

	//Time for cleanup before exit

	if err := cli.UnsubscribeAll(context.Background(), configs.Configs.TendermintClient.SubscriberName); err != nil {
		panic(err)
	}
	if err := cli.Stop(); err != nil {
		panic(err)
	}

	fmt.Println("\nCiao bello!")
}

func GrpcConnect() (*grpc.ClientConn, error) {

	tlsEnabled := os.Getenv("GRPC_TLS")
	GRPCServer := os.Getenv("GRPC_ADDRESS")

	fmt.Printf("\nConnecting to the GRPC [%s] \tTLS: [%s]", GRPCServer, tlsEnabled)

	if strings.ToLower(tlsEnabled) == "true" {
		creds := credentials.NewTLS(&tls.Config{})
		return grpc.Dial(GRPCServer, grpc.WithTransportCredentials(creds))
	}
	return grpc.Dial(GRPCServer, grpc.WithInsecure())

}

func SetBech32Prefixes() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(configs.Configs.Bech32Prefix.Account.Address, configs.Configs.Bech32Prefix.Account.PubKey)
	config.SetBech32PrefixForValidator(configs.Configs.Bech32Prefix.Validator.Address, configs.Configs.Bech32Prefix.Validator.PubKey)
	config.SetBech32PrefixForConsensusNode(configs.Configs.Bech32Prefix.Consensus.Address, configs.Configs.Bech32Prefix.Consensus.PubKey)
	config.Seal()
}

// // MakeEncodingConfig creates a new EncodingConfig with all modules registered
// func MakeEncodingConfig() params.EncodingConfig {
// 	encodingConfig := params.MakeEncodingConfig()
// 	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
// 	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
// 	// ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
// 	// ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
// 	return encodingConfig
// }
