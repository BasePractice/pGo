package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stdiopt/gowasm-experiments/arty/painter"
	"image/color"
	"log"
	"net/http"
)

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	manager := NewTokenManager()
	p, err := painter.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Init(painter.InitOP{Width: 10, Height: 10})

	err = p.HandleOP(painter.TextOP{
		Color: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		X:     5.0,
		Y:     5.0,
		Text:  "Hello world",
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer ws.Close()
		_ = ws.WriteJSON(painter.Message{Payload: painter.InitOP{
			Width:  p.Width(),
			Height: p.Height(),
			Data:   p.ImageData(),
		}})
		// цикл обработки сообщений
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			log.Printf("Received: %s", message)

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
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
