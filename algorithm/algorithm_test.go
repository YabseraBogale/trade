package algorithm

import (
	"fmt"
	"testing"
	"trade/algorithm"
)

func TestFetch(t *testing.T) {
	url := "http://127.0.0.1:5000/"
	res, err := algorithm.FetchVolumeAndClosePriceFromURL(url)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	fmt.Println("Successfully processed: length of data ", len(res))
}
