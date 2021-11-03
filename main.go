package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"stash/server"
	"time"
)

var (
	VERSION string = "1.0.0"
	HOST    string = "127.0.0.1:5050"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Stash v" + VERSION + " @IWick - Starting...")

	http.HandleFunc("/", server.Home)
	http.HandleFunc("/download/", server.Download)

	fmt.Println("[API] Listening on " + HOST)
	err := http.ListenAndServe(HOST, nil)
	if err != nil {
		fmt.Println("[PANIC] Failed to host server")
		fmt.Println("      | " + err.Error())
		return
	}
}
