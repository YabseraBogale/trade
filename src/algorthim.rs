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

pub struct TimeWeightedAveragePrice{
    pub share:f64,
    pub price:f64,
}

pub struct Order{
    pub tick:Tick,
    pub side:String,
    pub time_stamp:DateTime<Utc>,
    pub order_type:String,
}

pub fn percentage_of_volume(order_volume: f64, market_volume: f64)->f64 {
    return order_volume / (order_volume + market_volume);
}

pub fn public_market_volume(participation_rate: f64, market_volume: f64)->f64 {
    return participation_rate * market_volume / (1.0 - participation_rate);
}

pub fn volume_weighted_average_price(t:&[Tick]) -> Option<f64> {
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

