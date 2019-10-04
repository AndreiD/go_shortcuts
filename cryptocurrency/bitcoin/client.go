package btc

import (
	"watcher/config"

	rpcclient "github.com/stevenroose/go-bitcoin-core-rpc"

	log "github.com/sirupsen/logrus"
)

// BTCClient ..
var BTCClient *rpcclient.Client

// InitBTCClient ..
func InitBTCClient(conf *config.Config) error {

	log.Printf("BTC Client trying to connect to %s", conf.BTCMainnetRPC)

	connCfg := &rpcclient.ConnConfig{
		Host: conf.BTCMainnetRPC,
		User: "yyyyyy",
		Pass: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	BTCClient, err := rpcclient.New(connCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer BTCClient.Shutdown()

	// Get the current block count.
	blockCount, err := BTCClient.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BTC Connected ok to %s having block height: %d", conf.BTCMainnetRPC, blockCount)

	return nil
}