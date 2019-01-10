package http

import (
	"github.com/ritterhou/stinger/local/resource"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

var indexHtml string

func StartServer(port int) {
	go bandwidthTraffic()

	indexHtml = resource.GetContent("/html/index.html")

	// 首页
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, indexHtml)
	})

	// PAC文件获取
	pacConf := getPac()
	http.HandleFunc("/pac", func(w http.ResponseWriter, req *http.Request) {
		logrus.Printf("%s fetched PAC file", req.RemoteAddr)
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		io.WriteString(w, pacConf)
	})

	// WebSocket
	http.HandleFunc("/ws", ws)

	logrus.Printf("HTTP Server working on http://0.0.0.0:%d", port)
	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
	if err != nil {
		logrus.Fatal("ListenAndServe: ", err)
	}
}
