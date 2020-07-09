package report

import (
	"encoding/json"
	handlers "github.com/h3isenbug/url-analytics/http-handlers"
	"github.com/h3isenbug/url-analytics/services/analytics"
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
	"net/http"
)

type ReportHandler interface {
	GetReports(w http.ResponseWriter, r *http.Request)
}

type ReportHandlerV1 struct {
	log              log2.LogService
	analyticsService analytics.Service
}

func NewReportHandlerV1(log log2.LogService, analyticsService analytics.Service) ReportHandler {
	return &ReportHandlerV1{log: log, analyticsService: analyticsService}
}

func (handler ReportHandlerV1) GetReports(w http.ResponseWriter, r *http.Request) {
	var params = handlers.GetURLParams(r)
	var reports, err = handler.analyticsService.GetReports(params["short_path"])
	if err != nil {
		handler.log.Error("could not get reports for %s. reason: %s", params["short_path"], err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
