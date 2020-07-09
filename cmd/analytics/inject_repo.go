package main

import (
	"github.com/go-redis/redis"
	"github.com/h3isenbug/url-analytics/config"
	total "github.com/h3isenbug/url-analytics/repositories/total-views"
	unique "github.com/h3isenbug/url-analytics/repositories/unique-views"
	"github.com/jmoiron/sqlx"
)

func provideUniqueRepository(con *sqlx.DB) (unique.UniqueViewsRepository, error) {
	return unique.NewPostgresUniqueViewRepository(con)
}

func provideTodayViewsRepository(redis *redis.Client, archiveRepo total.TotalViewsRepository) total.TodayViewsRepository {
	return total.NewRedisTodayViewsRepository(redis, archiveRepo)
}

func provideTotalArchiveRepository(con *sqlx.DB) (total.TotalViewsRepository, error) {
	return total.NewPostgresTotalViewRepository(con)
}

func provideSQLXConnection() (*sqlx.DB, func(), error) {
	var con, err = sqlx.Open("postgres", config.Config.DSN)
	return con, func() { con.Close() }, err
}

func provideRedisClient() (*redis.Client, func()) {
	var client = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisServer,
		Password: config.Config.RedisPassword,
		DB:       config.Config.RedisDB,
	})
	return client, func() {
		client.Close()
	}
}
