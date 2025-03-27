package main

import (
	"embed"
	"html/template"
	"net/http"
)

type Page struct {
	resource string
	mime     string
	access   bool
	template bool
}

func ctor(resource, mime string, access, template bool) Page {
	return Page{resource, mime, access, template}
}

//go:embed resources
var resources embed.FS
var pages = map[string]Page{
	"/":             ctor("resources/access.gohtml", "text/html; charset=utf-8", false, true),
	"/game":         ctor("resources/game.gohtml", "text/html; charset=utf-8", true, true),
	"/main.css":     ctor("resources/main.css", "text/css; charset=utf-8", false, false),
	"/error":        ctor("resources/error.gohtml", "text/html; charset=utf-8", false, true),
	"/wasm_exec.js": ctor("resources/wasm_exec.js", "text/javascript; charset=utf-8", false, false),
	"/client.wasm":  ctor("resources/client.wasm", "application/wasm", false, false),
}

func getPage(name string) (Page, bool) {
	page, ok := pages[name]
	return page, ok
}

func handlePage(page Page, data map[string]interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", page.mime)
	if page.template {
		tpl, err := template.ParseFS(resources, page.resource)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := tpl.Execute(w, data); err != nil {
			return
		}
	} else {
		bytes, err := resources.ReadFile(page.resource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bytes)
	}
}

func handleRights(token token, page Page, w http.ResponseWriter) {
	handlePage(page, map[string]interface{}{"token": token.text}, w)
}

func handleAccess(w http.ResponseWriter, r *http.Request) {
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
