package main

import (
	//"fmt"
	"net/http"
)

func showHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from GenPortURL!"))
}
