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

func Find(dashId int, logname string, hostname string, agg string) (types.Counts, error) {
	duration := Logger.Time("/logs:time", time.Millisecond)
	where := `dash_id = ? and logname = ? and timestamp > now() - interval 7 day`
	values := []interface{}{dashId, logname}
	if hostname != "" {
		where += ` and hostname = ?`
		values = append(values, hostname)
	}
	aggmap := map[string]string{
		AggMinute: "toStartOfMinute",
		AggHour:   "toStartOfHour",
		AggDay:    "toStartOfDay",
	}
	aggfunc := aggmap[agg]
	if aggfunc == "" {
		aggfunc = aggmap[AggMinute]
	}
	sql := `
      select
        ` + aggfunc + `(timestamp) as ts,
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
	fmt.Println(sql)
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

func Find2(dashId int, logname string, hostname string, agg string) (types.Counts, error) {
	duration := Logger.Time("/logs:time", time.Millisecond)
	where := `dash_id = ? and logname = ? and timestamp > now() - interval 7 day`
	values := []interface{}{dashId, logname}
	if hostname != "" {
		where += ` and hostname = ?`
		values = append(values, hostname)
	}
	aggmap := map[string]string{
		AggMinute: "toStartOfMinute",
		AggHour:   "toStartOfHour",
		AggDay:    "toStartOfDay",
	}
	aggfunc := aggmap[agg]
	if aggfunc == "" {
		aggfunc = aggmap[AggMinute]
	}
	sql := `
      select
        toUnixTimestamp(` + aggfunc + `(timestamp)) as timestamp,
        hostname,
        keyname,
        sum(inc) as inc,
        max(max) as max,
        min(min) as min,
        sum(avg_sum) as avg_sum,
        sum(avg_num) as avg_num,
        sum(per_tkn) as per_tkn,
        sum(per_ttl) as per_ttl
      from counts
      where ` + where + `
      group by
        timestamp, hostname, keyname
      order by
        timestamp desc, hostname, keyname
    `
	fmt.Println(sql)

	counts := types.Counts{}
	err := clickhouse.Conn().Select(&counts, sql, values...)
	if err != nil {
		return nil, err
	}
	duration()
	Logger.Inc("/logs:cnt", 1)
	return counts, nil
}

func GetDashStats(dashIds []int) ([]*types.DashStatRow, error) {
	conn := clickhouse.Conn()
	sql := `
      SELECT dash_id, hostname, logname, version, count(*) AS cnt, max(toUnixTimestamp(timestamp)) AS updated
      FROM counts WHERE dash_id IN (?)
      GROUP BY dash_id, hostname, logname, version
    `
	stats := types.DashStatRows{}
	err := conn.Select(&stats, sql, dashIds)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
