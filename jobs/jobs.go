package jobs

import (
	"github.com/go-co-op/gocron"
	"github.com/h3isenbug/url-analytics/services/analytics"
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
	"time"
)

type JobRunner struct {
	scheduler *gocron.Scheduler
}

func NewJobRunner(logService log2.LogService, analyticsService analytics.Service) *JobRunner {
	var scheduler = gocron.NewScheduler(time.Now().Location())

	scheduler.Every(1).Day().At("00:30").Do(func() {
		if err := analyticsService.ArchiveYesterday(); err != nil {
			logService.Error("error while archiving yesterday: %s", err.Error())
		}
	})

	scheduler.Every(3).Hours().Do(func() {
		if err := analyticsService.CreateReports(); err != nil {
			logService.Error("error while creating unique reports: %s", err.Error())
		}
	})

	return &JobRunner{scheduler: scheduler}
}

func (runner JobRunner) Start() {
	runner.scheduler.StartAsync()
}

func (runner JobRunner) Stop() {
	runner.scheduler.Stop()
}
