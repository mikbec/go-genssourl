package main

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"strconv"

	"local/go-genssourl/ui"
)

var webCtxIdxs = make(map[string]int)

func main() {
	// first scan our Configuration
	scanConfiguration()

	// print configuration as JSON if wanted
	if myCfg.CliOpts.OptCfgAsJSON == true {
		printCfgAsJSON()
		defer func() { os.Exit(0) }()
		return
	}

	// print configuration as YAML if wanted
	if myCfg.CliOpts.OptCfgAsYAML == true {
		printCfgAsYAML()
		defer func() { os.Exit(0) }()
		return
	}

	// get a new mux
	mux := http.NewServeMux()

	// decide if we use our internal or extarnal pages/templates
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
	//mux.HandleFunc("/", doRedirect)
	//mux.HandleFunc("/home", showHome)
	for idx := 0; idx < len(myCfg.WebCtxs); idx++ {
		str := myCfg.WebCtxs[idx].ThisServerPath
		if str == "" {
			str = "/"
		}
		if idx2, ok := webCtxIdxs[str]; ok {
			log.Print("Warning: Index" + strconv.Itoa(idx) + " overwrites index " + strconv.Itoa(idx2))
		}
		webCtxIdxs[str] = idx
		log.Print("Adding web context '" + str + "' (index " + strconv.Itoa(idx) + ") to list of contexts ...")

		mux.HandleFunc(str, doRedirect)
	}

	// Install NotFound for all resources which are not found
	_, ok := webCtxIdxs["/"]
	if ok != true {
		mux.HandleFunc("/", returnNotFound)
	}

	// start server depending on commandline options
	// got idea from:
	//	https://muzzarelli.net/blog/2013/09/how-to-use-go-and-fastcgi/
	// Run as a local or remote web server
	if myCfg.CliOpts.OptCfgSvcWebTcp != "" {
		myCfg.CliOpts.OptCfgSvcFcgiStdIO = false
		// Run as HTTPS web server
		if (myCfg.CliOpts.OptCfgSvcWebTcpCertFile != "") && (myCfg.CliOpts.OptCfgSvcWebTcpKeyFile != "") {
			log.Print("HTTPS server is listening on https://" + myCfg.CliOpts.OptCfgSvcWebTcp)
			err = http.ListenAndServeTLS(
				myCfg.CliOpts.OptCfgSvcWebTcp,
				myCfg.CliOpts.OptCfgSvcWebTcpCertFile,
				myCfg.CliOpts.OptCfgSvcWebTcpKeyFile,
				mux)

			// Run as plain HTTP web server
		} else {
			log.Print("HTTP server is listening on http://" + myCfg.CliOpts.OptCfgSvcWebTcp)
			err = http.ListenAndServe(myCfg.CliOpts.OptCfgSvcWebTcp, mux)
		}

		// Run as FCGI via TCP
	} else if myCfg.CliOpts.OptCfgSvcFcgiTcp != "" {
		myCfg.CliOpts.OptCfgSvcFcgiStdIO = false
		listener, err := net.Listen("tcp", myCfg.CliOpts.OptCfgSvcFcgiTcp)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		log.Print("FCGI server is listening on tcp://" + myCfg.CliOpts.OptCfgSvcFcgiTcp)
		err = fcgi.Serve(listener, mux)

		// Run as FCGI via UNIX socket
	} else if myCfg.CliOpts.OptCfgSvcFcgiUnix != "" {
		myCfg.CliOpts.OptCfgSvcFcgiStdIO = false
		log.Print("FCGI server is listening on unix://" + myCfg.CliOpts.OptCfgSvcFcgiUnix)
		listener, err := net.Listen("unix", myCfg.CliOpts.OptCfgSvcFcgiUnix)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, mux)
		// Run as FCGI via standard I/O
	} else {
		myCfg.CliOpts.OptCfgSvcFcgiStdIO = true
		log.Print("FCGI server is listening on STD I/O")
		err = fcgi.Serve(nil, mux)
	}
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Server is going down...")
	}
}
