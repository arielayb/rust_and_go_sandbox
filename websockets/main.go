package main

import (
	"context"
	"fmt"
	"net/http"
	"websockets/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Task struct {
	Message string `json:"message"`
}

func main() {
	fmt.Println("Chatty Batty v0.1")
	// parent context
	ctx := context.Background()

	// store := app.SafeStore
	app := &app.App{
		Data:          *app.NewStore(),
		ParentContext: ctx,
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
