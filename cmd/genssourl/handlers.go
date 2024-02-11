package main

import (
	"fmt"
	"log"
	"net/http"

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
	server_protocol := myCfg.WebCtxs[idx].DstServerProtocol
	server_host := myCfg.WebCtxs[idx].DstServerHost
	server_port := myCfg.WebCtxs[idx].DstServerPort
	server_context := myCfg.WebCtxs[idx].DstServerCtx
	url_attr_username_key := myCfg.WebCtxs[idx].DstAttrKeyUsername
	url_attr_timestamp_key := myCfg.WebCtxs[idx].DstAttrKeyTimestamp
	url_attr_hash_key := myCfg.WebCtxs[idx].DstAttrKeyHash
	url_attr_id_key := myCfg.WebCtxs[idx].DstAttrKeyId

	username_val := myCfg.WebCtxs[idx].DstAttrValUsername
	//timestamp_val := myCfg.WebCtxs[idx].DstAttrValTimestamp
	timestamp_val := "2023-11-23T08:15:32Z"
	id_val := myCfg.WebCtxs[idx].DstAttrValId
	pub_pem_file := myCfg.WebCtxs[idx].DstServerCertPemFile
	hash_algo := myCfg.WebCtxs[idx].AlgorithmToUseForHash
	hash_val, _ := app.HexStringOfEncryptedHashValue(
		username_val+timestamp_val,
		hash_algo,
		pub_pem_file)

	urlString := fmt.Sprintf("%s://%s%s%s%s?%s=%s&%s=%s&%s=%s%s%s",
		server_protocol,
		server_host,
		gog.If(server_port == "", "", ":"), server_port,
		server_context,
		url_attr_username_key, username_val,
		url_attr_timestamp_key, timestamp_val,
		url_attr_hash_key, hash_val,
		gog.If(id_val == "", "", "&"+url_attr_id_key+"="), id_val)
	w.Write([]byte("Hello from GenPortURL! " + urlString))
}
