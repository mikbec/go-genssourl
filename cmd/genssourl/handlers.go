package main

import (
	"fmt"
	"net/http"

	"github.com/icza/gog"

	"local/go-genssourl/internal/app"
)

func showHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from home of GenSSOUrl!"))
}

func doRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// set attributes
	server_protocol := "https"
	server_host := "localhost"
	server_port := "5677"
	server_context := "bla"
	url_attr_username_key := "user"
	url_attr_timestamp_key := "ts"
	url_attr_hash_key := "hash"
	url_attr_id_key := "id"

	username_val := "user1"
	timestamp_val := "2023-11-23T08:15:32Z"
	pub_pem_file := "testdata/test.crt.pem"
	hash_val, _ := app.HexStringOfEncryptedHashValue(username_val+timestamp_val, "md5", pub_pem_file)
	id_val := "test"

	urlString := fmt.Sprintf("%s://%s%s%s/%s?%s=%s&%s=%s&%s=%s%s%s",
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
