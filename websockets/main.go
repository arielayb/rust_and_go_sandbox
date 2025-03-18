// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// var clients = make(map[*websocket.Conn]bool) // Connected clients
// var broadcast = make(chan []byte)            // Broadcast channel
// var mutex = &sync.Mutex{}                    // Protect clients map

// func wsHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Error upgrading:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	mutex.Lock()
// 	clients[conn] = true
// 	mutex.Unlock()

// 	for {
// 		// _, message, err := conn.ReadMessage()
// 		conn.ReadMessage()
// 		if err != nil {
// 			mutex.Lock()
// 			delete(clients, conn)
// 			mutex.Unlock()
// 			break
// 		}
// 		// fmt.Println("the message read: ", string(message))
// 		// broadcast <- message
// 	}
// }

// func handleMessages() {
// 	for {
// 		// Grab the next message from the broadcast channel
// 		message := <-broadcast

// 		// Send the message to all connected clients
// 		mutex.Lock()
// 		for client := range clients {
// 			err := client.WriteMessage(websocket.TextMessage, message)
// 			if err != nil {
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 		fmt.Println("the message to write: ", string(message))
// 		mutex.Unlock()
// 	}
// }

// func main() {
// 	http.HandleFunc("/ws", wsHandler)
// 	go handleMessages()
// 	fmt.Println("WebSocket server started on :8080")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		fmt.Println("Error starting server:", err)
// 	}
// }

package main

import (
	"fmt"
	"net/http"
	"websockets/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Task struct {
	Message string `json:"message"`
}

// // We'll need to define an Upgrader
// // this will require a Read and Write buffer size
// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,

// 	// We'll need to check the origin of our connection
// 	// this will allow us to make requests from our React
// 	// development server to here.
// 	// For now, we'll do no checking and just allow any connection
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
// func reader(conn *websocket.Conn) {
// 	for {
// 		// read in a message
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		// print out that message for clarity
// 		fmt.Println(string(p))

// 		if err := conn.WriteMessage(messageType, p); err != nil {
// 			log.Println(err)
// 			return
// 		}

// 	}
// }

func main() {
	fmt.Println("Chatty Batty v0.1")
	// store := app.SafeStore
	app := &app.App{
		Data: *app.NewStore(),
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Post("/alert", app.PostAlert)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome to the home page!"))
		})
	})

	go app.BroadcastMsg()

	// start the websocket
	http.HandleFunc("/ws", app.ServeWs)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
