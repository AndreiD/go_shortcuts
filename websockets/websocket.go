package bitstamp

import (
	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strings"
	"time"
)

var websocketConnection *websocket.Conn

// StartWS .
func StartWS() {

	u := url.URL{Scheme: "wss", Host: "ws.bitstamp.net", Path: ""}
	log.Printf("connecting....... %s", u.String())
	var err error
	websocketConnection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error("dial: %s", err)
		return
	}
	defer utils.Close(websocketConnection)
	log.Infof("connected OK..... %s", u.String())
	done := make(chan struct{})
	// listen for messages
	go func() {
		defer close(done)
		for {
			_, rawMessage, err := websocketConnection.ReadMessage()
			if err != nil {
				log.Error("error:", err)
				log.Warn("trying to reconnect...")
				time.Sleep(1 * time.Second)
				StartWS()
			}

			// first message is the success subscription
			if !strings.Contains(string(rawMessage), "subscription_succeeded") {
				updateBestBidAsk(rawMessage)
			}
		}
	}()

	// subscribe to trade data
	msg := `{"event": "bts:subscribe","data": {"channel": "order_book_btcusd"}}`
	if err := websocketConnection.WriteMessage(websocket.BinaryMessage, []byte(msg)); err != nil {
		log.Println(err)
		return
	}

	select {}
}

func updateBestBidAsk(data []byte) {
	var err error

	_bidPrice, _, _, err := jsonparser.Get(data, "data", "bids", "[0]", "[0]")
	if err != nil {
		log.Errorf("error parsing the best bid %s", err.Error())
		return
	}

	// etc .....
}
