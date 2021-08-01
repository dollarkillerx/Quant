package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// macro 宏观 分析 涨跌 区间， 马丁策略风险预警模型
// 1. h1 常在区间, 2. h1 最大涨跌区间

// 技术分析 要点
// 1. 判断当前 是跌势 还是 涨势
// 2. > 30% 回调为一个新的阶段

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	open, err := os.Open("./data_utf8/eurusd_20180101_20210728_h1.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer open.Close()

	reader := csv.NewReader(open)

	var up bool                    // 存储当前趋势是否是向上的
	var max float64                // 当前趋势状态下的 最大值
	var min float64                // 当前趋势状态下的 最小值
	lowList := make([]float64, 0)  // 下跌趋势存储
	highList := make([]float64, 0) // 上行趋势存储

	// 初始化
	nx := GetNX(reader)
	up, err = nx.Up()
	if err != nil {
		log.Fatalln(err)
	}

	if up {
		highList = append(highList, nx.Before, nx.After)
		max, min = nx.After, nx.Before
	} else {
		lowList = append(lowList, nx.Before, nx.After)
		max, min = nx.Before, nx.After
	}

	// 极值 统计
	extremumUp := make([]float64, 0)
	extremumLow := make([]float64, 0)

	// 以ASK 来确定
	for {
		now, err := GetKChart(reader)
		if err != nil {
			break
		}

		// 上行阶段
		if up {
			// 极点判断
			last := GetLast(highList)
			if last <= now.Ask {
				if max <= now.Ask { // 判断是否有新的突破
					max = now.Ask
				}
				highList = append(highList, now.Ask)
			} else {
				// 判断是否出现了新的趋势
				point := TurningPointUp(max, min)
				if now.Ask <= point {
					// 记录旧的趋势
					if max-min != 0 {
						extremumUp = append(extremumUp, max-min)
					}

					//新的趋势出现了
					up = false
					highList = highList[0:0] // 清空 上线数组 回收垃圾
					lowList = append(lowList, now.Ask)
					min, max = now.Ask, now.Ask
				} else {
					// 反之是回调
					highList = append(highList, now.Ask)
					//fmt.Println(now.Ask, "  ", point)
				}
			}

		} else {
			// 下行阶段
			last := GetLast(lowList)
			if last >= now.Ask {
				if min >= now.Ask {
					min = now.Ask
				}
				lowList = append(lowList, now.Ask)
			} else {
				// 判断是否出现了新的趋势
				point := TurningPointLow(max, min)
				if now.Ask >= point {
					if max-min != 0 {
						extremumLow = append(extremumLow, max-min)
					}

					//新的趋势出现了
					up = true
					lowList = lowList[0:0] // 清空 上线数组 回收垃圾
					highList = append(highList, now.Ask)
					min, max = now.Ask, now.Ask
				} else {
					// 反之是回调
					lowList = append(lowList, now.Ask)
				}
			}
		}
	}

	//fmt.Println("low: ", len(lowList))
	//fmt.Println("hig: ", len(highList))
	sort.Float64s(extremumUp)
	sort.Float64s(extremumLow)
	//fmt.Println(extremumUp)
	fmt.Printf("Ext UP: totol: %d  Min: %f Max: %f \n", len(extremumUp), extremumUp[0], extremumUp[len(extremumUp)-1])
	fmt.Printf("Ext Low: totol: %d Min: %f Max: %f \n", len(extremumLow), extremumLow[0], extremumLow[len(extremumLow)-1])

	//fmt.Printf("Ext UP: totol: %d  \n", len(extremumUp))
	//fmt.Printf("Ext Low: totol: %d \n", len(extremumLow))

	//utils.GenChatBarHtml(extremumUp, "EURUSD 1H 上升趋势统计", "EURUSD 1H 上升趋势统计", "eurusd_up_h1.html")
	//utils.GenChatBarHtml(extremumLow, "EURUSD 1H 下降趋势统计", "EURUSD 1H 下降趋势统计", "eurusd_low_h1.html")

	// GO数据统计的包太少了 没有泛型太TM搞人了

	// 平均数
	// 中位数
	// 众数
	tol := NumTol(0.008, extremumUp)
	fmt.Printf("extremumUp 平均数: %f 中位数: %f 众数占比: %f  \n", tol.Average, tol.Median, tol.Mode*100)
	tol = NumTol(0.008, extremumLow)
	fmt.Printf("extremumLow 平均数: %f 中位数: %f 众数占比: %f  \n", tol.Average, tol.Median, tol.Mode*100)

	marshal, err := json.Marshal(extremumLow)
	if err != nil {
		return
	}
	fmt.Printf(string(marshal))
	// 按照80点风险计算 爆仓率 4%
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

type Nx struct {
	Before float64 // 这根K
	After  float64 // 下根K
}

func GetNX(r *csv.Reader) Nx {
	resp := Nx{}

	before, err := GetKChart(r)
	if err != nil {
		return resp
	}
	resp.Before = before.Ask
	after, err := GetKChart(r)
	if err != nil {
		return resp
	}
	resp.After = after.Ask

	return resp
}

func (n *Nx) Up() (bool, error) {
	if n.Before == 0 || n.After == 0 {
		return false, io.EOF
	}

	if n.Before > n.After {
		return false, nil
	}

	return true, nil
}

func (n *Nx) Over() bool {
	if n.Before == 0 || n.After == 0 {
		return true
	}
	return false
}

// TurningPointUp 获取转折点 30%
func TurningPointUp(max, min float64) float64 {
	return max - (max-min)*0.3
}

// TurningPointLow 获取转折点 30%
func TurningPointLow(max, min float64) float64 {
	return (max-min)*0.3 + min
}

func GetLast(datas []float64) float64 {
	return datas[len(datas)-1]
}

type NumTolSt struct {
	Average float64
	Median  float64
	Mode    float64
}

func NumTol(r float64, list []float64) NumTolSt {
	return NumTolSt{
		Average: Average(list),
		Median:  Median(list),
		Mode:    Mode(r, list),
	}
}

func Average(list []float64) float64 {
	var r float64
	for _, v := range list {
		r += v
	}

	return r / float64(len(list))
}

func Median(list []float64) float64 {
	var r float64
	if len(list)%2 == 0 {
		r = list[len(list)/2]
	} else {
		r = (list[len(list)/2] + list[len(list)/2+1]) / 2
	}

	return r
}

func Mode(r float64, list []float64) float64 {
	var t int
	for _, v := range list {
		if v <= r {
			t += 1
		}
	}

	return float64(t) / float64(len(list))
}
