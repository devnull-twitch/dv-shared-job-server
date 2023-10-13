package sharedjob

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{}
var players = []*Player{}

func getPlayer(wsConn *websocket.Conn) *Player {
	for _, playerObj := range players {
		if playerObj.wsConn == wsConn {
			return playerObj
		}
	}
	return nil
}

func GetPlayers() []*Player {
	return slices.Clone(players)
}

type (
	clientWelcomeMessage struct {
		Conn     *websocket.Conn `json:"-"`
		Username string          `json:"username"`
	}
	clientMessage struct {
		StationID StationID       `json:"station_id"`
		Unsub     bool            `json:"unsub"`
		Conn      *websocket.Conn `json:"-"`
	}
	Player struct {
		wsConn         *websocket.Conn
		subbedStations []StationID
		Username       string
	}
	ProgressMessage struct {
		StationID StationID `json:"station_id"`
	}
)

func (p *Player) IsSubbedToStation(s StationID) bool {
	for _, subbedStation := range p.subbedStations {
		if subbedStation == s {
			return true
		}
	}

	return false
}

func (p *Player) GetSubbedStations() []StationID {
	return slices.Clone(p.subbedStations)
}

func HandleWebsocket(
	w http.ResponseWriter,
	r *http.Request,
	msgChan chan<- clientMessage,
) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(fmt.Errorf("could not upgrade connection: %w", err))
	}
	defer conn.Close()

	welcome := clientWelcomeMessage{}
	if err := conn.ReadJSON(&welcome); err != nil {
		logrus.WithError(err).Error("could not read welcome message")
		return
	}

	if welcome.Username == "" {
		logrus.WithError(err).Error("empty username")
		return
	}

	players = append(players, &Player{
		wsConn:         conn,
		Username:       welcome.Username,
		subbedStations: make([]StationID, 0),
	})

	for {
		v := clientMessage{}
		err := conn.ReadJSON(&v)
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				panic(fmt.Errorf("could not read message: %w", err))
			}

			cleanUpClosedConnection(conn)
			return
		}
		v.Conn = conn

		msgChan <- v
	}
}

func cleanUpClosedConnection(conn *websocket.Conn) {
	for index, playerObj := range players {
		if playerObj.wsConn == conn {
			players = slices.Delete(players, index, index+1)
			return
		}
	}
}

func StartWSProcessor() (chan<- clientMessage, chan<- ProgressMessage) {
	clientCh := make(chan clientMessage)
	progressCh := make(chan ProgressMessage)

	go func() {
		defer func() {
			logrus.Info("Worker shutdown")
		}()

		for {
			select {
			case msg := <-clientCh:
				playerObj := getPlayer(msg.Conn)
				if playerObj == nil {
					logrus.Error("websocket message from unknown player")
					return
				}

				logrus.WithFields(logrus.Fields{
					"station_id":  msg.StationID,
					"unsubscribe": msg.Unsub,
				}).Info("received client message")

				if msg.Unsub {
					for i, stationID := range playerObj.subbedStations {
						if stationID == msg.StationID {
							playerObj.subbedStations = slices.Delete(playerObj.subbedStations, i, i+1)
							break
						}
					}
				} else {
					playerObj.subbedStations = append(playerObj.subbedStations, msg.StationID)
				}
			case msg := <-progressCh:
				for _, playerObj := range players {
					if playerObj.IsSubbedToStation(msg.StationID) {
						if err := playerObj.wsConn.WriteJSON(msg); err != nil {
							logrus.WithError(err).Error("could not send progress message")
							continue
						}

						logrus.WithFields(logrus.Fields{
							"station_id": msg.StationID,
							"username":   playerObj.Username,
						}).Info("sent progress message")
					}
				}
			}
		}
	}()

	return clientCh, progressCh
}
