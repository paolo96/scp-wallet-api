package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"scp-app-api/spdbridge"
)

const verbose = false

var port = "14280"

func main() {

	if !checkSpd() {
		log.Fatal("spd daemon connection failed, check that:\n" +
			"- spd API is running at " + spdbridge.SpdApiURL + ":" + spdbridge.SpdApiPort + "\n" +
			"- spd consensus module is synced\n" +
			"- spd explorer module is loaded\n" +
			"- spd transaction pool module is loaded\n" +
			"- spd.patch has been applied\n" +
			"Command example: ./scpwalletapi [coinmarketcap api key] [spd api port] [spd api password] [custom port]")
	}

	StartDataSync()

	fmt.Println("Starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, buildRouter()))

}

//Checks if we can connect to spd
func checkSpd() bool {

	if len(os.Args) > 1 {
		CMCApiKey = os.Args[1]
	} else if CMCApiKey == "" {
		fmt.Println("No coinmarketcap API KEY found, fiat quotes will not be available to clients.")
	}
	if len(os.Args) > 2 {
		spdbridge.SpdApiPort = os.Args[2]
	}
	if len(os.Args) > 3 {
		spdbridge.SpdApiPassword = os.Args[3]
	}
	if len(os.Args) > 4 {
		port = os.Args[4]
	}

	consensus, err := spdbridge.GetConsensus()
	if err != nil {
		fmt.Printf("Test call to consensus failed with error: %v\n\n", err)
		return false
	} else if !consensus.Synced {
		fmt.Printf("Consensus is not synced yet\n\n")
		return false
	}
	_, err = spdbridge.GetTransactionPoolFees()
	if err != nil {
		fmt.Printf("Test call to transaction pool failed with error: %v\n\n", err)
		return false
	}
	_, err = spdbridge.ExplorerAddressesBatch([]string{})
	if err != nil {
		fmt.Printf("Test call to explorer failed with error: %v\n\n", err)
		return false
	}

	return true

}
