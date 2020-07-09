package analytics

import (
	"github.com/h3isenbug/url-analytics/config"
	"github.com/h3isenbug/url-analytics/repositories"
	total "github.com/h3isenbug/url-analytics/repositories/total-views"
	unique "github.com/h3isenbug/url-analytics/repositories/unique-views"
)

type Service interface {
	AddToAnalytics(shortPath, etag string, browser repositories.Browser, platform repositories.Platform) error
	CreateReports() error
	ArchiveYesterday() error
	GetReports(shortPath string) (*Reports, error)
}

type ServiceV1 struct {
	todayRepo  total.TodayViewsRepository
	totalRepo  total.TotalViewsRepository
	uniqueRepo unique.UniqueViewsRepository
}

type Reports struct {
	UniqueToday     *repositories.Report `json:"uniqueToday"`
	UniqueYesterday *repositories.Report `json:"uniqueYesterday"`
	UniqueLastWeek  *repositories.Report `json:"uniqueLastWeek"`
	UniqueLastMonth *repositories.Report `json:"uniqueLastMonth"`

	TotalToday     *repositories.Report `json:"totalToday"`
	TotalYesterday *repositories.Report `json:"totalYesterday"`
	TotalLastWeek  *repositories.Report `json:"totalLastWeek"`
	TotalLastMonth *repositories.Report `json:"totalLastMonth"`
}

func NewServiceV1(totalRepo total.TotalViewsRepository, todayRepo total.TodayViewsRepository, uniqueRepo unique.UniqueViewsRepository) Service {
	return &ServiceV1{totalRepo: totalRepo, todayRepo: todayRepo, uniqueRepo: uniqueRepo}
}

func (service ServiceV1) AddToAnalytics(shortPath, etag string, browser repositories.Browser, platform repositories.Platform) error {
	if err := service.todayRepo.AddView(shortPath, browser, platform, config.DaysSince2020()); err != nil {
		return err
	}

	if err := service.uniqueRepo.AddView(shortPath, etag, browser, platform, config.DaysSince2020()); err != nil {
		return err
	}

	return nil
}

func (service ServiceV1) CreateReports() error {
	var today = config.DaysSince2020()

	if err := service.uniqueRepo.CreateAndStoreReports(repositories.ReportToday, today); err != nil {
		return err
	}
	if err := service.uniqueRepo.CreateAndStoreReports(repositories.ReportYesterday, today); err != nil {
		return err
	}
	if err := service.uniqueRepo.CreateAndStoreReports(repositories.ReportLastWeek, today); err != nil {
		return err
	}
	if err := service.uniqueRepo.CreateAndStoreReports(repositories.ReportLastMonth, today); err != nil {
		return err
	}

	return nil
}

func (service ServiceV1) ArchiveYesterday() error {
	return service.todayRepo.MoveViews(config.DaysSince2020() - 1)
}

func (service ServiceV1) GetReports(shortPath string) (*Reports, error) {
	uniqueToday, err := service.uniqueRepo.GetReport(shortPath, repositories.ReportToday)
	if err != nil {
		return nil, err
	}

	uniqueYesterday, err := service.uniqueRepo.GetReport(shortPath, repositories.ReportYesterday)
	if err != nil {
		return nil, err
	}

	uniqueLastWeek, err := service.uniqueRepo.GetReport(shortPath, repositories.ReportLastWeek)
	if err != nil {
		return nil, err
	}

	uniqueLastMonth, err := service.uniqueRepo.GetReport(shortPath, repositories.ReportLastMonth)
	if err != nil {
		return nil, err
	}

	totalToday, err := service.todayRepo.GetReport(shortPath)
	if err != nil {
		return nil, err
	}

	var today = config.DaysSince2020()
	totalYesterday, err := service.totalRepo.GetReport(shortPath, today-1, today-1)
	if err != nil {
		return nil, err
	}
	totalLastWeek, err := service.totalRepo.GetReport(shortPath, today-7, today-1)
	if err != nil {
		return nil, err
	}
	totalLastMonth, err := service.totalRepo.GetReport(shortPath, today-30, today-1)
	if err != nil {
		return nil, err
	}

	return &Reports{
		UniqueToday:     uniqueToday,
		UniqueYesterday: uniqueYesterday,
		UniqueLastWeek:  uniqueLastWeek,
		UniqueLastMonth: uniqueLastMonth,
		TotalToday:      totalToday,
		TotalYesterday:  totalYesterday,
		TotalLastWeek:   totalLastWeek,
		TotalLastMonth:  totalLastMonth,
	}, nil
}
