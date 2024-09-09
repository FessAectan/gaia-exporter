package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "0.1.1"

type Config struct {
	Port    int
	GaiaUrl string
}

type Metrics struct {
	HighestBlockNumber      prometheus.Gauge
	CurrentBlockTimeDrift   prometheus.Gauge
	ConnectedPeers          prometheus.Gauge
	ConnectedPeersByVersion *prometheus.GaugeVec
}

type Application struct {
	Config *Config
	Logger *slog.Logger
	Prom   *Metrics
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &Application{
		Config: &Config{},
		Logger: logger,
		Prom: &Metrics{
			HighestBlockNumber: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "gaia_exporter",
				Name:      "highest_block_number",
				Help:      "Highest block number (latest_block_height from /status json).",
			}),
			CurrentBlockTimeDrift: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "gaia_exporter",
				Name:      "current_block_time_drift",
				Help:      "Current time in seconds minus latest_block_time in seconds from /status json.",
			}),
			ConnectedPeers: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "gaia_exporter",
				Name:      "connected_peers",
				Help:      "How many peers are now connected to the node.",
			}),
			ConnectedPeersByVersion: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "gaia_exporter",
				Name:      "connected_peers_by_version",
				Help:      "How many peers are now connected to the node by version.",
			},
				[]string{"version"}),
		},
	}

	flag.IntVar(&app.Config.Port, "port", 4000, "gaia-exporter port")
	flag.StringVar(&app.Config.GaiaUrl, "gaiaUrl", "http://127.0.0.1:26657", "Gaia URL")
	flag.Parse()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting gaia-exporter", "addr", srv.Addr, "version", version, "Gaia URL", app.Config.GaiaUrl)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
