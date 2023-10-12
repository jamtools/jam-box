package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var connections []*websocket.Conn

func addConnection(conn *websocket.Conn) {
	fmt.Println("adding connection")
	connections = append(connections, conn)
}

func removeConnection(conn *websocket.Conn) {
	fmt.Println("removing connection")
	result := []*websocket.Conn{}
	for _, c := range connections {
		if c != conn {
			result = append(result, c)
		}
	}

	connections = result
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Connection") != "Upgrade" {
		s := RenderWesocketScript()

		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(s))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		removeConnection(conn)
		return nil
	})

	addConnection(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func SendMessage(s string) {
	for _, conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(s))
		if err != nil {
			fmt.Println(err)
		}
	}
}
