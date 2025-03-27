package main

import (
	"embed"
	"html/template"
	"net/http"
)

type page struct {
	resource string
	mime     string
	access   bool
	template bool
}

func ctor(resource, mime string, access, template bool) page {
	return page{resource, mime, access, template}
}

//go:embed resources
var resources embed.FS
var pages = map[string]page{
	"/":             ctor("resources/access.gohtml", "text/html; charset=utf-8", false, true),
	"/game":         ctor("resources/game.gohtml", "text/html; charset=utf-8", true, true),
	"/main.css":     ctor("resources/main.css", "text/css; charset=utf-8", false, false),
	"/error":        ctor("resources/error.gohtml", "text/html; charset=utf-8", false, true),
	"/wasm_exec.js": ctor("resources/wasm_exec.js", "text/javascript; charset=utf-8", false, false),
	"/client.wasm":  ctor("resources/client.wasm", "application/wasm", false, false),
}

func handlePage(page page, data map[string]interface{}, w http.ResponseWriter) {
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

func handleRights(token token, page page, w http.ResponseWriter) {
	handlePage(page, map[string]interface{}{"token": token.text}, w)
}
