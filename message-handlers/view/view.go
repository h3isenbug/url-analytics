package view

import (
	"encoding/json"
	"github.com/avct/uasurfer"
	"github.com/h3isenbug/url-analytics/repositories"
	"github.com/h3isenbug/url-analytics/services/analytics"
	"github.com/h3isenbug/url-common/pkg/messages"
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
)

var browserConvert = map[uasurfer.BrowserName]repositories.Browser{
	uasurfer.BrowserChrome:  repositories.BrowserChrome,
	uasurfer.BrowserIE:      repositories.BrowserIE,
	uasurfer.BrowserSafari:  repositories.BrowserSafari,
	uasurfer.BrowserFirefox: repositories.BrowserFirefox,
}
var platformConvert = map[uasurfer.DeviceType]repositories.Platform{
	uasurfer.DeviceComputer: repositories.PlatformDesktop,
	uasurfer.DevicePhone:    repositories.PlatformMobile,
	uasurfer.DeviceTablet:   repositories.PlatformMobile,
}

type AnalyticsHandler struct {
	analyticsService analytics.Service
	log              log2.LogService
}

func NewAnalyticsMessageHandler(analyticsService analytics.Service, log log2.LogService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService, log: log}
}

func (handler AnalyticsHandler) ViewMessageHandler(body []byte) error {
	var message messages.URLViewedEvent
	if err := json.Unmarshal(body, &message); err != nil {
		return err
	}

	var userAgent = uasurfer.Parse(message.UserAgent)
	browser, found := browserConvert[userAgent.Browser.Name]
	if !found {
		browser = repositories.BrowserUnknown
	}

	platform, found := platformConvert[userAgent.DeviceType]
	if !found {
		platform = repositories.PlatformUnknown
	}

	return handler.analyticsService.AddToAnalytics(message.ShortPath, message.ETag, browser, platform)
}
