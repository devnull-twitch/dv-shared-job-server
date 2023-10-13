package ui

import (
	"fmt"

	"github.com/devnull-twitch/sharedjob-server"
)

func jobIdAttr(jobID string) string {
	return fmt.Sprintf("job-%s", jobID)
}

func jobTakeURL(jobID string) string {
	return fmt.Sprintf("/ui/jobs/%s/take", jobID)
}

func jobFinishURL(jobID string) string {
	return fmt.Sprintf("/ui/jobs/%s/finish", jobID)
}

func countSpawnedJobs(station *sharedjob.LogicStation) int {
	count := 0
	for _, job := range station.JobQueue {
		if job.IsSpawned() {
			count++
		}
	}
	return count
}
