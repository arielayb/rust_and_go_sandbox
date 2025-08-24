package app

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type UserInfo struct {
	UserID    string `json:"user_id"`
	Method    string `json:"method"`
	Message   string `json:"msg"`
	Global    bool   `json:"global"`
	WebSocket *websocket.Conn
}

type SafeStore struct {
	Clients    []UserInfo
	mu         *sync.Mutex
	storeCache bool
}

func NewStore() *SafeStore {
	return &SafeStore{
		Clients: []UserInfo{},
		mu:      &sync.Mutex{},
	}
}

func (ss *SafeStore) Set(userId string, ws *websocket.Conn) *UserInfo {
	userInfo := UserInfo{
		UserID:    userId,
		Method:    "USER_INFO",
		Message:   "",
		Global:    false,
		WebSocket: ws,
	}

	ss.mu.Lock()
	ss.Clients = append(ss.Clients, userInfo)
	ss.mu.Unlock()

	return &userInfo
}

func (ss *SafeStore) Get(ws *websocket.Conn) (int, UserInfo) {
	for index, client := range ss.GetAll() {
		if client.WebSocket == ws {
			return index, client
		}
	}

	return -1, UserInfo{}
}

func (ss *SafeStore) Remove(ws *websocket.Conn) {
	ss.mu.Lock()
	for index, client := range ss.Clients {
		if client.WebSocket == ws {
			ss.Clients = append(ss.Clients[:index], ss.Clients[index+1:]...)
		}
	}
	ss.mu.Unlock()
}

func (ss *SafeStore) GetAll() []UserInfo {
	currentClients := []UserInfo{}

	ss.mu.Lock()
	for _, client := range ss.Clients {
		currentClients = append(currentClients, client)
	}
	ss.mu.Unlock()

	return currentClients
}

func (ss *SafeStore) PrintAll() {
	ss.mu.Lock()
	for value := range ss.Clients {
		fmt.Println("the key: ", value)
	}
	ss.mu.Unlock()

	fmt.Println("the number of sessions: ", len(ss.Clients))
}
