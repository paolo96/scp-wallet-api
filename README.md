# SCP Wallet API
Backend for SCP Wallet app.

It exposes a simplified REST API which can be safely made public, since it shouldn't provide data related to the specific spd instance below.

Android App repo: https://github.com/paolo96/scp-wallet-android

## Requirements
*scpwalletapi* communicates with the [spd API](https://gitlab.com/scpcorp/ScPrime/-/blob/master/doc/API.md) to process data in a format useful to the SCP Wallet app.

To do so the following requirements must be met in order to run it:
* spd is listening on port 4280 or custom at localhost
* spd consensus module is synced
* spd explorer module is loaded
* spd transaction pool module is loaded
* TEMPORARY: spd.patch has been applied

## TEMPORARY Patch
*scpwalletapi* needs spd API to expose the endpoint /explorer/addresses/batch which is not included in the current version of spd.

A pull request will be made soon, but for now you can apply [spd.patch](spd.patch) to [ScPrime](https://gitlab.com/scpcorp/ScPrime) and build it from source.

## Get started
Build the package
```
go build -o scpwalletapi ./api
```

Run it
```
./scpwalletapi [coinmarketcap api key] [spd api port (default 4280)] [spd api password (default empty)] [custom port (default 14280)]
```

## Coinmarketcap API
To provide a SCP to Fiat exchange rate to the clients, it is recommended although not mandatory to provide a [Coinmarketcap API](https://coinmarketcap.com/api/documentation/v1/) key.

Requests are made every 5 minutes and cached. The API key should be used only for *scpwalletapi*, so that the free Basic plan will be enough, and it will not be Rate limited.

