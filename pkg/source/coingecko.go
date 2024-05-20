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

var _ Source = &CoinGeckoSource{}

type CoinGeckoPrice struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"current_price"`
}

// CoinGeckoSource represents a source for fetching prices from CoinGecko.
type CoinGeckoSource struct {
	Host      string
	Threshold float64
}

// NewCoinGeckoSource initializes a new CoinGeckoSource.
func NewCoinGeckoSource(config config.Config) Source {
	return &CoinGeckoSource{
		Host:      config.Sources.CoinGeckoHost,
		Threshold: config.Sources.CoinGeckoThreshold,
	}
}

// Name returns the name of the data source, which is "coingecko".
func (coinGeckoSource *CoinGeckoSource) Name() string {
	return "coingecko"
}

// Start begins the process of fetching prices, building isolation trees, and detecting anomalies.
// It sends errors, anomalies, and statistics through their respective channels.
func (coingeckoSource *CoinGeckoSource) Start(errorChan chan error, doneChan chan struct{}, anomalyChan chan string, statsChan chan string) {
	go func() {
		defer close(doneChan) // Ensure doneChan is closed when the goroutine exits
		for {
			prices, err := coingeckoSource.FetchPrices()
			if err != nil {
				errorChan <- fmt.Errorf("fetching prices: %v, %s", "CoinGeckoSource", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if len(prices) == 0 {
				errorChan <- fmt.Errorf("no prices fetched")
				time.Sleep(5 * time.Second)
				continue
			}

			slog.Debug("Fetched prices:", prices)

			rows := len(prices)
			data := mat.NewDense(rows, 1, nil)
			for i, price := range prices {
				data.Set(i, 0, price.Price)
			}

			numTrees := 100
			trees := make([]*procedure.IsolationTree, numTrees)
			height := int(math.Ceil(math.Log2(float64(rows))))
			for i := 0; i < numTrees; i++ {
				trees[i] = procedure.BuildTree(data, height)
			}

			scores := detect.DetectAnomalies(data, trees)

			// Determine threshold for anomalies
			sortedScores := append([]float64(nil), scores...)
			sort.Float64s(sortedScores)

			// Set the threshold to the 99.99th percentile (top 0.1% anomalies)
			thresholdIndex := int(0.9999 * float64(len(sortedScores)))
			if thresholdIndex >= len(sortedScores) {
				thresholdIndex = len(sortedScores) - 1
			}
			threshold := sortedScores[thresholdIndex]

			slog.Info("anomaly score", "source", coingeckoSource.Name(), "threshold", threshold)

			anomalyCount := 0
			for i, score := range scores {
				if score >= threshold {
					anomalyCount++
					anomalyChan <- fmt.Sprintf("source %s: data point %d - symbol: %s, value: %.2f, anomaly score: %.2f", coingeckoSource.Name(), i, prices[i].Symbol, data.At(i, 0), score)
				}
			}

			statsChan <- fmt.Sprintf("source:%s ,total items: %d, Anomalies: %d", coingeckoSource.Name(), rows, anomalyCount)

			time.Sleep(5 * time.Second)
		}
	}()
}

// FetchPrices is a placeholder function to fetch prices from Binance.
func (coinGeckoSource *CoinGeckoSource) FetchPrices() ([]Price, error) {
	resp, err := http.Get(coinGeckoSource.Host)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var coinGeckoPrice []CoinGeckoPrice
	if err := json.NewDecoder(resp.Body).Decode(&coinGeckoPrice); err != nil {
		return nil, err
	}
	prices := convertCoinGeckoToPrice(coinGeckoPrice)
	return prices, nil
}

// Conversion function
func convertCoinGeckoToPrice(coinGeckoPrices []CoinGeckoPrice) []Price {
	prices := make([]Price, len(coinGeckoPrices))
	for i, cgPrice := range coinGeckoPrices {
		prices[i] = Price{
			Symbol: cgPrice.Symbol,
			Price:  cgPrice.Price,
		}
	}
	return prices
}
