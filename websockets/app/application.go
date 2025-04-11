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

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (application *App) PostAlert(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	var task string
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//lastID++
	//tasks[lastID] = task
	application.Cache.BroadcastMsg <- task
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)

}

func (application *App) BroadcastMsg(ctx context.Context, userInfo *UserInfo) {
	ticker := time.NewTicker(pingPeriod)
	for {
		// Grab the next message from the broadcast channel
		select {
		case <-ticker.C:
			if err := userInfo.Conn.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
				return
			}
		case <-ctx.Done():
			fmt.Println("Closing write goroutine")
		}

		js, errjs := json.Marshal(userInfo.UserMsg)
		if errjs != nil {
			log.Fatal("Cannot pack the message as a JSON message!", "ERROR", errjs)
		}

		if userInfo.UserMsg != "" {
			// Send the message to all connected clients
			err := userInfo.Conn.WriteMessage(websocket.TextMessage, js)
			if err != nil {
				application.Cache.Remove()
			}
		}
	}
}

// define our WebSocket endpoint
func (application *App) ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

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

		userInfo := UserInfo{
			UserMsg: string(message),
			Conn:    ws,
		}

		// select {
		// case app.SafeStore.BroadcastMsg <- message:

		// }

		fmt.Println("what we got: ", string(message))
		application.Cache.Set(string(message), ws)

		go application.BroadcastMsg(application.ParentContext, &userInfo)

		application.Cache.PrintAll()
	}
}
