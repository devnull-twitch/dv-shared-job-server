package ui

import (
	"net/http"
	"slices"

	"github.com/devnull-twitch/sharedjob-server"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//go:generate jade -pkg=ui -writer -fmt -basedir ../jade pages parts

func AddUIHandlers(r *gin.Engine, processorCh chan<- sharedjob.ProgressMessage) {
	uiLog := logrus.WithField("module", "ui")

	ui := r.Group("/ui")
	{
		ui.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "jobs")
		})
		ui.GET("/jobs", func(c *gin.Context) {
			JobsView("Jobs", sharedjob.AllStations, c.Writer)
			c.Status(http.StatusOK)
		})
		ui.GET("/stations", func(c *gin.Context) {
			StationsView("Stations", sharedjob.AllStations, c.Writer)
			c.Status(http.StatusOK)
		})
		ui.GET("/jobs/:jobid/take", func(c *gin.Context) {
			jobID := c.Param("jobid")
			JobPartTake(jobID, c.Writer)
			c.Status(http.StatusOK)
		})
		ui.POST("/jobs/:jobid/take", func(c *gin.Context) {
			jobID := c.Param("jobid")
			assigneUser := c.PostForm("user")

			if assigneUser == "" {
				uiLog.Warn("No username provided")
				c.Status(http.StatusBadRequest)
				return
			}

			var ok bool
			if ok = sharedjob.ReserveJob(assigneUser, jobID); !ok {
				c.Status(http.StatusBadRequest)
				return
			}

			ok, unspawnedJobs, spawnedJobs, changedJobs := sharedjob.TakeJob(
				assigneUser,
				jobID,
				processorCh,
			)
			if !ok {
				c.Status(http.StatusBadRequest)
				return
			}

			uiChangedJobs := slices.Clone(changedJobs)
			uiChangedJobs = append(uiChangedJobs, unspawnedJobs...)
			uiChangedJobs = append(uiChangedJobs, spawnedJobs...)

			c.Status(http.StatusOK)
			JobPartUpdates(uiChangedJobs, nil, "", c.Writer)
		})
		ui.GET("/jobs/:jobid/finish", func(c *gin.Context) {
			jobID := c.Param("jobid")
			JobPartFinish(jobID, c.Writer)
			c.Status(http.StatusOK)
		})
		ui.POST("/jobs/:jobid/finish", func(c *gin.Context) {
			jobID := c.Param("jobid")
			assigneUser := c.PostForm("user")

			if assigneUser == "" {
				uiLog.Warn("No username provided")
				c.Status(http.StatusBadRequest)
				return
			}

			ok, unspawnedJobs, spawnedJobs, changedJobs, newJobs := sharedjob.FinishJob(
				assigneUser,
				jobID,
				processorCh,
			)
			if !ok {
				c.Status(http.StatusBadRequest)
				return
			}

			uiChangedJobs := slices.Clone(changedJobs)
			uiChangedJobs = append(uiChangedJobs, unspawnedJobs...)
			uiChangedJobs = append(uiChangedJobs, spawnedJobs...)

			c.Status(http.StatusOK)
			JobPartUpdates(uiChangedJobs, newJobs, jobID, c.Writer)
		})
		ui.GET("/connections", func(c *gin.Context) {
			ConnectedPlayersView("Players", sharedjob.GetPlayers(), c.Writer)
			c.Status(http.StatusOK)
		})
	}
}
