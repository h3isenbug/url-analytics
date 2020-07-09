package repositories

import (
	"errors"
)

type Browser int

const (
	BrowserUnknown Browser = iota
	BrowserChrome
	BrowserIE
	BrowserSafari
	BrowserFirefox
)

type Platform int

const (
	PlatformUnknown Platform = iota
	PlatformDesktop
	PlatformMobile
)

type ReportType int

const (
	ReportToday ReportType = iota + 1
	ReportYesterday
	ReportLastWeek
	ReportLastMonth
)

var (
	ErrNotFound = errors.New("not found")
)

type Report struct {
	BrowserChrome  int `db:"browser_chrome" json:"browserChrome"`
	BrowserIE      int `db:"browser_ie" json:"browserIE"`
	BrowserSafari  int `db:"browser_safari" json:"browserSafari"`
	BrowserFirefox int `db:"browser_firefox" json:"browserFirefox"`

	PlatformDesktop int `db:"platform_desktop" json:"platformDesktop"`
	PlatformMobile  int `db:"platform_mobile" json:"platformMobile"`

	TotalViews int `db:"total_views" json:"totalViews"`
}
