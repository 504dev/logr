CREATE TABLE IF NOT EXISTS counts
(
    `day` Date,
    `timestamp` DateTime,
    `dash_id` Int32,
    `hostname` String,
    `logname` String,
    `keyname` String,
    `inc` Nullable(Float32),
    `max` Nullable(Float32),
    `min` Nullable(Float32),
    `avg_sum` Float32,
    `avg_num` Int32,
    `per_tkn` Float32,
    `per_ttl` Float32
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(day)
ORDER BY (dash_id, logname, day, timestamp)