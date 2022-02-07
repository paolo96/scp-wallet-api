package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const CMCApiURL = "https://pro-api.coinmarketcap.com"
const CMCScpId = "4074"

var supportedFiats = []string{
	"EUR",
	"JPY",
	"GBP",
	"AUD",
	"CAD",
	"CHF",
	"CNY",
	"HKD",
	"NZD",
	"SEK",
	"INR",
}

var CMCApiKey = ""
var GetGeoApiKey = ""

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

	GetGetApiConversionResponse struct {
		Rates map[string]GetGetApiRate `json:"rates"`
	}

	GetGetApiRate struct {
		Rate string `json:"rate"`
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

//getUsdExchangeRates gets the USD to supportedFiats exchange rates from getgeoapi.com API
func getUsdExchangeRates() (*map[string]float64, error) {

	currencyList := strings.Join(supportedFiats, ",")

	client := http.Client{}
	request, err := http.NewRequest("GET", "https://api.getgeoapi.com/v2/currency/convert?api_key="+GetGeoApiKey+"&from=USD&to="+currencyList+"&format=json", nil)
	if err != nil {
		fmt.Println(err)
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	body, e := ioutil.ReadAll(response.Body)
	if verbose {
		fmt.Printf("GET getgeoapi RATES -> %v\n", string(body))
	}
	if e != nil {
		return nil, e
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {

		var data *GetGetApiConversionResponse
		e = json.Unmarshal(body, &data)
		if e != nil {
			return nil, e
		}

		var result = map[string]float64{}
		for currency, apiRate := range data.Rates {
			for _, supportedFiat := range supportedFiats {
				if supportedFiat == currency {
					exchangeRate, e := strconv.ParseFloat(apiRate.Rate, 8)
					if e == nil {
						result[currency] = exchangeRate
					}
					break
				}
			}
		}
		return &result, nil

	} else {
		return nil, errors.New("getgeoapi response code " + strconv.Itoa(response.StatusCode))
	}

}
