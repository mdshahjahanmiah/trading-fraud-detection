package main

import (
	"github.com/mdshahjahanmiah/fraud-detection/pkg/config"
	"github.com/mdshahjahanmiah/fraud-detection/pkg/source"
	"log/slog"
	"strings"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		slog.Error("error", "err", err)
		return
	}

	errCh := make(chan error)
	doneCh := make(chan struct{})
	anomalyChan := make(chan string)
	statsChan := make(chan string)

	enabledSources := strings.Split(conf.EnabledSources, ";")
	for _, src := range enabledSources {
		switch strings.TrimSpace(src) {
		case "binance":
			binanceSource := source.NewBinanceSource(conf)
			go binanceSource.Start(errCh, doneCh, anomalyChan, statsChan)
		case "coingecko":
			coinGeckoSource := source.NewCoinGeckoSource(conf)
			go coinGeckoSource.Start(errCh, doneCh, anomalyChan, statsChan)
		default:
			slog.Warn("unknown source", "source", src)
		}
	}

	for {
		select {
		case err := <-errCh:
			slog.Error("error", "err", err)
		case anomaly := <-anomalyChan:
			slog.Warn("found anomalies", "anomaly", anomaly)
		case stats := <-statsChan:
			slog.Info("anomalies statistics", "stats", stats)
		case <-doneCh:
			slog.Info("processing completed")
			return
		}
	}
}
