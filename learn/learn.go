package learn

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type BroadCast struct {
	Conn    *websocket.Conn
	Message []byte
}

var clientConnectionMap = make(map[*websocket.Conn]bool)
var broadcastChannel = make(chan BroadCast)
var mutex = &sync.Mutex{}

func Start() {
	log.Println("Web Sockets Service")
	go handleMessage()
	http.HandleFunc("/ws", wsHandler)
	log.Println("Starting web-socket server on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		fmt.Println("CheckOrigin : ", r.Header.Get("Origin"))
		fmt.Println("R.Host : ", r.Host)
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error while upgrading connection: %v", err)
		panic(err)
	}

	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close()
	mutex.Lock()
	clientConnectionMap[conn] = true
	mutex.Unlock()
	for {
		log.Printf("Remote address : %s \n\n", conn.RemoteAddr().String())
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clientConnectionMap, conn)
			mutex.Unlock()
			log.Printf("Error while reading message: %v", err)
			break
		}
		log.Printf("Received message: %s", string(msg))
		log.Printf("Message type: %d", msgType)

		// broacast
		broadcastChannel <- BroadCast{
			Conn:    conn,
			Message: msg,
		}

		// responseMessage := fmt.Sprintf("Message from server -> %s", string(msg))
		// err = conn.WriteMessage(websocket.TextMessage, []byte(responseMessage))
		// if err != nil {
		// 	log.Printf("Error while writing message: %v", err)
		// 	break
		// }
	}
}

func handleMessage() {
	for {
		broadcastMessage := <-broadcastChannel
		mutex.Lock()
		var err error
		for _conn := range clientConnectionMap {
			if _conn == broadcastMessage.Conn {
				continue
			}
			err = _conn.WriteMessage(websocket.TextMessage, broadcastMessage.Message)
			if err != nil {
				_conn.Close()
				delete(clientConnectionMap, _conn)
			}
		}
		mutex.Unlock()
	}
}
