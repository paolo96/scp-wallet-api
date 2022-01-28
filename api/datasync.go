package main

import (
	"fmt"
	"scp-app-api/spdbridge"
	"time"
)

const (
	networkSyncInterval   = 10 * time.Second
	fiatPriceSyncInterval = 300 * time.Second

	networkSyncErrorInterval = 60 * time.Second
)

var networkData *NetworkData = nil
var usdPrice *float64 = nil

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

//StartDataSync starts the caching of the data
func StartDataSync() {

	go syncNetworkData(nil)
	go syncFiatQuote()

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

func syncFiatQuote() {

	newData, err := getScpUsdQuote()
	if err != nil {
		fmt.Printf("Error while fetching fiat price: %v\n", err)
		time.Sleep(fiatPriceSyncInterval)
		go syncFiatQuote()
		return
	}
	usdPrice = newData

	time.Sleep(fiatPriceSyncInterval)
	go syncFiatQuote()

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
