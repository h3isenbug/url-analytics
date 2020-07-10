package total_views

import (
	"github.com/h3isenbug/url-analytics/repositories"
	"github.com/jmoiron/sqlx"
)

type PostgresTotalViewRepository struct {
	con *sqlx.DB
}

func NewPostgresTotalViewRepository(con *sqlx.DB) (TotalViewsRepository, error) {
	_, err := con.Exec(`CREATE TABLE IF NOT EXISTS total_views (
									short_path varchar(10),
									browser_chrome int,
									browser_ie int,
									browser_safari int,
									browser_firefox int,
									platform_desktop int,
									platform_mobile int,
									total_views int,
									day_since_2020 int,
									PRIMARY KEY (short_path, day_since_2020));`)
	if err != nil {
		return nil, err
	}

	return &PostgresTotalViewRepository{con: con}, nil
}

func (repo PostgresTotalViewRepository) AddViews(
	shortPath string,

	browserChrome,
	browserIE,
	browserSafari,
	browserFirefox,

	platformDesktop,
	platformMobile,

	totalViews int,

	day int,
) error {
	_, err := repo.con.Exec(`INSERT INTO total_views (
							short_path,
							browser_chrome, browser_ie, browser_safari, browser_firefox,
							platform_desktop, platform_mobile,
							total_views, day_since_2020) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
							ON CONFLICT (short_path, day_since_2020) DO UPDATE SET
								browser_chrome   = $2,
								browser_ie       = $3,
								browser_safari   = $4,
								browser_firefox  = $5,
								platform_desktop = $6,
								platform_mobile  = $7,
								total_views      = $8`,
		shortPath,
		browserChrome, browserIE, browserSafari, browserFirefox,
		platformDesktop, platformMobile,
		totalViews,
		day,
	)
	return err
}

func (repo PostgresTotalViewRepository) GetReport(shortPath string, fromDay, toDay int) (*repositories.Report, error) {
	var report repositories.Report
	var err = repo.con.Get(
		&report,
		`SELECT  COALESCE(SUM(browser_chrome), 0)  AS browser_chrome, 
						COALESCE(SUM(browser_ie), 0)      AS browser_ie, 
						COALESCE(SUM(browser_safari), 0)  AS browser_safari, 
						COALESCE(SUM(browser_firefox), 0) AS browser_firefox,
						 
						COALESCE(SUM(platform_desktop), 0) AS platform_desktop,
						COALESCE(SUM(platform_mobile), 0)  AS platform_mobile,
						COALESCE(SUM(total_views), 0)      AS total_views
				FROM total_views WHERE short_path=$1 AND day_since_2020 >= $2 AND day_since_2020 <= $3;`,
		shortPath, fromDay, toDay)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo PostgresTotalViewRepository) DeleteViewsOlderThan(day int) error {
	var _, err = repo.con.Exec(`DELETE FROM total_views WHERE day_since_2020 < $1`, day)
	return err
}
