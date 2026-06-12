use std::f64;

use chrono::{DateTime, Utc};
pub struct TransactionCost{
   pub commissions:f64,
    pub number_of_shares:i64,
    pub exchange_fee:f64,
    pub taxes_fee:f64,
    pub clearing_and_settlement_fees:f64,
    pub slippage:f64,
}
pub struct Tick{
    pub volume:f64,
    pub price:f64,
}

pub struct VolumeWeightedAveragePriceTracker{
    pub cumulative_volume:f64,
    pub cumulative_pv:f64,
}

impl VolumeWeightedAveragePriceTracker {
    pub fn new()-> Self{
        return Self { cumulative_volume: 0.0, cumulative_pv: 0.0 };
    }
    pub fn update(&mut self,price:f64,volume:f64)-> f64{
        self.cumulative_volume+=volume;
        self.cumulative_pv+=price*volume;
        return self.cumulative_pv/self.cumulative_volume;
    } 
}

pub struct TimeWeightedAveragePrice{
    pub share:f64,
    pub price:f64,
}

pub enum Side{
    Buy,
    Sell,
}
pub enum OrderType {
    Market,
    Limit{price:f64},
    Stop{trigger_price:f64},
}
pub struct Order{
    pub tick:Tick,
    pub side:Side,
    pub time_stamp:DateTime<Utc>,
    pub order_type:OrderType,
}

pub struct OrderBook{
    pub bids:Vec<(f64,f64)>,
    pub asks:Vec<(f64,f64)>,
}

impl OrderBook {
    pub fn mid_price(&self)-> Option<f64>{
        let best_bid=self.bids.first()?.0;
        let best_ask=self.asks.first()?.0;
        return Some((best_bid+best_ask)/2.0);
    }

    pub fn delay_cost(&self,shares_executed:f64,price_arrival:f64,price_decision:f64)->f64{
        return shares_executed*(price_arrival-price_decision);
    }

    pub fn tradin_cost(&self,shares_executed:f64,price_arrival:f64)->Option<f64>{
        
        match self.mid_price() {
            None=>{
                return None;
            },
            Some(ok)=>{
                return Some(shares_executed*(ok-price_arrival));
            }
        }

    }

    pub fn opportunity_cost(&self,share_desired:f64,share_executed:f64,price_close:f64,price_decision:f64)->f64{
        return (share_desired-share_executed)*(price_close-price_decision);
    }

    pub fn explicit_fees(t:TransactionCost)->f64{
        return t.clearing_and_settlement_fees+t.exchange_fee+t.taxes_fee+t.commissions;
    }


}

pub fn percentage_of_volume(order_volume: f64, market_volume: f64)->f64 {
    return order_volume / (order_volume + market_volume);
}

pub fn public_market_volume(participation_rate: f64, market_volume: f64)->f64 {
    return participation_rate * market_volume / (1.0 - participation_rate);
}

pub fn volume_weighted_average_price_historical(t:&[Tick]) -> Option<f64> {
   let mut sum_volume=0.0;
   let mut sum: f64 = 0.0;
   for tick in t{
        sum_volume+=tick.volume;
        sum+=tick.price*tick.volume;
   }    
   if sum_volume==0.0{
    return None;
   }
   return Some(sum/sum_volume);
}

pub fn bid_ask_spread(ask_price:f64,bid_price:f64)-> f64{
    return ask_price-bid_price;
}

pub fn implementation_shortfall(explicit_cost:f64,execution_cost:f64,slippage:f64,spread:f64,opportunity_cost:f64)->f64{
    return explicit_cost+(execution_cost*(slippage/spread))+opportunity_cost;
}

#[cfg(test)]
mod test{
    use super::*;
}
