CREATE TABLE IF NOT EXISTS counts
(
    `day` Date,
    `timestamp` DateTime,
    `dash_id` Int32,
    `hostname` String,
    `logname` String,
    `keyname` String,
    `version` String,
    `inc` Nullable(Float32),
    `max` Nullable(Float32),
    `min` Nullable(Float32),
    `avg_sum` Nullable(Float32),
    `avg_num` Nullable(Int32),
    `per_tkn` Nullable(Float32),
    `per_ttl` Nullable(Float32),
    `time_dur` Nullable(Int32)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(day)
ORDER BY (dash_id, logname, day, timestamp)
TTL day + INTERVAL 1 YEAR