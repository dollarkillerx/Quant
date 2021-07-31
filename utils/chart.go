package utils

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"os"
	"strconv"
)

func GenChatBarHtml(data []float64, title string, subtitle string, filename string) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    title,
		Subtitle: subtitle,
	}))

	axis := make([]string, 0)
	barData := make([]opts.BarData, 0)
	for _, v := range data {
		axis = append(axis, strconv.FormatFloat(v, 'E', -1, 32))
		barData = append(barData, opts.BarData{Value: v})
	}

	bar.SetXAxis(axis).AddSeries("K", barData)

	f, _ := os.Create(filename)
	bar.Render(f)
}
