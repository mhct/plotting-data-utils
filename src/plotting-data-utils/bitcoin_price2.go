package main

import (
	"fmt"
	"io/ioutil"
	"encoding/csv"
	"strings"
	"io"
	"log"
	"strconv"
	"github.com/wcharczuk/go-chart"
	"os"
	"image/png"
	"time"
	"math"
)

type StatsBuffer struct {
	priceAccumulated float64
	accumulatedVolume float64
	nbAccumulatedPrices float64
}

func (r *StatsBuffer) accumulate(price float64, volume float64) {
	r.priceAccumulated += price
	r.accumulatedVolume += volume
	r.nbAccumulatedPrices += 1.0
}

func (r *StatsBuffer) avgPrice() float64 {
	return r.priceAccumulated / r.nbAccumulatedPrices
}

func (r *StatsBuffer) resetBuffer() {
	r.priceAccumulated = 0.0
	r.accumulatedVolume = 0.0
	r.nbAccumulatedPrices = 0.0
}

func main() {
	bs, err := ioutil.ReadFile("coinbaseEUR.csv")

	if err != nil {
		return
	}


	bitcoinCSV := string(bs)

	//read the content of the CSV fields
	r := csv.NewReader(strings.NewReader(bitcoinCSV))

	//Graph parameters
	// n  number of data points from the CSV file  to be considered in the graph
	n := 14094495
	//n := 100
	timeInterval := 1000.0
	outputfileName := "bitcoin1000ms-dual.png"
	//
	//Graph parameters end
	//

	btcTime := make([]float64, n)
	btcPrice := make([]float64, n)
	btcVolume := make([]float64, n)

	//reads first record, to setup the initial timestamp
	record, err := r.Read()

	if err == io.EOF {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	initialTimestamp, err := strconv.ParseFloat(record[0], 64)
	currentTimestamp := 0.0
	if err != nil {
		panic(err)
	}
	btcStatsBuffer := StatsBuffer{accumulatedVolume:0.0, nbAccumulatedPrices: 0.0, priceAccumulated:0.0}
	index := 0
	for {
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		currentTimestamp, err = strconv.ParseFloat(record[0], 64)
		price, err := strconv.ParseFloat(record[1], 64)
		volume, err := strconv.ParseFloat(record[2], 64)

		btcStatsBuffer.accumulate(price, volume)

		if currentTimestamp - initialTimestamp <= timeInterval {
			continue
		} else {
			btcTime[index] = currentTimestamp
			btcPrice[index] = btcStatsBuffer.avgPrice()
			btcVolume[index] = btcStatsBuffer.accumulatedVolume
			index++
			initialTimestamp = currentTimestamp

			btcStatsBuffer.resetBuffer()
		}

		if err != nil {
			panic(err)
		}

		if i > n-2 {
			break
		}

		i++
	}

	graph := chart.Chart {
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true, //enables / displays the x-axis
			},
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				x := v.(float64)
				if x >= math.MaxInt64 || x <= math.MinInt64 { // <-- this works !
					fmt.Println("Conversion impossible: x is out of int64 range.")
					return ""
				}

				typedDate := time.Unix(int64(x), 0)
				return fmt.Sprintf("%d-%d-%d", typedDate.Day(), typedDate.Month(), typedDate.Year())
			},
		},

		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true, //enables / displays the y-axis
			},
		},

		//YAxisSecondary: chart.YAxis{
		//	Style: chart.Style{
		//		Show: true, //enables / displays the secondary y-axis
		//	},
		//},

		Series: []chart.Series {

			chart.ContinuousSeries{
				XValues: btcTime[0:index],
				YValues: btcPrice[0:index],
			},

			//chart.ContinuousSeries{
			//	YAxis: chart.YAxisSecondary,
			//	XValues: btcTime[0:index],
			//	YValues: btcVolume[0:index],
			//},
		},
		Title: "BTC-EUR",
		Width: 32000,
		Height: 1500,
	}

	//buffer := bytes.NewBuffer([]byte{})
	collector := &chart.ImageWriter{}
	err = graph.Render(chart.PNG, collector)
	if err != nil {
		panic(err)
	}


	image, err := collector.Image()
	if err != nil {
		panic(err)
	}

	f, err := os.Create(outputfileName)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, image); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("File Saved.")
}