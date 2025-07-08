package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type App struct {
	Cache         SafeStore
	ParentContext context.Context
	Post          *UserWebInfo
	ChanMsg       chan string
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	// Time allowed to write a message to the peer.
	//writeWait = 4 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 4 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 4) / 10
)

func (app *App) BroadcastMsg(ctx context.Context, userInfo *UserInfo, ws *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	for {
		userInfo.Message = ""
		// Grab the next message from the broadcast channel
		select {
		case <-ticker.C:
			if err := ws.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
				return
			}
		case <-ctx.Done():
			fmt.Println("Closing write goroutine")
		case userInfo.Message = <-app.ChanMsg:
			if app.Post.Message != "" {
				userInfo.Message = app.Post.Message
				userInfo.UserUUID = app.Post.UserUUID
			}

			js, errjs := json.Marshal(userInfo.Message)
			if errjs != nil {
				log.Fatal("Cannot pack the message as a JSON message!", "ERROR", errjs)
			}

			if userInfo.Message != "" {
				if userInfo.UserUUID == app.Cache.Get(userInfo.UserUUID, ws) {
					// Send the message to all connected clients
					log.Println("Sending the message: ", userInfo.Message)
					err := ws.WriteMessage(websocket.TextMessage, js)
					if err != nil {
						app.Cache.Remove()
					}
				}
			}
		}
	}
}

func (app *App) PostAlert(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "app/json")
	var task UserWebInfo
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Post = &task
	log.Println("post response: ", task)

	app.ChanMsg <- task.Message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// define our WebSocket endpoint
func (app *App) ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
	app.Cache.storeCache = false

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	var userSocketInfo *UserInfo
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var userInfo *UserWebInfo
		err = json.Unmarshal(message, &userInfo)
		if err != nil {
			fmt.Println("Cannot unmarshal the json message!!!")
		}

		switch userInfo.Method {
		case USER_INFO:
			if !app.Cache.storeCache {
				userSocketInfo = app.Cache.Set(userInfo.UserUUID, userInfo.Message, ws)
				// app.Cache.Set(userInfo.UserUUID, userInfo.AlertMsg, ws)
				app.Cache.storeCache = true
				go app.BroadcastMsg(app.ParentContext, userSocketInfo, ws)
				app.Cache.PrintAll()
			}
		}
	}

	if app.Cache.storeCache && userSocketInfo.UserUUID == app.Cache.Get(userSocketInfo.UserUUID, ws) {
		app.Cache.Remove()
	}
}
