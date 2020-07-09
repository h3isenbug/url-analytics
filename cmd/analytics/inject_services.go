package main

import (
	"github.com/h3isenbug/url-analytics/jobs"
	total "github.com/h3isenbug/url-analytics/repositories/total-views"
	unique "github.com/h3isenbug/url-analytics/repositories/unique-views"
	"github.com/h3isenbug/url-analytics/services/analytics"
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
	"os"
)

func provideAnalyticsService(totalRepo total.TotalViewsRepository, todayRepo total.TodayViewsRepository, uniqueRepo unique.UniqueViewsRepository) analytics.Service {
	return analytics.NewServiceV1(totalRepo, todayRepo, uniqueRepo)
}

func provideLogService() log2.LogService {
	return log2.NewLogServiceV1(os.Stdout)
}

func provideJobRunner(logService log2.LogService, analyticsService analytics.Service) (*jobs.JobRunner, func()) {
	var runner = jobs.NewJobRunner(logService, analyticsService)
	return runner, func() {
		runner.Stop()
	}
}
