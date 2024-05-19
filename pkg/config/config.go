package config

import (
	"flag"
	"fmt"
)

type Sources struct {
	BinanceHost        string
	CoinGeckoHost      string
	BinanceThreshold   float64
	CoinGeckoThreshold float64
}

type Config struct {
	Sources        Sources
	EnabledSources string
}

func Load() (Config, error) {
	flagSet := flag.NewFlagSet("", flag.ExitOnError)

	var sources Sources
	flagSet.StringVar(&sources.BinanceHost, "binance.source.hostname", "https://api.binance.com/api/v3/ticker/price", "binance source host name")
	flagSet.StringVar(&sources.CoinGeckoHost, "coingecko.source.hostname", "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd", "coingecko source host name")
	flagSet.Float64Var(&sources.BinanceThreshold, "binance.source.threshold", 0.999, "set the threshold for Binance (99.9th percentile)")
	flagSet.Float64Var(&sources.CoinGeckoThreshold, "coingecko.source.threshold", 0.9995, "set the threshold for CoinGecko (99.95th percentile)")

	// Define a variable for enabled sources
	enabledSources := flagSet.String("anomalies.enabled.sources", "binance;coingecko", "semi-colon delimited specifying enabled sources")

	// Parse the command line flags
	err := flagSet.Parse([]string{})
	if err != nil {
		return Config{}, fmt.Errorf("error parsing flags: %v", err)
	}

	config := Config{
		Sources:        sources,
		EnabledSources: *enabledSources,
	}

	return config, nil
}
