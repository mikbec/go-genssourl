package main

import (
	"log"
	"net/http"
	"os"

	"local/go-genssourl/ui"
)

func main() {
	mux := http.NewServeMux()

	path_to_static := "./ui/static/"
	_, err := os.Stat(path_to_static)
	if err == nil {
		// Create a file server which serves files out of the "./ui/static" directory.
		// Note that the path given to the http.Dir function is relative to the project
		// directory root.
		log.Print("Trying to use external path ...")
		fileServer := http.FileServer(http.Dir(path_to_static))

		// Use the mux.Handle() function to register the file server as the handler for
		// all URL paths that start with "/static/". For matching paths, we strip the
		// "/static" prefix before the request reaches the file server.
		mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	} else {
		// Take the ui.Content_static embedded filesystem and convert it to a
		// http.FS type so that it satisfies the http.FileSystem interface. We
		// then pass that to the http.FileServer() function to create the file
		// server handler.
		log.Print("Trying to use embedded path ...")
		fileServer := http.FileServer(http.FS(ui.Content_static))

		// Use the mux.Handle() function to register the file server as the handler for
		// all URL paths that start with "/static/". For matching paths, we strip the
		// "/static" prefix before the request reaches the file server.
		mux.Handle("/static/", fileServer)
	}

	// our static Route

	// our own routes
	mux.HandleFunc("/", showHome)

	log.Print("Starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
