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

func showHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
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
		http.NotFound(w, r)
		return
	}
	log.Print("doRedirect called for '" + r.URL.Path + "' ....")

	// set attributes
	dstServerProtocol := myCfg.WebCtxs[idx].DstServerProtocol
	dstServerHost := myCfg.WebCtxs[idx].DstServerHost
	dstServerPort := myCfg.WebCtxs[idx].DstServerPort
	dstServerCtx := myCfg.WebCtxs[idx].DstServerCtx
	dstAttrKeyUsername := myCfg.WebCtxs[idx].DstAttrKeyUsername
	dstAttrKeyTimestamp := myCfg.WebCtxs[idx].DstAttrKeyTimestamp
	dstAttrKeyHash := myCfg.WebCtxs[idx].DstAttrKeyHash
	dstAttrKeyId := myCfg.WebCtxs[idx].DstAttrKeyId
	proxyAttrRemoteUserName := myCfg.WebCtxs[idx].ProxyAttrRemoteUserName
	dstAttrValId := myCfg.WebCtxs[idx].DstAttrValId

	// set username from config or from request
	dstAttrValUsername := myCfg.WebCtxs[idx].DstAttrValUsername
	if dstAttrValUsername == "" {
		env := fcgi.ProcessEnv(r)
		dstAttrValUsername, ok = env[proxyAttrRemoteUserName]
		if ok != true {
			dstAttrValUsername = ""
		}
	}

	// set from config or from now
	dstAttrValTimestamp := myCfg.WebCtxs[idx].DstAttrValTimestamp
	if dstAttrValTimestamp == "" {
		t := time.Now()
		dstAttrValTimestamp = t.UTC().Format(myCfg.WebCtxs[idx].DstAttrValTimestampFormat)
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
	w.Write([]byte("Hello from GenPortURL! " + urlString))
}
