package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tmClient "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/archway-network/cosmologger/block"
	"github.com/archway-network/cosmologger/configs"
	"github.com/archway-network/cosmologger/database"
	"github.com/archway-network/cosmologger/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*--------------*/

// // Parse URL and set defaults
// func newParsedURL(remoteAddr string) (*parsedURL, error) {
// 	u, err := url.Parse(remoteAddr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// default to tcp if nothing specified
// 	if u.Scheme == "" {
// 		u.Scheme = protoTCP
// 	}

// 	pu := &parsedURL{
// 		URL:          *u,
// 		isUnixSocket: false,
// 	}

// 	if u.Scheme == protoUNIX {
// 		pu.isUnixSocket = true
// 	}

// 	return pu, nil
// }

// func makeHTTPDialer(remoteAddr string) (func(string, string) (net.Conn, error), error) {
// 	u, err := newParsedURL(remoteAddr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	protocol := u.Scheme

// 	// accept http(s) as an alias for tcp
// 	switch protocol {
// 	case protoHTTP, protoHTTPS:
// 		protocol = protoTCP
// 	}

// 	dialFn := func(proto, addr string) (net.Conn, error) {
// 		return net.Dial(protocol, u.GetDialAddress())
// 	}

// 	return dialFn, nil
// }

func main() {

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db := database.New(database.Postgres, psqlconn)
	defer db.Close()

	// conn, err := Connect()
	// if err != nil {
	// 	log.Fatalf("Did not connect: %s", err)
	// }
	// defer conn.Close()

	/*-------------*/

	SetBech32Prefixes()

	/*-------------*/

	// defaultTMURI := "tcp://rpc.cosmos.network:443"
	// defaultTMURI := "tcp://localhost:26657"
	// defaultTMURI := "tcp://192.168.188.26:26657"
	// defaultTMURI := "http://192.168.188.26:26657"
	// defaultTMURI := "https://65.21.229.173:443"
	// defaultTMURI := "ws://65.21.229.173:26657"
	// defaultTMURI := "tcp://35.196.115.108:31306" // Constantine
	// defaultTMURI := "tcp://rpc.cosmos.network:443"
	defaultTMURI := "https://rpc.cosmos.network:443"

	fmt.Println("Connecting to the websocket...")

	/*-----------------------*/

	// creds := credentials.NewTLS(&tls.Config{})

	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{
	// 			// InsecureSkipVerify: false,
	// 			ClientAuth: tls.VerifyClientCertIfGiven,
	// 		},
	// 	},
	// }
	// cli, err := tmClient.NewWithClient(defaultTMURI, "/websocket", client)
	cli, err := tmClient.New(defaultTMURI, "/websocket")
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the client...")

	if err := cli.Start(); err != nil {
		panic(err)
	}

	fmt.Println("Listening...")
	// Running the listeners
	tx.Start(cli, db)
	block.Start(cli, db)

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

	if err := cli.UnsubscribeAll(context.Background(), configs.Configs.SubscriberName); err != nil {
		panic(err)
	}
	if err := cli.Stop(); err != nil {
		panic(err)
	}

	fmt.Println("Ciao bello!")
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
