package app

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type UserInfo struct {
	USERID  string `json:"user_id"`
	Method  string `json:"method"`
	Message string `json:"msg"`
	Conn    *websocket.Conn
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

func (ss *SafeStore) Set(userId string, ws *websocket.Conn) *UserInfo {
	userInfo := UserInfo{
		USERID:  userId,
		Method:  "USER_INFO",
		Message: "",
		Conn:    ws,
	}

	ss.mu.Lock()
	ss.Clients.Enqueue(userInfo)
	ss.mu.Unlock()

	return &userInfo
}

func (ss *SafeStore) Get(USERID string, ws *websocket.Conn) string {
	for _, val := range ss.GetAll().in {
		if val.USERID == USERID && val.Conn == ws {
			return val.USERID
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
