package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/fcgi"
	"time"

	"github.com/icza/gog"

	"local/go-genssourl/internal/app"
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		w.Write([]byte("Errro(404): URL '" + r.URL.Path + "' ... resource not found."))
	}
}

func returnNotFound(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, r, http.StatusNotFound)
}

func showHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		//http.NotFound(w, r)
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	w.Write([]byte("Hello from home of GenSSOUrl!"))
}

func doRedirect(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	http.NotFound(w, r)
	//	return
	//}
	idx, ok := webCtxIdxs[r.URL.Path]
	if ok != true {
		//http.NotFound(w, r)
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	if myCfg.CliOpts.OptDebug >= 1 {
		log.Print("doRedirect called for '" + r.URL.Path + "' ....")
	}

	// set attributes
	dstServerProtocol := myCfg.WebCtxs[idx].DstServerProtocol
	dstServerHost := myCfg.WebCtxs[idx].DstServerHost
	dstServerPort := myCfg.WebCtxs[idx].DstServerPort
	dstServerCtx := myCfg.WebCtxs[idx].DstServerCtx
	dstAttrKeyUsername := myCfg.WebCtxs[idx].DstAttrKeyUsername
	dstAttrKeyTimestamp := myCfg.WebCtxs[idx].DstAttrKeyTimestamp
	dstAttrKeyHash := myCfg.WebCtxs[idx].DstAttrKeyHash
	dstAttrKeyId := myCfg.WebCtxs[idx].DstAttrKeyId
	proxyAttrRemoteUsername := myCfg.WebCtxs[idx].ProxyAttrRemoteUsername
	dstAttrValId := myCfg.WebCtxs[idx].DstAttrValId

	// set username from config or from request
	dstAttrValUsername := myCfg.WebCtxs[idx].DstAttrValUsername
	if dstAttrValUsername == "" && myCfg.WebCtxs[idx].ProxyAttrRemoteUsernames != nil {
		if myCfg.CliOpts.OptDebug >= 2 {
			log.Print("Trying ProxyAttrRemoteUsernames array ....")
		}
		for _, unVar := range myCfg.WebCtxs[idx].ProxyAttrRemoteUsernames {
			if myCfg.CliOpts.OptDebug >= 2 {
				log.Print("Trying ProxyAttrRemoteUsernames variable " + unVar + "....")
			}
			// first try FCGI environment variable
			env := fcgi.ProcessEnv(r)
			dstAttrValUsername, ok = env[unVar]
			if ok != true {
				// the try HTTP header variable
				dstAttrValUsername = r.Header.Get(unVar)
			}
			if dstAttrValUsername != "" {
				break
			}
		}
	} else if dstAttrValUsername == "" && proxyAttrRemoteUsername != "" {
		if myCfg.CliOpts.OptDebug >= 2 {
			log.Print("Trying ProxyAttrRemoteUsername variable " + proxyAttrRemoteUsername + " ....")
		}
		// first try FCGI environment variable
		env := fcgi.ProcessEnv(r)
		dstAttrValUsername, ok = env[proxyAttrRemoteUsername]
		if ok != true {
			// the try HTTP header variable
			dstAttrValUsername = r.Header.Get(proxyAttrRemoteUsername)
		}
	}
	if dstAttrValUsername == "" {
		log.Print("Warning: Could not find any authenticated username ... please check config for proxyAttrRemoteUsername or proxyAttrRemoteUsernames[].")
	}

	// set from config or from now
	dstAttrValTimestamp := myCfg.WebCtxs[idx].DstAttrValTimestamp
	dstAttrValTimestampFormat := myCfg.WebCtxs[idx].DstAttrValTimestampFormat
	dstAttrValTimezone := myCfg.WebCtxs[idx].DstAttrValTimezone
	if dstAttrValTimestamp == "" {
		t := time.Now()
		if dstAttrValTimezone == "" || dstAttrValTimezone == "UTC" {
			// use UTC time
			dstAttrValTimestamp = t.UTC().Format(dstAttrValTimestampFormat)
		} else {
			// Load the time zone location
			loc, err := time.LoadLocation(dstAttrValTimezone)
			if err != nil {
				log.Print("Warning: Could not load time zone '" + dstAttrValTimezone + "' ... please check config.")
				dstAttrValTimestamp = ""
			} else {
				// Get the current time at a location
				t := time.Now().In(loc)
				dstAttrValTimestamp = t.Format(dstAttrValTimestampFormat)
			}
		}
	}

	// calculate hash value
	dstServerCertPemFile := myCfg.WebCtxs[idx].DstServerCertPemFile
	algorithmToUseForHash := myCfg.WebCtxs[idx].AlgorithmToUseForHash
	hashVal, _ := app.HexStringOfEncryptedHashValue(
		dstAttrValUsername+dstAttrValTimestamp,
		algorithmToUseForHash,
		dstServerCertPemFile)

	// now generate url string
	urlString := fmt.Sprintf("%s://%s%s%s%s?%s=%s&%s=%s&%s=%s%s%s",
		dstServerProtocol,
		dstServerHost,
		gog.If(dstServerPort == "", "", ":"), dstServerPort,
		dstServerCtx,
		dstAttrKeyUsername, dstAttrValUsername,
		dstAttrKeyTimestamp, dstAttrValTimestamp,
		dstAttrKeyHash, hashVal,
		gog.If(dstAttrValId == "", "", "&"+dstAttrKeyId+"="), dstAttrValId)

	if myCfg.CliOpts.OptDebug >= 3 {
		w.Write([]byte("GenPortURL is redirecting to:\n" + urlString + "\n"))
	}
	http.Redirect(w, r, urlString, http.StatusSeeOther)
}
