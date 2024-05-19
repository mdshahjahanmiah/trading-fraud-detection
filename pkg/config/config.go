package config

import "flag"

type Sources struct {
	BinanceHost string
	Threshold   float64
}

type Config struct {
	Sources        Sources
	EnabledSources string
}

func Load() (Config, error) {
	flagSet := flag.NewFlagSet("", flag.ExitOnError)

	var sources Sources
	flagSet.StringVar(&sources.BinanceHost, "binance.source.hostname", "https://api.binance.com/api/v3/ticker/price", "binance source host name")
	flagSet.Float64Var(&sources.Threshold, "binance.source.threshold", 0.9995, "set the threshold to the 99th percentile (top 1% anomalies)")
	enabledSources := flagSet.String("anomalies.enabled.sources", "binance", "semi-colon delimited specifying enabled sources")

	config := Config{
		Sources:        sources,
		EnabledSources: *enabledSources,
	}

	return config, nil
}
