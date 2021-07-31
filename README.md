# Quant  

> 为啥感觉Forex市场非常闭塞 sharpe ratio 0.29都可以上台面?   所以有了此库 开源交易策略  大家一同学习 一同进步

### 提供语言支持
- mql5
- python
- golang

### data:

```
ask,ask_high,ask_low,bid,bid_high,bid_low,time

type CandlestickChart struct {
	Ask     float64 `json:"ask"`
	AskHigh float64 `json:"ask_high"`
	AskLow  float64 `json:"ask_low"`
	Bid     float64 `json:"bid"`
	BidHigh float64 `json:"bid_high"`
	BidLow  float64 `json:"bid_low"`
	Time    string  `json:"time"`
}
```