// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address")
var logObj = log.New(os.Stdout, "BattleShips: ", log.Ldate|log.Ltime)

func serveHome(w http.ResponseWriter, r *http.Request) {
	logObj.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	flag.Parse()

	initHubs()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/favicon.ico").Handler(http.FileServer(http.Dir("static")))

	r.HandleFunc("/", serveHome)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	srv := &http.Server{
		Handler: r,
		Addr:    *addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logObj.Print("Server Started and Listening")
	err := srv.ListenAndServe()
	if err != nil {
		logObj.Fatal("ListenAndServe: ", err)
	}
}
