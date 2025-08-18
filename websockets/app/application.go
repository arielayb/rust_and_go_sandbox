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
	Post          []UserWebInfo
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
	pingPeriod = (pongWait * 9) / 10
)

func (app *App) BroadcastMsg(ctx context.Context, userInfo *UserInfo) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := userInfo.WebSocket.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
				return
			}
		case <-ctx.Done():
			fmt.Println("Closing write goroutine")
		}

		// shallow copy the Post list
		//tempPost := app.Post
		if len(app.Post) > 0 {
			for index := range app.Post {
				if userInfo.USERID == app.Cache.Get(app.Post[index].UserID, userInfo.WebSocket) {
					if app.Post[index].Message != "" && !app.Post[index].Global {
						// Send the message to all connected clients
						log.Println("Sending the message: ", app.Post[index].Message)
						err := userInfo.WebSocket.WriteMessage(websocket.TextMessage, []byte(app.Post[index].Message))
						if err != nil {
							break
						} else {
							time.Sleep(1 * time.Second)
							// clear the index of the user information
							app.Post[index].Message = ""
						}
					} else if app.Post[index].Global {
						for _, client := range app.Cache.GetAll().In {
							// Send the message to all connected clients
							log.Println("Sending the message: ", app.Post[index].Message)
							err := client.WebSocket.WriteMessage(websocket.TextMessage, []byte(app.Post[index].Message))
							if err != nil {
								break
							}
						}
						time.Sleep(1 * time.Second)
						app.Post[index].Message = ""
					}
				}

				//clear the buffer list
				for _, msg := range app.Post {
					if msg.Message == "" {
						log.Println("removing message queue: ", app.Post)
						app.Post = app.Post[1:]
					}
				}
			}
		}
		//app.Post = tempPost
		log.Println("the message queue: ", app.Post)
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
	defer r.Body.Close()

	log.Println("the post message content: ", task)

	// create a queue of user messages
	app.Post = append(app.Post, task)

	log.Println("post response: ", app.Post)
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

		var userSocketInfo *UserInfo
		switch userInfo.Method {
		case USER_INFO:
			if !app.Cache.storeCache {
				userSocketInfo = app.Cache.Set(userInfo.UserID, ws)
				app.Cache.storeCache = true
				go app.BroadcastMsg(app.ParentContext, userSocketInfo)
			}
			app.Cache.PrintAll()
		}
	}

	if app.Cache.storeCache {
		app.Cache.Remove()
	}
}
