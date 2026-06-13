package algorithm

import "time"

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
	Tick      Tick
	Side      string
	TimeStamp time.Time
	Market    string
	Limit     float64
}

type OrderBook struct {
	Bids []float64
	Asks []float64
}

func (o *OrderBook) MidPrice() float64 {
	best_bid := o.Bids[0]
	best_ask := o.Asks[0]
	return (best_bid + best_ask) / 2.0
}

func (o *OrderBook) DelayCost(shares_executed float64, price_arrival float64, price_decision float64) float64 {
	return shares_executed * (price_arrival - price_decision)
}

func (o *OrderBook) TradingCost(shares_executed float64, price_arrival float64) float64 {
	ok := o.MidPrice()
	return shares_executed * (ok - price_arrival)

}

func (o *OrderBook) OpportunityCost(share_desired float64, share_executed float64, price_close float64, price_decision float64) float64 {
	return (share_desired - share_executed) * (price_close - price_decision)
}

func (o *OrderBook) ExplicitFees(t TransactionCost) float64 {
	return t.ExchangeFee + t.ClearingAndSettlementFees + t.Commissions + t.TaxesFee
}

func PercentageOfVolume(order_volume float64, market_volume float64) float64 {
	return order_volume / (order_volume + market_volume)
}

func PublicMarketVolume(participation_rate float64, market_volume float64) float64 {
	return participation_rate * market_volume / (1.0 - participation_rate)
}

func VolumeWeightedAveragePriceHistorical(t []Tick) float64 {
	sum_volume := 0.0
	sum := 0.0
	for _, tick := range t {
		sum_volume += tick.Volume
		sum += tick.Price * tick.Volume
	}
	return sum / sum_volume
}

func BidAskSpread(ask_price float64, bid_price float64) float64 {
	return ask_price - bid_price
}

func ImplementationShortfAll(explicit_cost float64, execution_cost float64, slippage float64, spread float64, opportunity_cost float64) float64 {
	return explicit_cost + (execution_cost * (slippage / spread)) + opportunity_cost
}
