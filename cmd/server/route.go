package main

import (
	"net/http"
)

func handleRoute(m tokenManager, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/login" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if username == "" || password == "" {
			http.Error(w, "username or password is empty", http.StatusBadRequest)
			return
		}
		access, err := m.Access(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "access", Value: access.text})
		//handlePage(pages["/error"], map[string]interface{}{"error": errors.New("not_access")}, w)
		http.Redirect(w, r, "/game", http.StatusFound)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
