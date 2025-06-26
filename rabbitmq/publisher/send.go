package main

import (
  "context"
  "log"
  "time"
  "fmt"
  "net/http"
  "github.com/gorilla/websocket"
  amqp "github.com/rabbitmq/amqp091-go"
)

var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
  CheckOrigin:     func(r *http.Request) bool { return true },
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
  // upgrade this connection to a WebSocket
  // connection
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
  }
  
  log.Println("Client Connected")
  
  err = ws.WriteMessage(1, []byte("Hi Client!"))
  if err != nil {
    log.Println(err)
  }
  // listen indefinitely for new messages coming
  // through on our WebSocket connection
  reader(ws)
}

func reader(conn *websocket.Conn) {
  for {
    // read in a message
    messageType, p, err := conn.ReadMessage()
    if err != nil {
      log.Println(err)
      return
    }
  
    // print out that message for clarity
    log.Println(string(p))
    if err := conn.WriteMessage(messageType, p); err != nil {
      log.Println(err)
      return
    }
  }
}

func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Home Page")
}

func setupRoutes() {
  http.HandleFunc("/", homePage)
  http.HandleFunc("/ws", wsEndpoint)
}

func main() {
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  failOnError(err, "Failed to connect to RabbitMQ")
  defer conn.Close()

  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  defer ch.Close()

  q, err := ch.QueueDeclare(
    "hello", // name
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  failOnError(err, "Failed to declare a queue")

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  body := "Hello World!"
  err = ch.PublishWithContext(ctx,
  "",     // exchange
  q.Name, // routing key
  false,  // mandatory
  false,  // immediate
  
  amqp.Publishing {
    ContentType: "text/plain",
    Body:        []byte(body),
  })
  
  failOnError(err, "Failed to publish a message")
  log.Printf(" [x] Sent %s\n", body)

  setupRoutes()
  log.Fatal(http.ListenAndServe(":8080", nil))
}

