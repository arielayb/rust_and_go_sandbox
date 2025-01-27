package app

import (
	"sync"

	"github.com/gorilla/websocket"
)

type UserInfo struct {
	UserCn string
}

type SafeStore struct {
	Clients      map[*UserInfo]*websocket.Conn
	mu           sync.Mutex
	BroadcastMsg chan string
}

func NewStore() *SafeStore {
	return &SafeStore{
		Clients:      make(map[*UserInfo]*websocket.Conn),
		mu:           sync.Mutex{},
		BroadcastMsg: make(chan string),
	}

}

func (ss *SafeStore) Set(user *UserInfo, client *websocket.Conn) {
	ss.mu.Lock()
	ss.Clients[user] = client
	ss.mu.Unlock()
}

func (ss *SafeStore) Remove(user *UserInfo, websock *websocket.Conn) {
	ss.mu.Lock()
	websock.Close()
	delete(ss.Clients, user)
	ss.mu.Unlock()
}

func (ss *SafeStore) GetAll() map[*UserInfo]*websocket.Conn {
	tempMap := make(map[*UserInfo]*websocket.Conn)
	ss.mu.Lock()
	for key, value := range ss.Clients {
		tempMap[key] = value
	}
	ss.mu.Unlock()

	return tempMap
}
