package algorithm_test

import (
	"fmt"
	"math/rand"
	"testing"
	"trade/algorithm"
)

func TestBackTest(t *testing.T) {

}

func TestFetchVolumeAndClosePriceFromURL(t *testing.T) {

	url := "http://127.0.0.1:5000/name"
	resName, err := algorithm.FetchNameList(url)
	if err != nil {
		t.Fatalf("Failed Error: %v", err)
	}
	fmt.Println("Successfully processed Symbole length of ", len(resName))
	index := rand.Intn(len(resName))
	url = "http://127.0.0.1:5000/" + resName[index].Name
	resVolumeClosePrice, err := algorithm.FetchVolumeAndClosePriceFromURL(url)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	fmt.Println("Successfully processed Volume and Closed Price ", len(resVolumeClosePrice))
}
