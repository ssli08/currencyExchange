package currency

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	vault "github.com/hashicorp/vault/api"
)

type Body struct {
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     rate   `json:"rates"`
}
type rate struct {
	AUD float32 `json:"AUD,omitempty"`
	NZD float32 `json:"NZD,omitempty"`
	CNY float32 `json:"CNY,omitempty"`
	USD float32 `json:"USD,omitempty"`
}

// get value of CNY/NZD and USD/NZD
func GetCurrencyRates(url, apiKey string) (Body, error) {
	/* {
	    "success": true,
	    "timestamp": 1519296206,
	    "base": "EUR",
	    "date": "2021-03-17",
	    "rates": {
	        "AUD": 1.566015,
	        "CAD": 1.560132,
	        "CHF": 1.154727,
	        "CNY": 7.827874,
	        "GBP": 0.882047,
	        "JPY": 132.360679,
	        "USD": 1.23396,
	    [...]
	    }
	} */

	// currency exchange api url
	url = fmt.Sprintf("%s=%s", url, apiKey)
	data, err := HttpProcess(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		return Body{}, fmt.Errorf("get data failed %s", err)
	}

	d := Body{}
	if err := json.Unmarshal(data, &d); err != nil {
		return Body{}, err
	}

	// return fmt.Sprintf("%v == %f(USD/NZD: %f)", time.Unix(d.Timestamp, 0), d.Rates.CNY/d.Rates.NZD, d.Rates.USD/d.Rates.NZD), nil
	return d, nil
}

// put token to vault server
func PutTokenToVault(client *vault.Client, secret string) error {

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	secretData := map[string]interface{}{"myPassword": secret}
	if _, err := client.KVv2("secret").Put(context.Background(), "exchangeAPI", secretData); err != nil {
		return err
	}
	log.Println("apiKey written successfully  to the vault")

	return nil
}

// get token from
func GetTokenFromVault(client *vault.Client) (string, error) {
	secret, err := client.KVv2("secret").Get(context.Background(), "exchangeAPI")
	if err != nil {
		return "", err
	}
	value, ok := secret.Data["myPassword"].(string)
	if !ok {
		return "", err
	}
	return value, nil
}
