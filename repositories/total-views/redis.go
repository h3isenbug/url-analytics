package total_views

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/h3isenbug/url-analytics/config"
	"github.com/h3isenbug/url-analytics/repositories"
	"strconv"
	"strings"
)

type RedisTodayViewsRepository struct {
	redisClient *redis.Client
	archiveRepo TotalViewsRepository
}

func NewRedisTodayViewsRepository(redisClient *redis.Client, archiveRepo TotalViewsRepository) TodayViewsRepository {
	return &RedisTodayViewsRepository{redisClient: redisClient, archiveRepo: archiveRepo}
}

var browserNames = map[repositories.Browser]string{
	repositories.BrowserChrome:  "chrome",
	repositories.BrowserIE:      "ie",
	repositories.BrowserSafari:  "safari",
	repositories.BrowserFirefox: "firefox",
}

var platformNames = map[repositories.Platform]string{
	repositories.PlatformDesktop: "desktop",
	repositories.PlatformMobile:  "mobile",
}

func (repo RedisTodayViewsRepository) AddView(
	shortPath string,
	browser repositories.Browser,
	platform repositories.Platform,
	day int,
) error {
	if err := repo.redisClient.HIncrBy(strconv.Itoa(day)+"_"+shortPath, "views", 1).Err(); err != nil {
		return err
	}

	if browser != repositories.BrowserUnknown {
		if err := repo.redisClient.HIncrBy(strconv.Itoa(day)+"_"+shortPath, "b_"+browserNames[browser], 1).Err(); err != nil {
			return err
		}
	}

	if platform != repositories.PlatformUnknown {
		if err := repo.redisClient.HIncrBy(strconv.Itoa(day)+"_"+shortPath, "p_"+platformNames[platform], 1).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (repo RedisTodayViewsRepository) MoveViews(day int) error {
	if err := repo.archiveRepo.DeleteViewsOlderThan(day - 30); err != nil {
		return err
	}

	iter := repo.redisClient.Scan(0, strconv.Itoa(day)+"_*", 0).Iterator()
	for iter.Next() {
		key := iter.Val()
		shortPath := strings.Split(key, "_")[1]

		result, err := repo.redisClient.HGetAll(key).Result()
		if err != nil {
			return err
		}

		chrome, _ := strconv.Atoi(result["b_chrome"])
		firefox, _ := strconv.Atoi(result["b_firefox"])
		ie, _ := strconv.Atoi(result["b_ie"])
		safari, _ := strconv.Atoi(result["b_safari"])

		desktop, _ := strconv.Atoi(result["p_desktop"])
		mobile, _ := strconv.Atoi(result["p_mobile"])

		views, _ := strconv.Atoi(result["views"])

		if err := repo.archiveRepo.AddViews(shortPath, chrome, ie, safari, firefox, desktop, mobile, views, day); err != nil {
			return err
		}

		if err := repo.redisClient.Del(key).Err(); err != nil {
			return err
		}

	}

	return nil
}

func (repo RedisTodayViewsRepository) GetReport(shortPath string) (*repositories.Report, error) {
	var today = strconv.Itoa(config.DaysSince2020())
	result, err := repo.redisClient.HGetAll(today + "_" + shortPath).Result()
	if errors.Is(err, redis.Nil) {
		return &repositories.Report{}, nil
	}
	if err != nil {
		return nil, err
	}

	return &repositories.Report{
		BrowserChrome:   atoi(result["b_chrome"]),
		BrowserIE:       atoi(result["b_ie"]),
		BrowserSafari:   atoi(result["b_safari"]),
		BrowserFirefox:  atoi(result["b_firefox"]),
		PlatformDesktop: atoi(result["p_desktop"]),
		PlatformMobile:  atoi(result["p_mobile"]),
		TotalViews:      atoi(result["views"]),
	}, nil
}

func atoi(ascii string) int {
	var i, _ = strconv.Atoi(ascii)
	return i
}
