package source

import (
	"encoding/json"
	"fmt"
	"github.com/mdshahjahanmiah/fraud-detection/pkg/config"
	"github.com/mdshahjahanmiah/fraud-detection/pkg/detect"
	"github.com/mdshahjahanmiah/fraud-detection/pkg/procedure"
	"gonum.org/v1/gonum/mat"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"time"
)

var _ Source = &BinanceSource{}

// BinancePrice represents a price data point.
type BinancePrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

// BinanceSource represents a source for fetching prices from Binance.
type BinanceSource struct {
	Host      string
	Threshold float64
}

// NewBinanceSource initializes a new BinanceSource.
func NewBinanceSource(config config.Config) Source {
	return &BinanceSource{
		Host:      config.Sources.BinanceHost,
		Threshold: config.Sources.BinanceThreshold,
	}
}

// Name returns the name of the data source, which is "binance".
func (binanceSource *BinanceSource) Name() string {
	return "binance"
}

// Start begins the process of fetching prices, building isolation trees, and detecting anomalies.
// It sends errors, anomalies, and statistics through their respective channels.
func (binanceSource *BinanceSource) Start(errorChan chan error, doneChan chan struct{}, anomalyChan chan string, statsChan chan string) {
	go func() {
		defer close(doneChan) // Ensure doneChan is closed when the goroutine exits
		for {
			// Fetch prices from Binance
			prices, err := binanceSource.FetchPrices()
			if err != nil {
				errorChan <- fmt.Errorf("fetching prices: %v, %s", binanceSource.Name(), err)
				time.Sleep(5 * time.Second)
				continue
			}

			// Check if any prices were fetched
			if len(prices) == 0 {
				errorChan <- fmt.Errorf("no prices fetched: %s", binanceSource.Name())
				time.Sleep(5 * time.Second)
				continue
			}

			slog.Debug("Fetched prices", "prices", prices)

			// Create a data matrix from the fetched prices
			rows := len(prices)
			data := mat.NewDense(rows, 1, nil)
			for i, price := range prices {
				data.Set(i, 0, price.Price)
			}

			// Build isolation trees for anomaly detection
			numTrees := 100
			trees := make([]*procedure.IsolationTree, numTrees)
			height := int(math.Ceil(math.Log2(float64(rows))))
			for i := 0; i < numTrees; i++ {
				trees[i] = procedure.BuildTree(data, height)
			}

			// Detect anomalies using the isolation trees
			scores := detect.DetectAnomalies(data, trees)

			// Sort the anomaly scores to determine the threshold
			sortedScores := append([]float64(nil), scores...)
			sort.Float64s(sortedScores)

			// Set the threshold to the 99th percentile (top 1% anomalies)
			thresholdIndex := int(binanceSource.Threshold * float64(len(sortedScores)))
			if thresholdIndex >= len(sortedScores) {
				thresholdIndex = len(sortedScores) - 1
			}
			threshold := sortedScores[thresholdIndex]

			slog.Info("anomaly score", "source", binanceSource.Name(), "threshold", threshold)

			// Count and report anomalies
			anomalyCount := 0
			for i, score := range scores {
				if score >= threshold {
					anomalyCount++
					anomalyChan <- fmt.Sprintf("source %s: data point %d - symbol: %s, value: %.2f, anomaly score: %.2f", binanceSource.Name(), i, prices[i].Symbol, data.At(i, 0), score)
				}
			}

			// Send statistics through the stats channel
			statsChan <- fmt.Sprintf("source: %s, total items: %d, anomalies: %d", binanceSource.Name(), rows, anomalyCount)

			time.Sleep(5 * time.Second)
		}
	}()
}

// FetchPrices is a placeholder function to fetch prices from Binance.
func (binanceSource *BinanceSource) FetchPrices() ([]Price, error) {
	resp, err := http.Get(binanceSource.Host)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var prices []Price
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, err
	}
	return prices, nil
}
