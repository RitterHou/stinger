package http

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/ritterhou/stinger/core/common"
	"github.com/ritterhou/stinger/local/socks"
	"log"
	"net/http"
	"time"
	"io"
	"github.com/ritterhou/stinger/local/resource"
	"strconv"
)

var download, upload uint64

// 计算带宽以及流量
func bandwidthTraffic() {
	log.Printf("Moniting bandwidth traffic.")

	ticker := time.NewTicker(1 * time.Second)
	lastDownload := socks.TotalDownload
	lastUpload := socks.TotalUpload
	for range ticker.C {
		t := time.Now()
		now := t.Format("2006-01-02 15:04:05")

		download = socks.TotalDownload - lastDownload
		upload = socks.TotalUpload - lastUpload
		if upload != 0 && download != 0 {
			fmt.Printf("%s %s ↓ %s ↑", now, common.ByteFormat(download), common.ByteFormat(upload))
			fmt.Printf("    (%s ↓ %s ↑)\n", common.ByteFormat(socks.TotalDownload), common.ByteFormat(socks.TotalUpload))
		}
		lastDownload = socks.TotalDownload
		lastUpload = socks.TotalUpload
	}
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// 流量跟踪数据
func traffic(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrade.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	messageType, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	seperater := string(p)
	log.Println("The seperater is", seperater)

	ticker := time.NewTicker(1 * time.Second)
	lastDownload := download
	lastUpload := upload
	for range ticker.C {
		if lastDownload != download || lastUpload != upload {
			lastDownload = download
			lastUpload = upload
			message := fmt.Sprintf("%s%s%s", common.ByteFormat(download), seperater, common.ByteFormat(upload))
			if err := conn.WriteMessage(messageType, []byte(message)); err != nil {
				log.Println(err)
				conn.Close()
				break
			}
		}
	}
	log.Println("Stop sending traffic to", conn.RemoteAddr())
}

var indexHtml string

func index(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, indexHtml)
}

func StartServer(port int) {
	go bandwidthTraffic()

	indexHtml = resource.GetContent("/html/index.html")

	// 首页
	http.HandleFunc("/", index)

	// PAC文件获取
	pacConf := getPac()
	http.HandleFunc("/pac", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s fetched PAC file\n", req.RemoteAddr)
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		io.WriteString(w, pacConf)
	})

	// 获取流量以及网速信息
	http.HandleFunc("/traffic", traffic)

	log.Printf("HTTP Server working on http://0.0.0.0:%d\n", port)
	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
