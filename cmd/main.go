package main

import (
	"github.com/mdshahjahanmiah/fraud-detection/pkg/config"
	"github.com/mdshahjahanmiah/fraud-detection/pkg/source"
	"log/slog"
	"strings"
)

func main() {
	// Load configuration settings
	conf, err := config.Load()
	if err != nil {
		slog.Error("error", "err", err)
		return
	}

	// Channels for error handling, completion signal, anomaly detection, and statistics
	errChan := make(chan error)
	doneChan := make(chan struct{})
	anomalyChan := make(chan string)
	statsChan := make(chan string)

	// Parse and initialize enabled data sources from configuration
	enabledSources := strings.Split(conf.EnabledSources, ";")
	for _, src := range enabledSources {
		switch strings.TrimSpace(src) {
		case "binance":
			binanceSource := source.NewBinanceSource(conf)
			go binanceSource.Start(errChan, doneChan, anomalyChan, statsChan)
		case "coingecko":
			coinGeckoSource := source.NewCoinGeckoSource(conf)
			go coinGeckoSource.Start(errChan, doneChan, anomalyChan, statsChan)
		default:
			slog.Warn("unknown source", "source", src)
		}
	}

	// Main loop to handle messages from channels
	for {
		select {
		case err := <-errChan:
			// Log any errors received
			slog.Error("error", "err", err)
		case anomaly := <-anomalyChan:
			// Log any anomalies detected
			slog.Warn("found anomalies", "anomaly", anomaly)
		case stats := <-statsChan:
			// Log statistics about anomalies
			slog.Info("anomalies statistics", "stats", stats)
		case <-doneChan:
			slog.Info("processing completed")
			return
		}
	}
}
