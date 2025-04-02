package app

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type UserInfo struct {
	UserUUID string
	Conn     *websocket.Conn
}

type SafeStore struct {
	Clients      Queue
	mu           sync.Mutex
	BroadcastMsg chan string
}

func NewStore() *SafeStore {
	return &SafeStore{
		Clients:      Queue{},
		mu:           sync.Mutex{},
		BroadcastMsg: make(chan string),
	}

}

func (ss *SafeStore) Set(userUUID string, ws *websocket.Conn) {
	userInfo := UserInfo{
		UserUUID: userUUID,
		Conn:     ws,
	}

	ss.mu.Lock()
	ss.Clients.Enqueue(userInfo)
	ss.mu.Unlock()
}

func (ss *SafeStore) Remove() {
	ss.mu.Lock()
	ss.Clients.Dequeue()
	ss.mu.Unlock()
}

func (ss *SafeStore) GetAll() Queue {
	ss.mu.Lock()
	tempMap := ss.Clients
	ss.mu.Unlock()

	return tempMap
}

func (ss *SafeStore) PrintAll() {
	ss.mu.Lock()
	for value := range ss.Clients.in {
		fmt.Println("the key: ", value)
	}
	ss.mu.Unlock()

	fmt.Println("the number of sessions: ", len(ss.Clients.in))
}
