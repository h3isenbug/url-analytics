//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/h3isenbug/url-analytics/jobs"
	consumer "github.com/h3isenbug/url-common/pkg/message-queue-consumer"
	pool "github.com/h3isenbug/url-common/pkg/worker-pool"
	"net/http"
)

type App struct {
	httpServer      *http.Server
	jobRunner       *jobs.JobRunner
	messageConsumer consumer.MessageQueueConsumer
	workerPool      pool.WorkerPool
}

func wireUp() (*App, func(), error) {
	wire.Build(provideJobRunner, provideApp, provideHTTPServer, provideMessageConsumer, provideAMQPChannel, provideMessageMux, provideWorkerPool, provideMessageHandlers, provideLogService, provideReportHandler, provideMuxRouter, provideAnalyticsService, provideUniqueRepository, provideTodayViewsRepository, provideTotalArchiveRepository, provideSQLXConnection, provideRedisClient)

	return &App{}, func() {

	}, nil
}

func provideApp(workerPool pool.WorkerPool, jobRunner *jobs.JobRunner, httpServer *http.Server, messageConsumer consumer.MessageQueueConsumer) *App {
	return &App{
		workerPool:      workerPool,
		httpServer:      httpServer,
		messageConsumer: messageConsumer,
		jobRunner:       jobRunner,
	}
}
