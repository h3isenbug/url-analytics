package total_views

import (
	"github.com/h3isenbug/url-analytics/repositories"
)

type TotalViewsRepository interface {
	AddViews(
		shortPath string,

		browserChrome, browserIE, browserSafari, browserFirefox,

		platformDesktop, platformMobile,

		totalViews int,
		day int,
	) error

	DeleteViewsOlderThan(day int) error
	GetReport(shortPath string, fromDay, toDay int) (*repositories.Report, error)
}

type TodayViewsRepository interface {
	AddView(
		shortPath string,
		browser repositories.Browser,
		platform repositories.Platform,
		day int,
	) error

	MoveViews(day int) error
	GetReport(shortPath string) (*repositories.Report, error)
}
