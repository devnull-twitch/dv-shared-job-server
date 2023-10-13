package main

import (
	"net"
	"net/http"

	"github.com/devnull-twitch/sharedjob-server"
	"github.com/devnull-twitch/sharedjob-server/ui"
	"github.com/gin-gonic/gin"
)

func main() {
	sharedjob.Setup()

	clientCh, processorCh := sharedjob.StartWSProcessor()

	r := gin.Default()
	r.GET("/station/:station", func(c *gin.Context) {
		username := c.Query("username")

		stationCode := sharedjob.StationID(c.Param("station"))
		var jobs []sharedjob.Job
		if username != "" {
			jobs = sharedjob.GetAllStationJobsForUsername(stationCode, username)
		} else {
			jobs = sharedjob.GetAllStationJobs(stationCode)
		}
		c.JSON(200, jobs)
	})
	r.POST("/job/:job_id/reserve", func(c *gin.Context) {
		userPayload := &sharedjob.UserIDPayload{}
		if err := c.BindJSON(userPayload); err != nil {
			// make better error handling
			panic("no valid user payload")
		}

		jobID := c.Param("job_id")
		if !sharedjob.ReserveJob(userPayload.Name, jobID) {
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusOK)
	})
	r.POST("/job/:job_id/take", func(c *gin.Context) {
		userPayload := &sharedjob.UserIDPayload{}
		if err := c.BindJSON(userPayload); err != nil {
			// make better error handling
			panic("no valid user payload")
		}

		jobID := c.Param("job_id")
		if ok, _, _, _ := sharedjob.TakeJob(
			userPayload.Name,
			jobID,
			processorCh,
		); !ok {
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusOK)
	})
	r.GET("/fakeprogress/:station", func(c *gin.Context) {
		stationCode := sharedjob.StationID(c.Param("station"))
		processorCh <- sharedjob.ProgressMessage{StationID: stationCode}
	})
	r.POST("/job/:job_id/finish", func(c *gin.Context) {
		userPayload := &sharedjob.UserIDPayload{}
		if err := c.BindJSON(userPayload); err != nil {
			// make better error handling
			panic("no valid user payload")
		}

		jobID := c.Param("job_id")
		if ok, _, _, _, _ := sharedjob.FinishJob(userPayload.Name, jobID, processorCh); !ok {
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusOK)
	})

	r.Any("/ws", func(c *gin.Context) {
		sharedjob.HandleWebsocket(c.Writer, c.Request, clientCh)
	})

	ui.AddUIHandlers(r, processorCh)

	tcpListener, err := net.Listen("tcp", ":8083")
	if err != nil {
		panic(err)
	}
	err = r.RunListener(tcpListener)
	if err != nil {
		panic(err)
	}
}
