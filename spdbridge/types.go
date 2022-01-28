package spdbridge

type (
	AddressesBatchParams struct {
		Addresses []string `json:"addresses"`
	}

	TransactionFeesResp struct {
		MinFee string `json:"minimum"`
		MaxFee string `json:"maximum"`
	}

	ConsensusResp struct {
		Synced bool   `json:"synced"`
		Height uint64 `json:"height"`
	}

	AddressesBatchResp struct {
		Addresses []ExplorerAddress `json:"addresses"`
	}

	TransactionPoolResp struct {
		Transactions []RawTransaction `json:"transactions"`
	}
)

type (
	ExplorerAddress struct {
		Address      string                `json:"address"`
		Transactions []ExplorerTransaction `json:"transactions"`
	}

	ExplorerTransaction struct {
		RawTransaction RawTransaction `json:"rawtransaction"`
		BlockTimestamp uint64         `json:"blocktimestamp"`
		Id             string         `json:"id"`
		Height         uint64         `json:"height"`
	}

	RawTransaction struct {
		ScpInputs  []ScpInput  `json:"siacoininputs"`
		ScpOutputs []ScpOutput `json:"siacoinoutputs"`
		MinerFees  []string    `json:"minerfees"`
	}

	TransactionOutput struct {
		Id             string `json:"id"`
		RelatedAddress string `json:"relatedaddress"`
		Value          string `json:"value"`
	}

	ScpOutput struct {
		Value      string `json:"value"`
		UnlockHash string `json:"unlockhash"`
		Id         string `json:"id"`
	}

	ScpInput struct {
		ParentId         string           `json:"parentid"`
		UnlockConditions UnlockConditions `json:"unlockconditions"`
	}

	UnlockConditions struct {
		Timelock           uint64         `json:"timelock"`
		SignaturesRequired uint64         `json:"signaturesrequired"`
		PublicKeys         []ScpPublicKey `json:"publickeys"`
	}

	ScpPublicKey struct {
		Key string `json:"key"`
	}
)
