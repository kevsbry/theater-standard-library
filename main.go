package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/kevsbry/theater/api"
	"github.com/kevsbry/theater/handler"
)

func main() {
	// file handler
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/mediaSource/", http.StripPrefix("/mediaSource/", http.FileServer(http.Dir("E:/videos/Movies"))))
	http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(`D:\myVideos`))))

	//get local ip address
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	//handler
	http.HandleFunc("/theater/", handler.Theater)
	http.HandleFunc("/api/", api.Theater)

	//start server
	port := ":9998"
	address := strings.Split(fmt.Sprint(localAddr), ":")
	go fmt.Println(fmt.Sprintf("<------ %s%s ------>", address[0], port))
	serverErr := http.ListenAndServe(port, nil)

	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
