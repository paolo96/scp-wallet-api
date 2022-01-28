package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"scp-app-api/spdbridge"
)

const standardFailResponse = "{\"status\":\"ko\"}"
const standardSuccessResponse = "{\"status\":\"ok\"}"

//getScPrimeDataHandler handles requests to /scprime/data
//Returns the cached network data and the cached SCP/USD exchange rate
func getScPrimeDataHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	data, _ := GetNetworkData()
	usdPrice, _ := GetFiatPrice()

	jsonResp, err := json.Marshal(NetworkDataResponse{
		NetworkData: data,
		ScpPrice:    usdPrice,
	})
	if err != nil {
		http.Error(w, standardFailResponse, 500)
		return
	}

	w.Write(jsonResp)
}

//getTransactionsHandler handles requests to /transactions/batch
//Returns the transactions related to the addresses requested
func getAddressesTransactionsBatchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, standardFailResponse, 400)
		return
	}

	var params TransactionsBatchParams
	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, standardFailResponse, 400)
		return
	}

	explorerAddresses, err := spdbridge.ExplorerAddressesBatch(params.Addresses)
	if err != nil {
		http.Error(w, standardFailResponse, 500)
		return
	}

	unconfirmedTransactions, err := spdbridge.GetTransactionPool()
	if err != nil {
		http.Error(w, standardFailResponse, 500)
		return
	}

	transactions := filterTransactions(params, explorerAddresses, unconfirmedTransactions)
	jsonResp, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, standardFailResponse, 500)
		return
	}

	w.Write(jsonResp)
}

//getTransactionsHandler handles requests to /transactions
//Verifies and then broadcasts the transaction set provided
//The data for the verification is sent separately from the data for the broadcast. This is done because they
//need different formatting, and it's quicker for the client App to provide it rather than reformat it server side
func newTransactionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, standardFailResponse, 400)
		return
	}

	var newTransaction NewTransactionParams
	err = json.Unmarshal(body, &newTransaction)
	if err != nil {
		http.Error(w, standardFailResponse, 400)
		return
	}

	result, err := spdbridge.ConsensusValidateTxns([]byte(newTransaction.ValidateData))
	if err != nil || !result {
		http.Error(w, standardFailResponse, 400)
		return
	}

	resultBroadcast, err := spdbridge.TransactionPoolRaw(newTransaction.BroadcastData.Parents, newTransaction.BroadcastData.Transaction)
	if err != nil || !resultBroadcast {
		http.Error(w, standardFailResponse, 400)
		return
	}

	fmt.Fprintf(w, standardSuccessResponse)
}

func filterTransactions(params TransactionsBatchParams, explorerAddresses *spdbridge.AddressesBatchResp, unconfirmedTransactions *spdbridge.TransactionPoolResp) (transactions TransactionsBatchResp) {

	for _, explorerAddress := range explorerAddresses.Addresses {
		for _, explorerTransaction := range explorerAddress.Transactions {
			transaction := newTransactionFromExplorer(explorerTransaction)
			transactions.Transactions = append(transactions.Transactions, transaction)
		}
	}

utl:
	for _, unconfirmedTransaction := range unconfirmedTransactions.Transactions {
		for _, output := range unconfirmedTransaction.ScpOutputs {
			for _, address := range params.Addresses {
				if output.UnlockHash == address {
					transaction := newTransactionFromUnconfirmed(unconfirmedTransaction)
					transactions.Transactions = append(transactions.Transactions, transaction)
					continue utl
				}
			}
		}
		for _, input := range unconfirmedTransaction.ScpInputs {
			for _, publicKey := range input.UnlockConditions.PublicKeys {
				for _, filterPublicKey := range params.PublicKeys {
					if publicKey.Key == filterPublicKey {
						transaction := newTransactionFromUnconfirmed(unconfirmedTransaction)
						transactions.Transactions = append(transactions.Transactions, transaction)
						continue utl
					}
				}
			}
		}
	}
	return transactions

}
