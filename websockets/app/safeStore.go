package app

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	USER_UUID string = "USER_UUID"
)

type UserInfo struct {
	UserUUID string `json:"user_uuid"`
	Method   string `json:"method"`
	Message  string `json:"message"`
	Conn     *websocket.Conn
}

type SafeStore struct {
	Clients      Queue
	mu           *sync.Mutex
	UserInfo     *UserInfo
	BroadcastMsg chan string
	storeCache   bool
}

func NewStore() *SafeStore {
	return &SafeStore{
		Clients:      Queue{},
		mu:           &sync.Mutex{},
		UserInfo:     &UserInfo{},
		BroadcastMsg: make(chan string),
	}

}

func (ss *SafeStore) Set(userMsg string, ws *websocket.Conn) {
	userInfo := UserInfo{
		UserUUID: userMsg,
		Conn:     ws,
	}

	ss.mu.Lock()
	ss.Clients.Enqueue(userInfo)
	ss.mu.Unlock()
}

func (ss *SafeStore) Get(userUUID string) string {
	var currentMap Stack
	if len(ss.Clients.in) != 0 {
		currentMap = ss.GetAll().in
	} else {
		currentMap = ss.GetAll().out
	}
	for _, val := range currentMap {
		if val.UserUUID == userUUID {
			return val.UserUUID
		}

	}

	return "user not found!"
}

func (ss *SafeStore) Remove() {
	ss.mu.Lock()
	ss.Clients.Dequeue()
	copy(ss.Clients.in, ss.Clients.out)
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
