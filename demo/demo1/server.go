package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

type Tick struct {
	Key string `json:"key"`
}

func main() {
	app := gin.Default()

	file, err := ioutil.ReadFile("./history/resse.json")
	if err != nil {
		log.Fatalln(err)
	}

	var his History
	err = json.Unmarshal(file, &his)
	if err != nil {
		log.Fatalln(err)
	}

	for i, v := range his.Data.Items {
		if v.OpenPrice-v.ClosePrice <= 0 {
			if v.Profit > 0 {
				his.Data.Items[i].Buy = false
			} else {
				his.Data.Items[i].Buy = true
			}
		} else {
			if v.Profit > 0 {
				his.Data.Items[i].Buy = true
			} else {
				his.Data.Items[i].Buy = false
			}
		}
	}

	app.POST("/tick", func(ctx *gin.Context) {
		var t Tick
		err := ctx.BindJSON(&t)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(t)

		ctx.JSONP(200, gin.H{"msg": "world"})
	})

	if err := app.Run("0.0.0.0:8085"); err != nil {
		log.Fatalln(err)
	}
}

type History struct {
	Code int         `json:"code"`
	Data HistoryData `json:"data"`
}

type HistoryData struct {
	Total       int            `json:"total"`
	TotalLots   int            `json:"total_lots"`
	TotalProfit int            `json:"total_profit"`
	Items       []HistoryItems `json:"items"`
}

type HistoryItems struct {
	Buy        bool    `json:"buy"`
	ID         int     `json:"id"`
	TradeID    int     `json:"trade_id"`
	Symbol     string  `json:"symbol"`
	Digits     int     `json:"digits"`
	Volume     float64 `json:"volume"`
	OpenTime   int     `json:"open_time"`
	OpenPrice  float64 `json:"open_price"`
	CloseTime  int     `json:"close_time"`
	ClosePrice float64 `json:"close_price"`
	TP         float64 `json:"tp"`
	Profit     float64 `json:"profit"`
	Ex         struct {
		StandardSymbol string  `json:"standard_symbol"`
		StandardLots   float64 `json:"standard_lots"`
	} `json:"ex"`
}
