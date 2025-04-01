package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ms struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Data   string `json:"data"`
}

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	manager := NewTokenManager()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer ws.Close()
		cookie, err := r.Cookie("access")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		tok, err := manager.Token(cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		_ = ws.WriteJSON(ms{Width: 8, Height: 9, Data: "0,0,1,1,1,1,1,0,1,1,1,0,0,0,1,0,1,2,4,3,0,0,1,0,1,1,1,0,3,2,1,0,1,2,1,1,3,0,1,0,1,0,1,0,2,0,1,1,1,3,0,5,3,3,2,1,1,0,0,0,2,0,0,1,1,1,1,1,1,1,1,1"})
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("[%s]: %s", tok.username, message)

			if err := ws.WriteMessage(messageType, message); err != nil {
				log.Println(err)
				break
			}
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, ok := getPage(r.URL.Path)
		if !ok {
			handleRoute(manager, w, r)
			return
		}
		if page.access {
			access, err := r.Cookie("access")
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					handleAccess(w, r)
					return
				default:
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			token, err := manager.Token(access.Value)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			handleRights(*token, page, w)
		} else {
			handlePage(page, map[string]interface{}{}, w)
		}
	})
	http.FileServer(http.FS(resources))
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
