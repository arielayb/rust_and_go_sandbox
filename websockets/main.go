package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
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
	ctx, stop := context.WithCancel(context.Background())

	go func() {
		log.Println(http.ListenAndServe("localhost:8082", nil))
	}()

	app := &app.App{
		Cache:         *app.NewStore(),
		ParentContext: ctx,
		Post:          []app.UserWebInfo{},
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Post("/alert", app.PostAlert)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome to the home page!"))
		})
	})

	// start the websocket
	r.HandleFunc("/ws", app.ServeWs)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}

	stop()
}
