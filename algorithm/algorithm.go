package algorithm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Constituent struct {
	Name              string
	CurrentPrice      float64
	SharesOutStanding float64
}

type Portfolio struct {
	Cash           float64
	PositionShares float64
}

type BackTest struct {
	Portfolio Portfolio
	Fees      TransactionCost
}

type TransactionCost struct {
	Commissions               float64
	NumberOfShares            int
	ExchangeFee               float64
	TaxesFee                  float64
	ClearingAndSettlementFees float64
	Slippage                  float64
}

type TickerConstituent struct {
	Name             string  `json:"name"`
	ShareOutStanding float64 `json:"shares_outstanding"`
}

func FetchConstituent(apiURL string) ([]TickerConstituent, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var constituent []TickerConstituent
	err = json.NewDecoder(resp.Body).Decode(&constituent)
	if err != nil {
		return nil, err
	}
	return constituent, nil

}

func CalculateMarketCapIndex(currentAssts []Constituent, baseMarketCap float64, baseIndexValue float64) (float64, error) {
	if baseMarketCap <= 0 {
		return 0, fmt.Errorf("base market cap must be greater than zero")
	}
	var currentTotalMarketCap float64

	for _, assets := range currentAssts {
		currentTotalMarketCap += assets.CurrentPrice * assets.SharesOutStanding
	}

	indexValue := (currentTotalMarketCap / baseMarketCap) * baseIndexValue
	return indexValue, nil

}

type InitalPrice struct {
	Name      string
	BasePrice float64
}

func CalculateEqualWeightIndex(currentAssets []Constituent, basePrices map[string]float64) (float64, error) {
	if len(currentAssets) <= 0 {
		return 0, fmt.Errorf("no asset provided")
	}

	var totalReturn float64
	var countedAssets float64

	for _, asset := range currentAssets {
		basePrice, exists := basePrices[asset.Name]
		if exists && basePrice > 0 {
			totalReturn += asset.CurrentPrice / basePrice
			countedAssets++
		}
	}
	if countedAssets == 0 {
		return 0, fmt.Errorf("no matching base price found")
	}

	return (totalReturn / countedAssets), nil
}

func NewBackTestEngine(initalCash float64, fee TransactionCost) *BackTest {
	return &BackTest{
		Portfolio: Portfolio{
			Cash:           initalCash,
			PositionShares: 0,
		},
		Fees: fee,
	}
}

func (b *BackTest) Run(data []PriceLevel, targetShares float64, side string) (*Order, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no historical data provided for backtest")
	}
	order := &Order{
		Side:           strings.ToUpper(side),
		SharesDesired:  targetShares,
		SharesExecuted: 0,
		PriceDecision:  data[0].Price,
		PriceArrival:   data[0].Price,
		PriceClose:     data[len(data)-1].Price,
	}
	vwapTracker := New()
	var totalExecutionCost float64
	fmt.Printf("Starting Backtest: %s %.2f shares. Initial Cash: $%.2f\n", order.Side, targetShares, b.Portfolio.Cash)
	for i, tick := range data {
		currentVwap := vwapTracker.Update(tick.Price, tick.Volume)
		if order.SharesExecuted > order.SharesDesired {
			participationRate := 0.10
			sliceShare := tick.Volume * participationRate

			if order.SharesExecuted+sliceShare > order.SharesDesired {
				sliceShare = order.SharesDesired - order.SharesExecuted
			}
			if sliceShare > 0 {
				if order.Side == "BUY" {
					cost := sliceShare * tick.Price
					if b.Portfolio.Cash >= cost {
						b.Portfolio.Cash -= cost
						b.Portfolio.PositionShares += sliceShare
						order.SharesExecuted += sliceShare
						totalExecutionCost += sliceShare * tick.Price
					}
				} else if order.Side == "SELL" {
					if b.Portfolio.PositionShares >= sliceShare {
						b.Portfolio.PositionShares -= sliceShare
						b.Portfolio.Cash += sliceShare * tick.Price
						order.SharesExecuted += sliceShare
						totalExecutionCost += sliceShare * tick.Price

					}
				}
				fmt.Printf("[Tick %d] Market Price: %.2f | Running VWAP: %.2f | Executed Slice: %.2f\n",
					i, tick.Price, currentVwap, sliceShare)
			}
		}
	}
	if order.SharesExecuted > 0 {
		order.PriceExecuted = totalExecutionCost / order.SharesExecuted
	}
	return order, nil
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

type Symbole struct {
	Name              string  `json:"Name"`
	SharesOutstanding float64 `json:"SharesOutstanding"`
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

func FetchVolumeAndClosePriceFromURL(apiURL string) ([]PriceLevel, error) {
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

func FetchNameList(apiURL string) ([]Symbole, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	symbole := []Symbole{}

	err = json.NewDecoder(resp.Body).Decode(&symbole)
	if err != nil {
		return nil, err
	}

	return symbole, nil

}
