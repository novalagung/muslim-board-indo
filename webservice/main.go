package main

import (
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"muslimboard-api.novalagung.com/pkg/logger"
	router "muslimboard-api.novalagung.com/router"
)

func main() {
	logger.Init()
	defer sentry.Flush(2 * time.Second)

	http.HandleFunc("/muslimboard-api", router.MuslimboardApi)

	port := "0.0.0.0:" + os.Getenv("PORT")
	logger.Log.Infoln("listening to", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		logger.Log.Fatal(err)
	}
}
