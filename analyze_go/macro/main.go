package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// macro 宏观 分析 涨跌 区间， 马丁策略风险预警模型
// 1. h1 常在区间, 2. h1 最大涨跌区间

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	open, err := os.Open("./data_utf8/eurusd_20180101_20210728_h1.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer open.Close()

	reader := csv.NewReader(open)

	for {
		chart, err := GetKChart(reader)
		if err != nil {
			break
		}

		fmt.Println(chart)
	}
}

type CandlestickChart struct {
	Ask     float64 `json:"ask"`
	AskHigh float64 `json:"ask_high"`
	AskLow  float64 `json:"ask_low"`
	Bid     float64 `json:"bid"`
	BidHigh float64 `json:"bid_high"`
	BidLow  float64 `json:"bid_low"`
	Time    string  `json:"time"`
}

func GetKChart(r *csv.Reader) (CandlestickChart, error) {
	read, err := r.Read()
	if err != nil {
		return CandlestickChart{}, err
	}

	split := strings.Split(read[0], "\t")

	if len(split) != 7 {
		return CandlestickChart{}, io.EOF
	}

	ask, err := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
	if err != nil {
		fmt.Println(split[0])
		log.Println(err)
		return CandlestickChart{}, err
	}

	askHigh, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		log.Println(err)
		return CandlestickChart{}, err
	}

	askLow, err := strconv.ParseFloat(split[2], 64)
	if err != nil {
		log.Println(err)
		return CandlestickChart{}, err
	}

	bid, err := strconv.ParseFloat(split[3], 64)
	if err != nil {
		log.Println(err)
		return CandlestickChart{}, err
	}

	bidHigh, err := strconv.ParseFloat(split[4], 64)
	if err != nil {
		log.Println(err)
		return CandlestickChart{}, err
	}

	bidLow, err := strconv.ParseFloat(split[5], 64)
	if err != nil {
		log.Println(err)
		return CandlestickChart{}, err
	}

	kChart := CandlestickChart{
		Ask:     ask,
		AskHigh: askHigh,
		AskLow:  askLow,
		Bid:     bid,
		BidHigh: bidHigh,
		BidLow:  bidLow,
		Time:    split[6],
	}

	return kChart, nil
}
