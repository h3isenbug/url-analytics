package unique_views

import "github.com/h3isenbug/url-analytics/repositories"

type UniqueViewsRepository interface {
	AddView(shortPath string, etag string, browser repositories.Browser, platform repositories.Platform, day int) error
	CreateAndStoreReports(reportType repositories.ReportType, today int) error
	GetReport(shortPath string, reportType repositories.ReportType) (*repositories.Report, error)
}
