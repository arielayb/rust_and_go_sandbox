package app

import (
	"sync"
	"websockets/model"

	"github.com/gorilla/websocket"
)

type SafeStore struct {
	Clients map[*websocket.Conn]model.DataPacket
	mu      sync.Mutex
}

func NewStore() *SafeStore {
	return &SafeStore{
		Clients: make(map[*websocket.Conn]model.DataPacket),
		mu:      sync.Mutex{},
	}

}

func (ss *SafeStore) Set(client *websocket.Conn, data model.DataPacket) {
	ss.mu.Lock()
	ss.Clients[client] = data
	ss.mu.Unlock()
}

func (ss *SafeStore) Remove(client *websocket.Conn) {
	ss.mu.Lock()
	delete(ss.Clients, client)
	ss.mu.Unlock()
}
