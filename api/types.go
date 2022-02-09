package main

import (
	"scp-app-api/spdbridge"
)

type (
	NetworkDataResponse struct {
		ScpPrice         *float64            `json:"scpPrice"`
		NetworkData      *NetworkData        `json:"networkData"`
		USDExchangeRates *map[string]float64 `json:"usdExchangeRates"`
	}

	NewTransactionParams struct {
		BroadcastData BroadcastData `json:"broadcastData"`
		ValidateData  string        `json:"validateData"`
	}

	TransactionsBatchParams struct {
		Addresses  []string `json:"addresses"`
		PublicKeys []string `json:"publickeys"`
	}

	TransactionsBatchResp struct {
		Transactions []Transaction `json:"transactions"`
	}
)

type (
	NetworkData struct {
		ConsensusHeight uint64 `json:"consensusHeight"`
		MinFee          string `json:"minFee"`
		MaxFee          string `json:"maxFee"`
	}

	BroadcastData struct {
		Parents     string `json:"parents"`
		Transaction string `json:"transaction"`
	}

	Transaction struct {
		ScpInputs      []spdbridge.ScpInput  `json:"siacoininputs"`
		ScpOutputs     []spdbridge.ScpOutput `json:"siacoinoutputs"`
		MinerFees      []string              `json:"minerfees"`
		Height         uint64                `json:"height"`
		BlockTimestamp uint64                `json:"blocktimestamp"`
		Id             string                `json:"id"`
	}
)

func newTransactionFromExplorer(eT spdbridge.ExplorerTransaction) (t Transaction) {
	t.ScpInputs = eT.RawTransaction.ScpInputs
	t.ScpOutputs = eT.RawTransaction.ScpOutputs
	t.MinerFees = eT.RawTransaction.MinerFees
	t.Id = eT.Id
	t.BlockTimestamp = eT.BlockTimestamp
	t.Height = eT.Height
	return t
}

func newTransactionFromUnconfirmed(rT spdbridge.RawTransaction) (t Transaction) {
	t.ScpInputs = rT.ScpInputs
	t.ScpOutputs = rT.ScpOutputs
	t.MinerFees = rT.MinerFees
	return t
}
