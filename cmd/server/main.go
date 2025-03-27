package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

func handleAccess(m tokenManager, w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFS(resources, "resources/access.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	data := map[string]interface{}{
		"agent": r.Header.Get("User-Agent"),
	}
	if err := tpl.Execute(w, data); err != nil {
		return
	}
}

func main() {
	manager := newTokenManager()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page, ok := pages[r.URL.Path]
		if !ok {
			handleRoute(manager, w, r)
			return
		}
		if page.access {
			access, err := r.Cookie("access")
			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					handleAccess(manager, w, r)
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
