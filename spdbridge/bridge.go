package spdbridge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var SpdApiPort = "4280"
var SpdApiPassword = ""

const SpdApiURL = "http://127.0.0.1"

const headerJSON = "application/json"
const headerAgent = "ScPrime-Agent"

const verbose = false

//GetConsensus performs a GET request ScPrime API endpoint /consensus
func GetConsensus() (*ConsensusResp, error) {
	resp, e := getRequest("/consensus")
	if e != nil {
		return nil, e
	}

	var data ConsensusResp
	e = json.Unmarshal(resp, &data)
	if e != nil {
		return nil, e
	}

	return &data, nil
}

//GetTransactionPoolFees performs a GET request ScPrime API endpoint /tpool/fee
func GetTransactionPoolFees() (*TransactionFeesResp, error) {
	resp, e := getRequest("/tpool/fee")
	if e != nil {
		return nil, e
	}

	var data TransactionFeesResp
	e = json.Unmarshal(resp, &data)
	if e != nil {
		return nil, e
	}

	return &data, nil
}

//GetTransactionPool performs a GET request ScPrime API endpoint /tpool/fee
func GetTransactionPool() (*TransactionPoolResp, error) {
	resp, e := getRequest("/tpool/transactions")
	if e != nil {
		return nil, e
	}

	var data TransactionPoolResp
	e = json.Unmarshal(resp, &data)
	if e != nil {
		return nil, e
	}

	return &data, nil
}

//TransactionPoolRaw performs a POST request ScPrime API endpoint /tpool/raw
func TransactionPoolRaw(parents string, transaction string) (bool, error) {
	requestData := url.Values{}
	requestData.Set("parents", parents)
	requestData.Set("transaction", transaction)

	_, e := postRequestForm("/tpool/raw", requestData)
	if e != nil {
		return false, e
	}

	return true, nil
}

//ConsensusValidateTxns performs a POST request ScPrime API endpoint /consensus/validate/transactionset
func ConsensusValidateTxns(txnsData []byte) (bool, error) {
	_, e := postRequestJSON("/consensus/validate/transactionset", txnsData)
	if e != nil {
		return false, e
	}

	return true, nil
}

//ExplorerAddressesBatch performs a POST request ScPrime API endpoint /explorer/transactions/batch
func ExplorerAddressesBatch(addresses []string) (*AddressesBatchResp, error) {

	jsonRequest, e := json.Marshal(AddressesBatchParams{
		Addresses: addresses,
	})
	if e != nil {
		return nil, e
	}

	resp, e := postRequestJSON("/explorer/addresses/batch", jsonRequest)
	if e != nil {
		return nil, e
	}

	var data AddressesBatchResp
	e = json.Unmarshal(resp, &data)
	if e != nil {
		return nil, e
	}

	return &data, nil
}

//getRequest performs a GET request to ApiURL/path tailored to ScPrime API
func getRequest(path string) ([]byte, error) {

	req, e := http.NewRequest("GET", SpdApiURL+":"+SpdApiPort+path, nil)
	if e != nil {
		return nil, e
	}

	req.Header.Set("User-Agent", headerAgent)
	req.SetBasicAuth("", SpdApiPassword)

	client := &http.Client{}
	response, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer response.Body.Close()

	body, e := ioutil.ReadAll(response.Body)
	if verbose {
		fmt.Printf("GET %v -> %v\n", path, string(body))
	}
	if e != nil {
		return nil, e
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {

		return body, nil

	} else {
		return nil, errors.New(path + " error " + strconv.Itoa(response.StatusCode))
	}
}

//postRequestJSON performs a POST request with params in a JSON body to ApiURL/path tailored to ScPrime API
func postRequestJSON(path string, params []byte) ([]byte, error) {

	req, e := http.NewRequest("POST", SpdApiURL+":"+SpdApiPort+path, bytes.NewBuffer(params))
	if e != nil {
		return nil, e
	}

	req.Header.Set("User-Agent", headerAgent)
	req.Header.Set("Content-Type", headerJSON)
	req.SetBasicAuth("", SpdApiPassword)

	client := &http.Client{}
	response, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer response.Body.Close()

	body, e := ioutil.ReadAll(response.Body)
	if verbose {
		fmt.Printf("POST %v -> %v\n", path, string(body))
	}
	if e != nil {
		return nil, e
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return body, nil
	} else {
		return nil, errors.New(path + " error " + strconv.Itoa(response.StatusCode))
	}
}

//postRequestForm performs a POST request with form params to ApiURL/path tailored to ScPrime API
func postRequestForm(path string, params url.Values) ([]byte, error) {

	req, e := http.NewRequest("POST", SpdApiURL+":"+SpdApiPort+path, strings.NewReader(params.Encode()))
	if e != nil {
		return nil, e
	}

	req.Header.Set("User-Agent", headerAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("", SpdApiPassword)

	client := &http.Client{}
	response, e := client.Do(req)
	if e != nil {
		return nil, e
	}
	defer response.Body.Close()

	body, e := ioutil.ReadAll(response.Body)
	if verbose {
		fmt.Printf("POST %v -> %v\n", path, string(body))
	}
	if e != nil {
		return nil, e
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return body, nil
	} else {
		return nil, errors.New(path + " error " + strconv.Itoa(response.StatusCode))
	}
}
