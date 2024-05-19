package source

type Source interface {
	Name() string
	FetchPrices() ([]Price, error)
	Start(errCh chan error, doneCh chan struct{}, anomalyChan chan string, statsChan chan string)
}

// Price represents a price data point.
type Price struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}
