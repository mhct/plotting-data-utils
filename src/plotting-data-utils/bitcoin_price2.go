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
	//n := 1000
	timeInterval := 200.0
	outputfileName := "bitcoin200ms.png"
	//
	//Graph parameters end
	//

	btcTime := make([]float64, n)
	btcPrice := make([]float64, n)
	//btcVolumePoints := make(plotter.XYs, n)

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
	priceAccumulated := 0.0
	nbAccumulatedPrices := 0.0
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
		//volume, err := strconv.ParseFloat(record[2], 64)

		priceAccumulated += price
		nbAccumulatedPrices = nbAccumulatedPrices + 1.0

		if currentTimestamp - initialTimestamp <= timeInterval {
			continue
		} else {
			btcTime[index] = currentTimestamp
			btcPrice[index] = priceAccumulated / nbAccumulatedPrices
			index++
			initialTimestamp = currentTimestamp
			priceAccumulated = 0.0
			nbAccumulatedPrices = 0.0
		}

		if err != nil {
			panic(err)
		}

		if i > n-2 {
			break
		}

		i++
	}
	//fmt.Println(btcPrice)

	//viridisByY := func(xr, yr chart.Range, index int, x, y float64) drawing.Color {
	//	return chart.Viridis(y, yr.GetMin(), yr.GetMax())
	//}
	fmt.Println(bla(int64(1429753354)))

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

		Series: []chart.Series {

			chart.ContinuousSeries{
				XValues: btcTime[0:index],
				YValues: btcPrice[0:index],
			},
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


func bla(v interface{}) string {
	typed := v.(int64)
	typedDate := time.Unix(typed, 0)
	return fmt.Sprintf("%d-%d\n%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
}