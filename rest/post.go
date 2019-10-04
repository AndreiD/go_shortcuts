package bitstamp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// POSTExample .... COMPLETE POST EXAMPLE .
func POSTExample(config *configs.ViperConfiguration) {

	baseURL := "www.bitstamp.net"
	urlPath := "/api/v2/balance/"

	apiKey := config.Get("bitstamp.apiKey")
	secret := config.Get("bitstamp.secret")
	randUUID, _ := uuid.NewRandom()
	xAuthNonce := randUUID.String()
	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)

	stringToSign := "BITSTAMP" + " " + apiKey + "POST" + baseURL + urlPath + "" + xAuthNonce + timestamp + "v2"

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	sha := hex.EncodeToString(h.Sum(nil))
	sha = strings.ToUpper(sha)

	client := http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest("POST", "https://"+baseURL+urlPath, nil)
	if err != nil {
		log.Error(err)
		return
	}

	req.Header.Set("X-Auth", "BITSTAMP "+apiKey)
	req.Header.Set("X-Auth-Signature", sha)
	req.Header.Set("X-Auth-Nonce", xAuthNonce)
	req.Header.Set("X-Auth-Timestamp", timestamp)
	req.Header.Set("X-Auth-Version", "v2")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	reqWithDeadline := req.WithContext(ctx)
	response, err := client.Do(reqWithDeadline)
	if err != nil {
		log.Error(err)
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("Server responded with %s", string(data))

	// --------- BALANCE USD ---------
	balanceUSD, _, _, err := jsonparser.Get(data, "usd_available")
	if err != nil {
		log.Error(err)
		return
	}

}
