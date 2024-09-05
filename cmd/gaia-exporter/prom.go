package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func (app *Application) promHandler() http.Handler {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		app.Prom.HighestBlockNumber,
		app.Prom.CurrentBlockTimeDrift,
		app.Prom.ConnectedPeers,
		app.Prom.ConnectedPeersByVersion,
	)

	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}
