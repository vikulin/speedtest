package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
)

func getData(w http.ResponseWriter, r *http.Request) {
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 8 * 1024 * 1024
	}
	chunkSize, err := strconv.Atoi(r.URL.Query().Get("chunkSize"))
	if err != nil {
		chunkSize = 64 * 1024
	}
	data := make([]byte, chunkSize)
	chunks := size / chunkSize
	for i := 0; i < chunks; i++ {
		w.Write(data)
	}
}

func setEmpty(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(""))
}

func getIP(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return
	}
	w.Write([]byte(userIP.String()))
}

func main() {
	port := "80"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	basePath := "web"
	if len(os.Args) > 2 {
		basePath = os.Args[2]
	}

	http.Handle("/", http.FileServer(http.Dir(basePath)))
	http.HandleFunc("/download", getData)
	http.HandleFunc("/ip", getIP)
	http.HandleFunc("/ping", setEmpty)
	http.HandleFunc("/upload", setEmpty)

	fmt.Print("Server listening at http://0.0.0.0:" + port + "\n")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
