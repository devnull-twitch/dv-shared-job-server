package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/devnull-twitch/sharedjob-server"
	"github.com/sirupsen/logrus"
)

// Helper script to finish some tasks to get to a certain point.

func main() {
	httpClient := &http.Client{}

	for _, jobId := range os.Args[1:] {
		logrus.WithField("job_id", jobId).Info("progressing job")
		makeProgressReq(httpClient, jobId, "reserve")
		time.Sleep(time.Millisecond * 100)
		makeProgressReq(httpClient, jobId, "take")
		time.Sleep(time.Millisecond * 100)
		makeProgressReq(httpClient, jobId, "finish")
		time.Sleep(time.Millisecond * 100)
	}

	logrus.Info("Done")
}

func makeProgressReq(client *http.Client, jobId, actionType string) {
	fakeUser := &sharedjob.UserIDPayload{
		Name: "Insomnia",
	}
	jsonBytes, _ := json.Marshal(fakeUser)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:8083/job/%s/%s", jobId, actionType),
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}
