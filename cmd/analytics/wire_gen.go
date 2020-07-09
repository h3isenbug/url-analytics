// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/h3isenbug/url-analytics/jobs"
	"github.com/h3isenbug/url-common/pkg/message-queue-consumer"
	"github.com/h3isenbug/url-common/pkg/worker-pool"
	"net/http"
)

import (
	_ "github.com/lib/pq"
)

// Injectors from inject_app.go:

func wireUp() (*App, func(), error) {
	logService := provideLogService()
	workerPool, cleanup := provideWorkerPool(logService)
	db, cleanup2, err := provideSQLXConnection()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	totalViewsRepository, err := provideTotalArchiveRepository(db)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	client, cleanup3 := provideRedisClient()
	todayViewsRepository := provideTodayViewsRepository(client, totalViewsRepository)
	uniqueViewsRepository, err := provideUniqueRepository(db)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	service := provideAnalyticsService(totalViewsRepository, todayViewsRepository, uniqueViewsRepository)
	jobRunner, cleanup4 := provideJobRunner(logService, service)
	reportHandler := provideReportHandler(logService, service)
	router := provideMuxRouter(reportHandler)
	server, cleanup5 := provideHTTPServer(router)
	channel, cleanup6, err := provideAMQPChannel()
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	v := provideMessageHandlers(logService, service)
	messageMux := provideMessageMux(workerPool, v)
	messageQueueConsumer, cleanup7, err := provideMessageConsumer(channel, messageMux)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	app := provideApp(workerPool, jobRunner, server, messageQueueConsumer)
	return app, func() {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// inject_app.go:

type App struct {
	httpServer      *http.Server
	jobRunner       *jobs.JobRunner
	messageConsumer consumer.MessageQueueConsumer
	workerPool      pool.WorkerPool
}

func provideApp(workerPool pool.WorkerPool, jobRunner *jobs.JobRunner, httpServer *http.Server, messageConsumer consumer.MessageQueueConsumer) *App {
	return &App{
		workerPool:      workerPool,
		httpServer:      httpServer,
		messageConsumer: messageConsumer,
		jobRunner:       jobRunner,
	}
}
