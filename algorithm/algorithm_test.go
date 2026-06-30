package algorithm_test

import (
	"fmt"
	"math/rand"
	"testing"
	"trade/algorithm"
)

func TestBackTest(t *testing.T) {
	url := "http://localhost:8080/name"
	resName, err := algorithm.FetchNameList(url)
	if err != nil {
		t.Fatalf("Failed Error at resName %v", err)
	}
	index := rand.Intn(len(resName))
	tickerURL := "http://localhost:8080/" + resName[index].Name
	histroicalData, err := algorithm.FetchVolumeAndClosePriceFromURL(tickerURL)
	if err != nil {
		t.Fatalf("Failed Error at histroicalData %v", err)
	}
	fee := algorithm.TransactionCost{
		Commissions:               0.0,
		ExchangeFee:               0.0,
		TaxesFee:                  0.0,
		ClearingAndSettlementFees: 0.0,
	}
	engine := algorithm.NewBackTestEngine(1000, fee)
	orderResults, err := engine.Run(histroicalData, 500.0, "BUY")
	if err != nil {
		t.Fatalf("Back Test error at orderResults %v", err)
	}
	shortFall, err := orderResults.ImplementationShortfall(fee)
	if err != nil {
		t.Fatalf("Short fall error at shortFall %v", err)
	}
	fmt.Println("Total Slippage & Fees (Implementation Shortfall):", shortFall)
	fmt.Printf("\n--- BACKTEST RESULTS FOR %s ---\n", resName[index].Name)
	fmt.Printf("Shares Desired:  %.2f | Shares Executed: %.2f\n", orderResults.SharesDesired, orderResults.SharesExecuted)
	fmt.Printf("Decision Price:  $%.2f  | Avg Executed Price: $%.2f\n", orderResults.PriceDecision, orderResults.PriceExecuted)
}

func TestFetchVolumeAndClosePriceFromURL(t *testing.T) {

	url := "http://localhost:8080/name"
	resName, err := algorithm.FetchNameList(url)
	if err != nil {
		t.Fatalf("Failed Error at resName %v", err)
	}
	fmt.Println("Successfully processed Symbole length of ", len(resName))
	index := rand.Intn(len(resName))
	url = "http://127.0.0.1:5000/" + resName[index].Name
	resVolumeClosePrice, err := algorithm.FetchVolumeAndClosePriceFromURL(url)
	if err != nil {
		t.Fatalf("Failed Error at resVolumeClosePrice %v", err)
	}
	fmt.Println("Successfully processed Volume and Closed Price ", len(resVolumeClosePrice))
}
