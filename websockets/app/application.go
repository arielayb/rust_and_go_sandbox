package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type App struct {
	Data SafeStore
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
	writeWait = 4 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 6 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func (app *App) PostAlert(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	var task string
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//lastID++
	//tasks[lastID] = task
	app.Data.BroadcastMsg <- task
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)

}

func (app *App) BroadcastMsg() {
	ticker := time.NewTicker(pingPeriod)
	for {
		message := ""
		tempMap := app.Data.GetAll()
		// Grab the next message from the broadcast channel
		select {
		case message = <-app.Data.BroadcastMsg:
			fmt.Println("we got a message: ", message)
		case <-ticker.C:
			// client.SetWriteDeadline(time.Now().Add(writeWait))
			for _, client := range tempMap {
				if err := client.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
					return
				}
			}
		}

		js, errjs := json.Marshal(message)
		if errjs != nil {
			log.Fatal("Cannot pack the message as a JSON message!", "ERROR", errjs)
		}

		// Send the message to all connected clients
		for key, client := range tempMap {
			if len(message) != 0 {
				err := client.WriteMessage(websocket.TextMessage, js)
				if err != nil {
					app.Data.Remove(key, client)
				}
				fmt.Println("the message to write: ", string(message))
			}
		}

	}
}

// define our WebSocket endpoint
func (app *App) ServeWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	// stdoutDone := make(chan struct{})
	// go app.Ping(ws, stdoutDone)

	var userData UserInfo
	app.Data.Set(&userData, ws)
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		} else {
			fmt.Println("what we got: ", string(message))
			userData.UserCn = string(message)
		}
		// message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		// app.Data.BroadcastMsg <- string(message)

		// print out that message for clarity
		app.Data.PrintAll()
	}

}
