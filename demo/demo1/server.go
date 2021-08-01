package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Tick struct {
	Key string `json:"key"`
}

func main() {
	app := gin.Default()

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
