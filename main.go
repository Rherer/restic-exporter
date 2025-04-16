package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Settings struct {
	HTTP_BASE_PATH string        `env:"HTTP_BASE_PATH" envDefault:"/metrics"`
	HTTP_BASE_PORT int           `env:"HTTP_BASE_PORT" envDefault:"8080"`
	CHECK_INTERVAL time.Duration `env:"CHECK_INTERVAL" envDefault:"30m"`
}

var Config Settings = Settings{}
var checkResult int = -1

func main() {
	getSettings()
	// Testing
	os.Setenv("RESTIC_REPOSITORY", "/mnt/data/distrobox/restic-exporter/home/restic")
	os.Setenv("RESTIC_PASSWORD", "1")

	// Run a check at first, we can safely ignore the errors on this function
	checkResult, _ = runCheck()

	// Register a Ticket for periodic checks
	registerTicker()

	// Prometheus handles all calls for us (Using the custom Collector)
	registerHTTP()
}

// Only settings for the exporter, other ENV-Variables get passed to restic
// See struct on top of this file, for all options
func getSettings() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
}

func registerTicker() {
	ticker := time.NewTicker(Config.CHECK_INTERVAL)

	go func(ticker *time.Ticker, checkResult *int) {
		for _ = range ticker.C {
			result, err := runCheck()
			if err != nil {
				fmt.Println("Check failed with error:", err)
			}
			*checkResult = result
		}
	}(ticker, &checkResult)
}

// Register custom metrics collector and start the http server
func registerHTTP() {
	registry := prometheus.NewPedanticRegistry()
	collector := newCollector()
	registry.MustRegister(collector)

	http.Handle(Config.HTTP_BASE_PATH, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))
	fmt.Println("Serving metrics at :" + strconv.Itoa(Config.HTTP_BASE_PORT) + Config.HTTP_BASE_PATH)
	err := http.ListenAndServe(":"+strconv.Itoa(Config.HTTP_BASE_PORT), nil)
	if err != nil {
		panic(err)
	}
}
