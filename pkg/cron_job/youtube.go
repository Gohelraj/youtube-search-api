package cron_job

import (
	"github.com/Gohelraj/youtube-search-api/config"
	"github.com/Gohelraj/youtube-search-api/pkg/youtube"
	"log"
)

func (c CronJob) FetchYoutubeVideosAndAddToQueue() {
	_, err := c.CronObj.AddFunc(config.Conf.CronSpecsToFetchVideos, func() {
		log.Println("Fetching youtube videos")
		youtube.SearchVideosFromYoutube(config.Conf.VideoKeyword, c.PgxPool)
		log.Println("Queued youtube videos")
	})
	if err != nil {
		log.Fatalf("error adding cron job: %v\n", err)
	}
}