package main

import (
	"log"
	"testing"
)

//Provide a valid getgeo api key to run TestGetFiatExchangeRates
const GetGeoApiKeyTest = ""

func TestGetFiatExchangeRates(t *testing.T) {

	GetGeoApiKey = GetGeoApiKeyTest

	response, e := getUsdExchangeRates()
	if e != nil {
		log.Fatal(e)
	}

fiatsLoop:
	for _, currency := range supportedFiats {
		for rateCurrency := range response.Rates {
			if currency == rateCurrency {
				continue fiatsLoop
			}
		}
		log.Fatal("Currency " + currency + " not found in response")
	}

}
