package app

import (
	"fmt"
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

func (ss *SafeStore) PrintAll() {
	ss.mu.Lock()
	for key, _ := range ss.Clients {
		fmt.Println("the key: ", key)
	}
	ss.mu.Unlock()

	fmt.Println("the number of sessions: ", len(ss.Clients))
}
