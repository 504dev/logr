package count

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/dbs/clickhouse"
	"github.com/504dev/logr/dbs/clickhouse/batcher"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	AGG_MINUTE   = "m"
	AGG_5MINUTES = "5m"
	AGG_HOUR     = "h"
	AGG_DAY      = "d"
)

type CountRepo struct {
	conn    *sqlx.DB
	batcher *batcher.Batcher[*_types.Count]
}

func NewCountRepo() *CountRepo {
	return &CountRepo{
		conn: clickhouse.Conn(),
	}
}

func (repo *CountRepo) Insert(count *_types.Count) error {
	return repo.BatchInsert([]*_types.Count{count})
}

func (repo *CountRepo) BatchInsert(batch []*_types.Count) error {
	const sql = `
		INSERT INTO counts (
		    day, timestamp, dash_id, hostname, logname, keyname, version,
		    inc, max, min, avg_sum, avg_num, per_tkn, per_ttl, time_dur
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	tx, err := repo.conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range batch {
		values := (*types.CountVector)(v).AsVector()
		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repo *CountRepo) Find(filter types.Filter, agg string) (types.Counts, error) {
	where := `dash_id = ? and logname = ?`
	values := []interface{}{filter.DashId, filter.Logname}
	if filter.Hostname != "" {
		where += ` and hostname = ?`
		values = append(values, filter.Hostname)
	}
	if filter.Keyname != "" {
		where += ` and keyname = ?`
		values = append(values, filter.Keyname)
	}
	if filter.Version != "" {
		where += ` and version = ?`
		values = append(values, filter.Version)
	}
	if filter.Timestamp != nil {
		if filter.Timestamp[0] != 0 {
			where += " AND timestamp > ?"
			values = append(values, filter.Timestamp[0])
		}
		if filter.Timestamp[1] != 0 {
			where += " AND timestamp <= ?"
			values = append(values, filter.Timestamp[1])
		}
	}
	aggmap := map[string][]string{
		AGG_MINUTE:   {"toStartOfMinute", "6 hour"},
		AGG_5MINUTES: {"toStartOfFiveMinute", "1 day"},
		AGG_HOUR:     {"toStartOfHour", "14 day"},
		AGG_DAY:      {"toStartOfDay", "366 day"},
	}
	aggvalues, ok := aggmap[agg]
	if !ok {
		aggvalues = aggmap[AGG_MINUTE]
	}
	where += " and timestamp > now() - interval " + aggvalues[1]
	sql := `
      select
        ` + aggvalues[0] + `(timestamp) as ts,
        hostname,
        keyname,
        sum(inc),
        max(max),
        min(min),
        sum(avg_sum),
        sum(avg_num),
        sum(per_tkn),
        sum(per_ttl)
      from counts
      where ` + where + `
      group by
        ts, hostname, keyname
      order by
        ts desc, hostname, keyname
    `
	Logger.Debug("%v %v", sql, values)

	rows, err := repo.conn.Query(sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := types.Counts{}

	for rows.Next() {
		var timestamp time.Time
		var hostname, keyname string
		var inc, max, min, avgSum, perTotal, perTaken *float64
		var avgNum *int
		err = rows.Scan(&timestamp, &hostname, &keyname, &inc, &max, &min, &avgSum, &avgNum, &perTaken, &perTotal)
		if err != nil {
			return nil, err
		}
		metrics := _types.Metrics{}
		if inc != nil {
			metrics.Inc = &_types.Inc{Val: *inc}
		}
		if max != nil {
			metrics.Max = &_types.Max{Val: *max}
		}
		if min != nil {
			metrics.Min = &_types.Min{Val: *min}
		}
		if avgNum != nil && avgSum != nil {
			metrics.Avg = &_types.Avg{Sum: *avgSum, Num: *avgNum}
		}
		if perTaken != nil && perTotal != nil {
			metrics.Per = &_types.Per{Total: *perTotal, Taken: *perTaken}
		}
		counts = append(counts, &_types.Count{
			Timestamp: timestamp.Unix(),
			Hostname:  hostname,
			Keyname:   keyname,
			Metrics:   metrics,
		})
	}

	err = rows.Err()

	return counts, err
}

func (repo *CountRepo) StatsByLogname(dashId int, logname string) ([]*types.DashStatRow, error) {
	sql := `
      SELECT hostname, version, count(*) AS cnt, max(toUnixTimestamp(timestamp)) AS updated
      FROM counts
      WHERE dash_id = ? AND logname = ? AND day > toDate(now() - interval 7 day)
      GROUP BY hostname, version
    `
	Logger.Debug("%v %v", sql, []interface{}{dashId, logname})
	stats := types.DashStatRows{}
	err := repo.conn.Select(&stats, sql, dashId, logname)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (repo *CountRepo) StatsByDashboard(dashId int) ([]*types.DashStatRow, error) {
	sql := `
      SELECT logname, count(*) AS cnt
      FROM counts
      WHERE
        dash_id = ? AND
        day >= toDate(now() - interval 1 day) AND
        timestamp > now() - interval 1 hour
      GROUP BY logname
    `
	stats := types.DashStatRows{}
	err := repo.conn.Select(&stats, sql, dashId)
	if err != nil {
		return nil, err
	}
	Logger.Debug("%v %v", sql, dashId)
	return stats, nil
}
