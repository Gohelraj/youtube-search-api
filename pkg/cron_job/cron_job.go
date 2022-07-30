package cron_job

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	CronObj *cron.Cron
	PgxPool *pgxpool.Pool
}

func NewCronJobObject(c *cron.Cron, pgxPool *pgxpool.Pool) CronJob {
	return CronJob{
		CronObj: c,
		PgxPool: pgxPool,
	}
}

func Init(pgxPool *pgxpool.Pool) {
	c := cron.New(
		cron.WithParser(
			cron.NewParser(
				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)))
	cronObj := NewCronJobObject(c, pgxPool)
	cronObj.FetchYoutubeVideosAndAddToQueue()
	c.Start()
}
