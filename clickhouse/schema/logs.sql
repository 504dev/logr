CREATE TABLE IF NOT EXISTS logs
(
    `day` Date,
    `timestamp` UInt64,
    `dash_id` Int32,
    `hostname` String,
    `logname` String,
    `level` String,
    `message` String,
    `pid` Int32,
    `version` String
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(day)
ORDER BY (dash_id, logname, day, timestamp)
TTL day + INTERVAL 1 DAY