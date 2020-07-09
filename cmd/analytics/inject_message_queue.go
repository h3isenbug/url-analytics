package main

import (
	"github.com/h3isenbug/url-analytics/config"
	"github.com/h3isenbug/url-analytics/message-handlers/view"
	"github.com/h3isenbug/url-analytics/services/analytics"
	mux "github.com/h3isenbug/url-common/pkg/event-mux"
	consumer "github.com/h3isenbug/url-common/pkg/message-queue-consumer"
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
	pool "github.com/h3isenbug/url-common/pkg/worker-pool"
	"github.com/streadway/amqp"
	"time"
)

func provideMessageHandlers(logService log2.LogService, analyticsService analytics.Service) map[string]func([]byte) error {
	var messageHandler = view.NewAnalyticsMessageHandler(analyticsService, logService)

	return map[string]func([]byte) error{
		"view": messageHandler.ViewMessageHandler,
	}
}

func provideWorkerPool(logService log2.LogService) (pool.WorkerPool, func()) {
	workerPool := pool.NewWorkerPoolV1(logService, config.Config.WorkerCount, time.Millisecond*50)
	return workerPool, workerPool.GracefulShutdown
}

func provideMessageMux(workerPool pool.WorkerPool, handlers map[string]func([]byte) error) mux.MessageMux {
	var messageMux = mux.NewMessageMuxV1(workerPool)
	for k, v := range handlers {
		messageMux.SetHandler(k, v)
	}

	return messageMux
}
func provideAMQPChannel() (*amqp.Channel, func(), error) {
	con, err := amqp.Dial(config.Config.RabbitServer)
	if err != nil {
		return nil, nil, err
	}

	channel, err := con.Channel()

	return channel, func() {
		channel.Close()
		con.Close()
	}, err
}

func provideMessageConsumer(channel *amqp.Channel, messageMux mux.MessageMux) (consumer.MessageQueueConsumer, func(), error) {
	var cons, err = consumer.NewRabbitMQQueueConsumerV1(channel, func(tag uint64) { channel.Ack(tag, false) }, messageMux, config.Config.RabbitQueueName)
	return cons, func() { cons.GracefulShutdown() }, err
}
