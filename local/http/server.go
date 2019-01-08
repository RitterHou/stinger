package http

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

func StartServer(port int) {
	pacConf := getPac()
	http.HandleFunc("/pac", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s fetched PAC file\n", req.RemoteAddr)
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		io.WriteString(w, pacConf)
	})
	log.Printf("HTTP Server working on http://0.0.0.0:%d\n", port)
	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
