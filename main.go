package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

)

type ContainerOpts struct {
	Namespace string	`yaml:"namespace"`
	Pod       string	`yaml:"pod"`
	Container string	`yaml:"container"`
	Command   []string	`yaml:"command"`
	TTY       bool		`yaml:"tty"`
	Stdin     bool		`yaml:"stdin"`
}


func homePage (w http.ResponseWriter, r *http.Request) {
	log.Println("Home Page Refreshed")
	if r.URL.Path != "/" {
		http.Error(w, "Page Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func getTestControllerOpts() *ContainerOpts {
	var cmd []string
	cmd = append(cmd, "ls ;")
	return &ContainerOpts{
		Namespace: "default",
		Pod:       "ws-test-1",
		Container: "ws-test-1",
		Command:   cmd,
		TTY:       true,
		Stdin:     true,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsController(writer http.ResponseWriter, request *http.Request) {
	contConfig := getTestControllerOpts()
	fmt.Printf("%+v", contConfig )

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage,[]byte("Hi Fiest conection"))
	if err != nil {
		fmt.Printf("Error %+v",err)
	}
	conn.Close()
}


func main() {
	fmt.Println("Hidden Note !")
	http.HandleFunc("/ws-example", wsController)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":2900",nil)
}