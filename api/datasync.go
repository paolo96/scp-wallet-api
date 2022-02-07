package main

import (
	"fmt"
	"scp-app-api/spdbridge"
	"time"
)

const (
	networkSyncInterval          = 10 * time.Second
	usdPriceSyncInterval         = 300 * time.Second
	usdExchangeRatesSyncInterval = 1000 * time.Second

	networkSyncErrorInterval = 60 * time.Second
)

var networkData *NetworkData = nil
var usdPrice *float64 = nil
var exchangeRates *map[string]float64 = nil

//GetNetworkData returns the cached ScPrime network data
func GetNetworkData() (*NetworkData, error) {

	if networkData == nil {
		newData, err := downloadNetworkData()
		if err != nil {
			fmt.Printf("Error while fetching spd network data: %v\n", err)
			return nil, err
		}
		networkData = newData
	}
	return networkData, nil

}

//GetFiatPrice returns the cached SCP/USD exchange rate
func GetFiatPrice() (*float64, error) {

	if usdPrice == nil {
		newData, err := getScpUsdQuote()
		if err != nil {
			fmt.Printf("Error while fetching fiat price: %v\n", err)
			return nil, err
		}
		usdPrice = newData
	}
	return usdPrice, nil

}

//GetUsdExchangeRates returns the cached USD to supportedFiats exchange rate
func GetUsdExchangeRates() (*map[string]float64, error) {

	if exchangeRates == nil {
		newData, err := getUsdExchangeRates()
		if err != nil {
			fmt.Printf("Error while fetching usd exchange rates: %v\n", err)
			return nil, err
		}
		exchangeRates = newData
	}
	return exchangeRates, nil

}

//StartDataSync starts the caching of the data
func StartDataSync() {

	go syncNetworkData(nil)
	go syncUsdQuote()
	go syncUsdExchangeRates()

}

func syncNetworkData(changedHeight *func(uint64, uint64)) {

	newData, err := downloadNetworkData()
	if err != nil {
		if verbose {
			fmt.Println("Waiting for daemon to resync")
		}
		time.Sleep(networkSyncErrorInterval)
		go syncNetworkData(changedHeight)
		return
	}

	var oldHeight uint64
	if networkData != nil {
		oldHeight = networkData.ConsensusHeight
	}
	networkData = newData
	if oldHeight < newData.ConsensusHeight && changedHeight != nil {
		(*changedHeight)(oldHeight, newData.ConsensusHeight)
	}

	time.Sleep(networkSyncInterval)
	go syncNetworkData(changedHeight)

}

func syncUsdQuote() {

	newData, err := getScpUsdQuote()
	if err != nil {
		fmt.Printf("Error while fetching fiat price: %v\n", err)
		time.Sleep(usdPriceSyncInterval)
		go syncUsdQuote()
		return
	}
	usdPrice = newData

	time.Sleep(usdPriceSyncInterval)
	go syncUsdQuote()

}

func syncUsdExchangeRates() {

	newData, err := getUsdExchangeRates()
	if err != nil {
		fmt.Printf("Error while fetching usd exchange rates: %v\n", err)
		time.Sleep(usdExchangeRatesSyncInterval)
		go syncUsdExchangeRates()
		return
	}
	exchangeRates = newData

	time.Sleep(usdExchangeRatesSyncInterval)
	go syncUsdExchangeRates()

}

//downloadNetworkData downloads and aggregates the following data from spd:
//consensus height, min transaction fee, max transaction fee
func downloadNetworkData() (*NetworkData, error) {

	consensus, err := spdbridge.GetConsensus()
	if err != nil || !consensus.Synced {
		return nil, err
	}

	fees, err := spdbridge.GetTransactionPoolFees()
	if err != nil {
		return nil, err
	}

	var newData = NetworkData{
		ConsensusHeight: consensus.Height,
		MinFee:          fees.MinFee,
		MaxFee:          fees.MaxFee,
	}

	if verbose {
		fmt.Printf("Successfully retrieved network data %v\n", newData)
	}
	return &newData, nil

}
