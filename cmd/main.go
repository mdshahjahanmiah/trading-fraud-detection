package main

import (
	"github.com/mdshahjahanmiah/fraud-detection/pkg/config"
	"github.com/mdshahjahanmiah/fraud-detection/pkg/source"
	"log/slog"
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

	binanceSource := source.NewBinanceSource(conf)
	go binanceSource.Start(errCh, doneCh, anomalyChan, statsChan)

	for {
		select {
		case err := <-errCh:
			slog.Error("error", "err", err)
		case anomaly := <-anomalyChan:
			slog.Info("found anomalies", "anomaly", anomaly)
		case stats := <-statsChan:
			slog.Info("statistics", "stats", stats)
		case <-doneCh:
			slog.Info("processing completed")
			return
		}
	}
}
