package log

import (
	_types "github.com/504dev/logr-go-client/types"
	"github.com/504dev/logr/dbs/clickhouse/queue"
	"time"
)

type LogRepo struct {
	queue *queue.Queue
}

func (repo *LogRepo) RunQueue() *LogRepo {
	repo.queue.Run()
	return repo
}

func (repo *LogRepo) StopQueue() error {
	return repo.queue.Stop()
}

func (repo *LogRepo) Store(log *_types.Log) error {
	day := time.Unix(0, log.Timestamp).UTC().Format("2006-01-02")
	values := []interface{}{day, log.Timestamp, log.DashId, log.Hostname, log.Logname, log.Level, log.Message, log.Pid, log.Version}
	repo.queue.Push(values)
	return nil
}
