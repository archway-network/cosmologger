package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/archway-network/cosmologger/block"
	"github.com/archway-network/cosmologger/configs"
	"github.com/archway-network/cosmologger/database"
	"github.com/archway-network/cosmologger/dbinit"
	"github.com/archway-network/cosmologger/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	tmClient "github.com/tendermint/tendermint/rpc/client/http"
)

/*--------------*/
func main() {

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

	/*-------------*/

	SetBech32Prefixes()

	/*-------------*/

	fmt.Printf("\nConnecting to the websocket... ")

	wsURI := os.Getenv("RPC_ADDRESS")

	wsURI = "tcp://192.168.188.26:26657"
	// wsURI = "ws://65.21.229.173:26657"
	// wsURI = "wss://rpc.cosmos.network:443"
	// wsURI = "https://rpc.augusta-1.archway.tech:443"

	/*-----------------------*/

	// creds := credentials.NewTLS(&tls.Config{})

	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{
	// 			InsecureSkipVerify: true,
	// 			// ClientAuth:         tls.VerifyClientCertIfGiven,
	// 			// Certificates:       []tls.Certificate{},
	// 		},
	// 	},
	// }
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

	// Due to some limitations of the RPC APIs we need to call GRPC ones as well
	grpcCnn, err := GrpcConnect()
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	defer grpcCnn.Close()

	/*------------------*/

	fmt.Println("Listening...")
	// Running the listeners
	tx.Start(cli, grpcCnn, db)
	block.Start(cli, grpcCnn, db)

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
