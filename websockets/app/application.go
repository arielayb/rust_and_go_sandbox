package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	for {
		// Grab the next message from the broadcast channel
		select {
		case message := <-app.Data.BroadcastMsg:
			fmt.Println("we got a message: ", message)

			js, err := json.Marshal(message)
			if err != nil {
				log.Fatal("Cannot pack the message as a JSON message!", "ERROR", err)
			}

			// Send the message to all connected clients
			tempMap := app.Data.GetAll()
			for key, client := range tempMap {
				err := client.WriteMessage(websocket.TextMessage, js)
				if err != nil {
					app.Data.Remove(key, client)
				}

			}
			fmt.Println("the message to write: ", string(message))
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

	var userData UserInfo
	// todo: implement process to extract current user info from the browser.
	app.Data.Set(&userData, ws)
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	for {
		// read in a message
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println("what we got: ", string(msg))
		userData = UserInfo{UserCn: string(msg)}
		app.Data.PrintAll()
		app.Data.BroadcastMsg <- string(msg)
	}

}
