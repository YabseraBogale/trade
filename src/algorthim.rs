use std::f64;

use chrono::{DateTime, Utc};
pub struct TransactionCost{
    commissions:f64,
    number_of_shares:i64,
    exchange_fee:f64,
    taxes_fee:f64,
    clearing_and_settlement_fees:f64,
    slippage:f64,
}
pub struct Tick{
    volume:f64,
    price:f64,
}

pub struct TimeWeightedAveragePrice{
    share:f64,
    price:f64,
}

pub struct Order{
    tick:Tick,
    side:String,
    time_stamp:DateTime<Utc>,
    order_type:String,
}

pub fn percentage_of_volume(order_volume: f64, market_volume: f64)->f64 {
    return order_volume / (order_volume + market_volume);
}

pub fn public_market_volume(participation_rate: f64, market_volume: f64)->f64 {
    return participation_rate * market_volume / (1.0 - participation_rate);
}

pub fn volume_weighted_averagePrice(t:&[Tick]) -> f64 {
   let mut sum_volume=0.0;
   let mut sum: f64 = 0.0;
   for tick in t{
        sum_volume+=tick.volume;
        sum+=tick.price*tick.volume;
   }    
   return sum/sum_volume;
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

