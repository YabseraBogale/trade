package algorithm

import (
	"crypto/rand"
	"fmt"
	"testing"
	"trade/algorithm"
)

func TestFetchVolumeAndClosePriceFromURL(t *testing.T) {

	url := "http://127.0.0.1:5000/name"
	resName, err := algorithm.FetchNameList(url)
	if err != nil {
		t.Fatalf("Failed Error: %v", err)
	}

	index, err := rand.Int(len(resName))
	if err != nil {
		t.Fatalf("Failed Error: %v", err)
	}
	url = "http://127.0.0.1:5000/" + resName[index]
	resVolumeClosePrice, err := algorithm.FetchVolumeAndClosePriceFromURL(url)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	fmt.Println("Successfully processed: length of data ", len(resVolumeClosePrice))
}
