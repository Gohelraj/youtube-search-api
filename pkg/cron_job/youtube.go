package cron_job

import (
	"github.com/Gohelraj/youtube-search-api/config"
	"github.com/Gohelraj/youtube-search-api/pkg/youtube"
	"log"
)

// FetchYoutubeVideosAndAddToQueue fetches youtube videos and adds to queue
func (c CronJob) FetchYoutubeVideosAndAddToQueue() {
	_, err := c.CronObj.AddFunc(config.Conf.CronSpecsToFetchVideos, func() {
		log.Println("Fetching youtube videos")
		youtube.SearchVideosFromYoutubeAndAddToQueue(config.Conf.VideoKeyword, c.PgxPool)
	})
	if err != nil {
		log.Fatalf("error adding cron job: %v\n", err)
	}
}