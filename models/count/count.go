package count

import (
	"fmt"
	"github.com/504dev/logr/clickhouse"
	. "github.com/504dev/logr/logger"
	"github.com/504dev/logr/types"
	"time"
)

const (
	AggMinute = "m"
	AggHour   = "h"
	AggDay    = "d"
)

func Find(filter types.Filter, agg string) (types.Counts, error) {
	duration := Logger.Time("/logs:time", time.Millisecond)
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
	if filter.Timestamp[0] != 0 {
		where += " AND timestamp > ?"
		values = append(values, filter.Timestamp[0])
	}
	if filter.Timestamp[1] != 0 {
		where += " AND timestamp <= ?"
		values = append(values, filter.Timestamp[1])
	}
	aggmap := map[string][]string{
		AggMinute: {"toStartOfMinute", "1 day"},
		AggHour:   {"toStartOfHour", "30 day"},
		AggDay:    {"toStartOfDay", "360 day"},
	}
	aggvalues, ok := aggmap[agg]
	if !ok {
		aggvalues = aggmap[AggMinute]
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
	fmt.Println(sql, values)

	rows, err := clickhouse.Conn().Query(sql, values...)
	if err != nil {
		return nil, err
	}

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
		metrics := types.Metrics{}
		if inc != nil {
			metrics.Inc = &types.Inc{Val: *inc}
		}
		if max != nil {
			metrics.Max = &types.Max{Val: *max}
		}
		if min != nil {
			metrics.Min = &types.Min{Val: *min}
		}
		if avgNum != nil && avgSum != nil {
			metrics.Avg = &types.Avg{Sum: *avgSum, Num: *avgNum}
		}
		if perTaken != nil && perTotal != nil {
			metrics.Per = &types.Per{Total: *perTotal, Taken: *perTaken}
		}
		counts = append(counts, &types.Count{
			Timestamp: timestamp.Unix(),
			Hostname:  hostname,
			Keyname:   keyname,
			Metrics:   metrics,
		})
	}
	duration()
	Logger.Inc("/logs:cnt", 1)
	return counts, nil
}

func GetDashStats(dashId int) ([]*types.DashStatRow, error) {
	conn := clickhouse.Conn()
	sql := `
      SELECT hostname, logname, version, count(*) AS cnt, max(toUnixTimestamp(timestamp)) AS updated
      FROM counts
      WHERE dash_id = ?
      GROUP BY hostname, logname, version
    `
	stats := types.DashStatRows{}
	err := conn.Select(&stats, sql, dashId)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
