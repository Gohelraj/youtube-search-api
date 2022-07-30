package main

import (
	"fmt"
	"github.com/Gohelraj/youtube-search-api/api/route"
	"github.com/Gohelraj/youtube-search-api/config"
	"github.com/Gohelraj/youtube-search-api/db"
	"github.com/Gohelraj/youtube-search-api/pkg/cron_job"
	"github.com/Gohelraj/youtube-search-api/pkg/youtube"
	"log"
	"net/http"
	"time"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v\n", err)
	}

	pgxPool := db.Connect(config.Conf.DatabaseURL)
	// closes db connection after the server is shut down
	defer pgxPool.Close()

	port := fmt.Sprintf(":%s", config.Conf.Port)
	// Start the server
	srv := &http.Server{
		Addr:    port,
		Handler: route.InitializeRouter(pgxPool),
		// IdleTimeout is the maximum amount of time to wait for the
		// next request when keep-alives are enabled.
		IdleTimeout: 10 * time.Minute,
	}

	// Start amqp consumer to process youtube videos from queue
	go youtube.ProcessYoutubeVideosFromQueue(pgxPool)
	// start event scheduler on app start
	go cron_job.Init(pgxPool)

	log.Printf("Server listening on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}
