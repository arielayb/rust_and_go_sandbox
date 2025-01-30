package app

import (
	"bytes"
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
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

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
	message := ""
	for {
		tempMap := app.Data.GetAll()
		// Grab the next message from the broadcast channel
		for key, client := range tempMap {
			select {
			case message = <-app.Data.BroadcastMsg:
				fmt.Println("we got a message: ", message)

				js, errjs := json.Marshal(message)
				if errjs != nil {
					log.Fatal("Cannot pack the message as a JSON message!", "ERROR", errjs)
				}

				// Send the message to all connected clients
				// tempMap := app.Data.GetAll()
				// for key, client := range tempMap {
				client.SetWriteDeadline(time.Now().Add(writeWait))
				err := client.WriteMessage(websocket.TextMessage, js)
				if err != nil {
					app.Data.Remove(key, client)
				}

				//}
				fmt.Println("the message to write: ", string(message))

			case <-ticker.C:
				msg := "keepalive"
				client.SetWriteDeadline(time.Now().Add(writeWait))
				if err := client.WriteMessage(websocket.PingMessage, []byte(msg)); err != nil {
					return
				}
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

	stdoutDone := make(chan struct{})
	go app.Ping(ws, stdoutDone)

	var userData UserInfo
	// todo: implement process to extract current user info from the browser.
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
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
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		app.Data.BroadcastMsg <- string(message)

		// print out that message for clarity
		fmt.Println("what we got: ", string(message))
		userData = UserInfo{UserCn: string(message)}
		app.Data.PrintAll()
	}

}

func (app *App) Ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}
