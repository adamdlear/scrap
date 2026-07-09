package main

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/net/websocket"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func wsHandler(ws *websocket.Conn) {
	for {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error: %v\n", err)
		}
		words := strings.Split(msg, " ")
		slices.Reverse(words)
		revMsg := strings.Join(words, " ")
		websocket.Message.Send(ws, revMsg)
	}
}

func main() {
	http.HandleFunc("/health", health)
	http.Handle("/ws", websocket.Handler(wsHandler))

	fmt.Println("Listening on port :8080")
	http.ListenAndServe(":8080", nil)
}
