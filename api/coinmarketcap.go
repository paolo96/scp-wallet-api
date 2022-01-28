package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const CMCApiURL = "https://pro-api.coinmarketcap.com"
const CMCScpId = "4074"

var CMCApiKey = ""

type (
	CMCQuoteResponse struct {
		Quotes map[string]CMCQuote `json:"data"`
	}

	CMCQuote struct {
		Quote map[string]CMCQuoteData `json:"quote"`
	}

	CMCQuoteData struct {
		Price float64 `json:"price"`
	}
)

//getScpUsdQuote grabs the SCP/USD exchange rate from coinmarketcap API if an API key is provided
func getScpUsdQuote() (*float64, error) {

	if CMCApiKey == "" {
		return nil, errors.New("no API key provided for CMC")
	}

	req, e := http.NewRequest("GET", CMCApiURL+"/v1/cryptocurrency/quotes/latest?id="+CMCScpId, nil)
	if e != nil {
		return nil, e
	}

	req.Header.Set("X-CMC_PRO_API_KEY", CMCApiKey)

	client := &http.Client{}
	response, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer response.Body.Close()

	body, e := ioutil.ReadAll(response.Body)
	if verbose {
		fmt.Printf("GET CMC quote -> %v\n", string(body))
	}
	if e != nil {
		return nil, e
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {

		var data *CMCQuoteResponse
		e = json.Unmarshal(body, &data)
		if e != nil {
			return nil, e
		}
		for _, q := range data.Quotes {
			if val, ok := q.Quote["USD"]; ok {
				return &val.Price, nil
			}
		}

		return nil, errors.New("CMC quote not found in json")

	} else {
		return nil, errors.New("CMC response code " + strconv.Itoa(response.StatusCode))
	}

}
