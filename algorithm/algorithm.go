package algorithm

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type TransactionCost struct {
	Commissions               float64
	NumberOfShares            int
	ExchangeFee               float64
	TaxesFee                  float64
	ClearingAndSettlementFees float64
	Slippage                  float64
}

type Tick struct {
	Volume float64
	Price  float64
}

type VolumeWeightedAveragePriceTracker struct {
	CumulativeVolume float64
	CumulativePv     float64
}

func New() VolumeWeightedAveragePriceTracker {
	return VolumeWeightedAveragePriceTracker{
		CumulativeVolume: 0.0,
		CumulativePv:     0.0,
	}
}

func (v *VolumeWeightedAveragePriceTracker) Update(price float64, volume float64) float64 {
	v.CumulativeVolume += volume
	v.CumulativePv += price * volume
	return v.CumulativePv / v.CumulativeVolume
}

type TimeWeightedAveragePrice struct {
	Share float64
	Price float64
}

type Order struct {
	Side           string
	SharesDesired  float64
	SharesExecuted float64
	PriceDecision  float64
	PriceArrival   float64
	PriceExecuted  float64
	PriceClose     float64
}

type PriceLevel struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

type OrderBook struct {
	Bids []PriceLevel
	Asks []PriceLevel
}

func (o *OrderBook) MidPrice() float64 {
	bestBid := o.Bids[0].Price
	bestAsk := o.Asks[0].Price
	return (bestBid + bestAsk) / 2.0
}

func (o *OrderBook) DelayCost(sharesExecuted float64, priceArrival float64, priceDecision float64) float64 {
	return sharesExecuted * (priceArrival - priceDecision)
}

func (o *OrderBook) TradingCost(sharesExecuted float64, priceArrival float64) float64 {
	ok := o.MidPrice()
	return sharesExecuted * (ok - priceArrival)

}

func (o *OrderBook) OpportunityCost(shareDesired float64, sharedExecuted float64, priceClose float64, priceDecision float64) float64 {
	return (shareDesired - sharedExecuted) * (priceClose - priceDecision)
}

func (t *TransactionCost) TotalExplicitFees() float64 {
	return t.ExchangeFee + t.ClearingAndSettlementFees + t.Commissions + t.TaxesFee
}

func (o *Order) ImplementationShortfall(fees TransactionCost) (float64, error) {
	side := strings.ToUpper(o.Side)
	if side != "BUY" && side != "SELL" {
		return 0.0, errors.New("invalid order side: must be BUY or SELL")
	}

	var executionCost float64
	if side == "BUY" {
		executionCost = o.SharesExecuted * (o.PriceExecuted - o.PriceDecision)
	} else {
		executionCost = o.SharesExecuted * (o.PriceDecision - o.PriceExecuted)
	}

	var opportunityCost float64
	sharesUnexecuted := o.SharesDesired - o.SharesExecuted
	if sharesUnexecuted > 0 {
		if side == "BUY" {
			opportunityCost = sharesUnexecuted * (o.PriceClose - o.PriceDecision)
		} else {
			opportunityCost = sharesUnexecuted * (o.PriceDecision - o.PriceClose)
		}
	}

	explicitCost := fees.TotalExplicitFees()

	totalIS := executionCost + opportunityCost + explicitCost
	return totalIS, nil
}

func PercentageOfVolume(orderVolume float64, marketVolume float64) float64 {
	return orderVolume / (orderVolume + marketVolume)
}

func PublicMarketVolume(participationRate float64, marketVolume float64) float64 {
	return participationRate * marketVolume / (1.0 - participationRate)
}

func VolumeWeightedAveragePriceHistorical(t []Tick) float64 {
	sumVolume := 0.0
	sum := 0.0
	for _, tick := range t {
		sumVolume += tick.Volume
		sum += tick.Price * tick.Volume
	}
	if sumVolume == 0.0 {
		return 0.0
	}
	return sum / sumVolume
}

func BidAskSpread(askPrice float64, bidPrice float64) float64 {
	return askPrice - bidPrice
}

func FetchAndParseFromURL(apiURL string) ([]PriceLevel, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	priceLevels := []PriceLevel{}
	err = json.NewDecoder(resp.Body).Decode(&priceLevels)
	if err != nil {
		return nil, err
	}

	return priceLevels, nil

}
