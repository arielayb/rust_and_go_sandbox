package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var testApp *App

func TestMain(m *testing.M) {
	ctx, stop := context.WithCancel(context.Background())
	testApp = &App{
		Cache:         *NewStore(),
		ParentContext: ctx,
		Post:          []UserWebInfo{},
	}
	exitVal := m.Run()
	os.Exit(exitVal)
	stop()
}

func TestSafeCacheStore(t *testing.T) {
	// Create test server with the echo handler.
	server := httptest.NewServer(http.HandlerFunc(testApp.ServeWs))
	serverPost := httptest.NewServer(http.HandlerFunc(testApp.PostAlert))
	defer server.Close()
	defer serverPost.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	upgradeToWs := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(upgradeToWs, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	testData := UserWebInfo{
		UserID:  "guest1",
		Method:  "USER_INFO",
		Message: "hello1",
		Global:  false,
	}

	jsonMsg, errJson := json.Marshal(&testData)
	if errJson != nil {
		t.Fatalf("Error: %v", errJson.Error())
	}

	if err := (ws.WriteMessage(websocket.TextMessage, jsonMsg)); err != nil {
		t.Fatalf("%v", err)
	}
	time.Sleep(1 * time.Second)
	clients := testApp.Cache.Clients
	fmt.Println("the Queue: ", clients)

	assert.Equal(t, len(clients), 1)
}
