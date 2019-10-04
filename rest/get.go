package rest

import (
	"github.com/go-delve/delve/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// GetRequest executes a GET request
func GetRequest(url string) ([]byte, error) {

	client := http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	reqWithDeadline := req.WithContext(ctx)
	response, err := client.Do(reqWithDeadline)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	return data, err

}

// GetRequestWithStuff .
func GetEntrustStatus(entrustID string, conf *config.Config) {
	url := baseURL + "/exchange/entrust/controller/website/EntrustController/getEntrustById?marketId=4040&entrustId=" + entrustID

	client := http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
		return
	}
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	req.Header.Set("Timestamp", timestamp)
	req.Header.Set("Apiid", conf.BWAPIAccessKey)
	req.Header.Set("Sign", utils.GetMD5Hash(conf.BWAPIAccessKey+timestamp+"entrustId"+entrustID+"marketId4040"+conf.BWAPISecretKEY))

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

	payload, _, _, err := jsonparser.Get(data, "datas")
	if err != nil {
		log.Error(err)
		return
	}

	// Status
	// -2 funds unfrozen failed -1 user funds are insufficient 0 start 1 cancel 2 transaction success 3 transaction part 4 cancel
	log.Println("--------------------------")
	log.Println("Entrust Status -> %s", string(payload))
	log.Println("--------------------------")
}
