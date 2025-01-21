package app

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SafeStore struct {
	Clients map[string][]*websocket.Conn
	mu      sync.Mutex
}

func NewStore() *SafeStore {
	return &SafeStore{
		Clients: make(map[string][]*websocket.Conn),
		mu:      sync.Mutex{},
	}

}

func (ss *SafeStore) Set(user string, client *websocket.Conn) {
	ss.mu.Lock()
	ss.Clients[user] = append(ss.Clients[user], client)
	ss.mu.Unlock()
}

func (ss *SafeStore) Remove(user string) {
	ss.mu.Lock()
	delete(ss.Clients, user)
	ss.mu.Unlock()
}

func (ss *SafeStore) GetAll() map[string][]*websocket.Conn {
	ss.mu.Lock()
	tempMap := map[string][]*websocket.Conn{}
	for key, value := range ss.Clients {
		tempMap[key] = value
	}
	ss.mu.Unlock()

	return tempMap
}
