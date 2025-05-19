package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Struct for additional Settings, RESTIC Environment variables get passed to restic directly
type Settings struct {
	HTTP_BASE_PATH string        `env:"HTTP_BASE_PATH" envDefault:"/metrics"`
	HTTP_BASE_PORT int           `env:"HTTP_BASE_PORT" envDefault:"8080"`
	CHECK_INTERVAL time.Duration `env:"CHECK_INTERVAL" envDefault:"30m"`
	NO_CHECK       bool          `env:"NO_CHECK" envDefault:"false"`
	USE_REPO_PATH  bool          `env:"USE_REPO_PATH" envDefault:"false"`
	USE_LATEST_N   int           `env:"USE_LATEST_N" envDefault:"1"`
}

var Config Settings = Settings{}
var checkResult int = -1

func main() {
	getSettings()

	// Immediately panic, if no repo found at location
	checkIfRepoExists()

	// Only run checks, if not disabled via env NO_CHECK
	if !Config.NO_CHECK {
		// Register a Ticker for periodic checks
		registerTicker()
	}

	// Prometheus handles all calls for us (Using the custom Collector)
	// So the metrics will get refreshed on every http request (Except check, see CHECK_INTERVAL)
	registerHTTP()
}

// Only settings for the exporter, other ENV-Variables get passed to restic
// See struct Settings for more information
func getSettings() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
}

// Registers a Ticker with the configured CHECK_INTERVAL to periodically run checks
func registerTicker() {
	ticker := time.NewTicker(Config.CHECK_INTERVAL)

	go func(ticker *time.Ticker, checkResult *int) {
		// Small workaround, so the ticket triggers once at the beginning, see: https://github.com/golang/go/issues/17601
		for ; true; <-ticker.C {
			result, err := runCheck()
			if err != nil {
				fmt.Println("Check failed with error:", err)
			}
			*checkResult = result
		}
	}(ticker, &checkResult)
}

// Register custom metrics collector and start the http server on configured port and path (see: HTTP_BASE_PORT and HTTP_BASE_PATH)
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
