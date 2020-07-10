package unique_views

import (
	"database/sql"
	"errors"
	"github.com/h3isenbug/url-analytics/repositories"
	"github.com/jmoiron/sqlx"
)

type PostgresUniqueViewRepository struct {
	con *sqlx.DB
}

func NewPostgresUniqueViewRepository(con *sqlx.DB) (UniqueViewsRepository, error) {
	_, err := con.Exec(`CREATE TABLE IF NOT EXISTS unique_views (short_path varchar(10), etag varchar(40), browser int, platform int, days_since_2020 int, PRIMARY KEY(short_path, etag));`)
	if err != nil {
		return nil, err
	}

	_, err = con.Exec(`CREATE TABLE IF NOT EXISTS unique_views_report (short_path varchar(10), report_type int, browser_chrome int, browser_ie int, browser_safari int, browser_firefox int, platform_desktop int, platform_mobile int, total_views int, updated_at TIMESTAMPTZ, PRIMARY KEY (short_path, report_type));`)
	if err != nil {
		return nil, err
	}

	return &PostgresUniqueViewRepository{con: con}, nil
}

func (repo PostgresUniqueViewRepository) AddView(shortPath string, etag string, browser repositories.Browser, platform repositories.Platform, day int) error {
	_, err := repo.con.Exec(
		`INSERT INTO unique_views (short_path, etag, browser, platform, days_since_2020) 
					VALUES ($1, $2, $3, $4, $5) ON CONFLICT (short_path, etag) DO UPDATE SET days_since_2020=$5`,
		shortPath, etag, browser, platform, day)
	return err
}

func (repo PostgresUniqueViewRepository) CreateAndStoreReports(reportType repositories.ReportType, today int) error {
	var from = map[repositories.ReportType]int{
		repositories.ReportToday: today,
		repositories.ReportYesterday:          today - 1,
		repositories.ReportLastWeek:           today - 7,
		repositories.ReportLastMonth:          today - 30,
	}[reportType]
	var to = map[repositories.ReportType]int{
		repositories.ReportToday:     today,
		repositories.ReportYesterday: today - 1,
		repositories.ReportLastWeek:  today - 1,
		repositories.ReportLastMonth: today - 1,
	}[reportType]

	_, err := repo.con.Exec(
		`
		INSERT INTO unique_views_report
			(SELECT short_path,
					$3 AS report_type,
					COALESCE(SUM(case when browser = 1 then 1 else 0 end), 0) AS browser_chrome,
					COALESCE(SUM(case when browser = 2 then 1 else 0 end), 0) AS browser_ie,
					COALESCE(SUM(case when browser = 3 then 1 else 0 end), 0) AS browser_safari,
					COALESCE(SUM(case when browser = 4 then 1 else 0 end), 0) AS browser_firefox,
						
					COALESCE(SUM(case when platform = 1 then 1 else 0 end), 0) AS platform_desktop,
					COALESCE(SUM(case when platform = 2 then 1 else 0 end), 0) AS platform_mobile,
					COUNT(*) AS total_views,
					CURRENT_TIMESTAMP as updated_at
					
					from unique_views WHERE days_since_2020 >= $1 AND days_since_2020 <= $2 group by short_path)
			ON CONFLICT (short_path, report_type) DO UPDATE SET
				browser_chrome  = EXCLUDED.browser_chrome,
				browser_ie      = EXCLUDED.browser_ie,
				browser_safari  = EXCLUDED.browser_safari,
				browser_firefox = EXCLUDED.browser_firefox,

				platform_desktop = EXCLUDED.platform_desktop,
				platform_mobile  = EXCLUDED.platform_mobile,
				
				total_views = EXCLUDED.total_views;`,
		from, to, int(reportType),
	)
	return err
}

func (repo PostgresUniqueViewRepository) GetReport(shortPath string, reportType repositories.ReportType) (*repositories.Report, error) {
	var report repositories.Report
	var err = repo.con.Get(
		&report,
		`SELECT browser_chrome, browser_ie, browser_safari, browser_firefox, 
						platform_desktop, platform_mobile, total_views
						 FROM unique_views_report WHERE short_path=$1 AND report_type=$2;`,
		shortPath, int(reportType))

	if errors.Is(err, sql.ErrNoRows) {
		return &report, nil
	}

	if err != nil {
		return nil, err
	}

	return &report, nil
}
