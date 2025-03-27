package main

import (
	"errors"
	"fmt"
	"net/http"
)

func main() {
	manager := NewTokenManager()
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
