package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/h3isenbug/url-analytics/config"
	handlers "github.com/h3isenbug/url-analytics/http-handlers"
	"github.com/h3isenbug/url-analytics/http-handlers/report"
	"github.com/h3isenbug/url-analytics/services/analytics"
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
	"net/http"
)

func provideHTTPServer(router *mux.Router) (*http.Server, func()) {
	server := &http.Server{
		Addr:    ":" + config.Config.Port,
		Handler: router,
	}
	return server, func() { server.Shutdown(context.Background()) }
}

func provideReportHandler(logService log2.LogService, analyticsService analytics.Service) report.ReportHandler {
	return report.NewReportHandlerV1(logService, analyticsService)
}

func provideMuxRouter(reportHandler report.ReportHandler) *mux.Router {
	router := mux.NewRouter()
	router.Use(handlers.GorillaMuxURLParamMiddleware)
	router.Methods("GET").Path("/analytics/{short_path}").HandlerFunc(reportHandler.GetReports)
	return router
}
