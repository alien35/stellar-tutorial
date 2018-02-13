package main

import (
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	"log"
	"net/http"
	"github.com/stellar/go/clients/horizon"
)

func fillAccounts(addresses []string) {
	for _, address := range addresses {
		friendBotResp, err := http.Get("https://horizon-testnet.stellar.org/friendbot?addr=" + address)
		if err != nil {
			log.Fatal(err)
		}
		defer friendBotResp.Body.Close()
	}
}
func logBalances(addresses []string) {
	for _, address := range addresses {
		account, err := horizon.DefaultTestNetClient.LoadAccount(address)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Balances for address:", address)
		for _, balance := range account.Balances {
			log.Println(balance)
		}
	}
}

func trust(asset build.Asset, limit string, address string, seed string) {
   tx, err := build.Transaction(
      build.SourceAccount{address},
      build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
      build.TestNetwork,
      build.Trust(asset.Code, asset.Issuer, build.Limit(limit)),
   )
    if err != nil { log.Fatal(err) }
    signAndSubmit(tx, seed, "Create orange trustline")
}

func signAndSubmit(tx *build.TransactionBuilder, seed string, privateMemo string) {
    txe, err := tx.Sign(seed)
   if err != nil { log.Fatal(err) }
   txeB64, err := txe.Base64()
   if err != nil { log.Fatal(err) }
   _, err = horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
   if err != nil {
      log.Fatal(err)
   }
   log.Println("Transaction successfully submitted:", privateMemo)
}

func makeOrangePurchaseOffer(rate build.Rate, amount build.Amount, address string, seed string) {
    tx, err := build.Transaction(
      build.SourceAccount{address},
      build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
      build.TestNetwork,
      build.CreateOffer(rate, amount),
   )
    if err != nil { log.Fatal(err) }
    signAndSubmit(tx, seed, "Make purchase offer")
}

func makeOrangeSellOffer(rate build.Rate, amount build.Amount, address string, seed string) {
    tx, err := build.Transaction(
      build.SourceAccount{address},
      build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
      build.TestNetwork,
      build.CreateOffer(rate, amount),
   )
    if err != nil { log.Fatal(err) }
    signAndSubmit(tx, seed, "Make sell offer")
}

func confirmPayment(asset build.Asset, amount string, destinationAddress string, seed string) {
    tx, err := build.Transaction(
    		build.SourceAccount{seed},
    		build.TestNetwork,
    		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
    		build.Payment(
   			    build.Destination{AddressOrSeed: destinationAddress},
   			    build.CreditAmount{asset.Code, asset.Issuer, amount},
   		    ),
    )
    if err != nil { log.Fatal(err) }
    signAndSubmit(tx, seed, "Oranges arrived")
}

func main() {
	seller, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	buyer, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	addresses := []string{seller.Address(), buyer.Address()}
	fillAccounts(addresses)
	logBalances(addresses)
	OrangeCreditAsset := build.CreditAsset("Orange", seller.Address())
    trust(OrangeCreditAsset, "500", buyer.Address(), buyer.Seed())
    sellOrangeRate := build.Rate{Selling: OrangeCreditAsset, Buying: build.NativeAsset(), Price: "0.5"}
    buyOrangeRate := build.Rate{Selling: build.NativeAsset(), Buying: OrangeCreditAsset, Price: "2"}
    makeOrangeSellOffer(sellOrangeRate, "20", seller.Address(), seller.Seed())
    makeOrangePurchaseOffer(buyOrangeRate, "20", buyer.Address(), buyer.Seed())
    logBalances(addresses)
    confirmPayment(OrangeCreditAsset, "20", seller.Address(), buyer.Seed())
    logBalances(addresses)
}
