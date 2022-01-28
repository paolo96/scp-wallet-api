package main

import (
	"github.com/julienschmidt/httprouter"
)

func buildRouter() *httprouter.Router {

	version := "/:version"

	router := httprouter.New()
	router.GET(version+"/scprime/data", getScPrimeDataHandler)
	router.POST(version+"/addresses/transactions/batch", getAddressesTransactionsBatchHandler)
	router.POST(version+"/transactions", newTransactionHandler)

	return router

}
