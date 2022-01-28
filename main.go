package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/archway-network/cosmologger/configs"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	// conn, err := Connect()
	// if err != nil {
	// 	log.Fatalf("Did not connect: %s", err)
	// }
	// defer conn.Close()

	/*-------------*/

	// SetBech32Prefixes()

	/*-------------*/

	// c := tx.NewServiceClient()
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(configs.Configs.GRPC.CallTimeout))
	// defer cancel()

	// response, err := c.

	// client.EventsClient.Subscribe()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(configs.Configs.GRPC.CallTimeout))
	defer cancel()

	// defaultTMURI := "tcp://rpc.cosmos.network:443"
	// defaultTMURI := "tcp://localhost:26657"
	defaultTMURI := "tcp://192.168.188.26:26657"

	fmt.Println("Connecting to the thing...")
	cli, err := http.New(defaultTMURI, "/websocket")
	if err != nil {
		panic(err)
	}
	fmt.Println(cli)

	fmt.Println("Starting the client...")
	err = cli.Start()
	if err != nil {
		panic(err)
	}

	// Make a closure so you can actually close this when you finish everything
	stopClient := func() {
		err := cli.UnsubscribeAll(ctx, "helpers")
		if err != nil {
			panic(err)
		}
		err = cli.Stop()
		if err != nil {
			panic(err)
		}
	}

	// eventChan, err := cli.Subscribe(ctx, "helpers", tmTypes.QueryForEvent(tmTypes.EventNewBlock).String())
	eventChan, err := cli.Subscribe(ctx, "helpers", tmTypes.QueryForEvent(tmTypes.EventTx).String())
	if err != nil {
		panic(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Kill, //nolint
		os.Interrupt)

	for {
		fmt.Println("Waiting for the new signal")
		select {
		case evRes := <-eventChan:
			processTx(evRes)
		case <-signalChannel:
			stopClient()
		}
	}

}

func processTx(evr coretypes.ResultEvent) {

	fmt.Printf("\n\n\t\t\t=================\n\n\n")

	jsonString, err := json.MarshalIndent(evr, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("jsonString: %s\n", jsonString)

	tx := evr.Data.(tmTypes.EventDataTx)
	// tx, err :=
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Printf("\n\nTX: %+v\n", tx)

	// fmt.Printf("evr.Data.(tmTypes.EventDataTx).Result: %v\n\n-----------\n", evr.Data.(tmTypes.EventDataTx).Result)
	// fmt.Printf("evr.Data.(tmTypes.EventDataTx).TxResult: %v\n", evr.Data.(tmTypes.EventDataTx).TxResult)

	// fmt.Printf("\n\n\t\t\t=----------------\n\n\n")

	// fmt.Printf("evr.Events: \n")
	// for i := range evr.Events {
	// 	fmt.Printf("\n\nType: %v\n", i)
	// 	for j := range evr.Events[i] {
	// 		fmt.Printf("\tAttr: %v ==> %s\n", j, evr.Events[i][j])
	// 	}
	// }
}

func Connect() (*grpc.ClientConn, error) {

	if configs.Configs.GRPC.TLS {
		creds := credentials.NewTLS(&tls.Config{})
		return grpc.Dial(configs.Configs.GRPC.Server, grpc.WithTransportCredentials(creds))
	}
	return grpc.Dial(configs.Configs.GRPC.Server, grpc.WithInsecure())

}

func SetBech32Prefixes() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(configs.Configs.Bech32Prefix.Account.Address, configs.Configs.Bech32Prefix.Account.PubKey)
	config.SetBech32PrefixForValidator(configs.Configs.Bech32Prefix.Validator.Address, configs.Configs.Bech32Prefix.Validator.PubKey)
	config.SetBech32PrefixForConsensusNode(configs.Configs.Bech32Prefix.Consensus.Address, configs.Configs.Bech32Prefix.Consensus.PubKey)
	config.Seal()
}
