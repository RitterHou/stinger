package http

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/ritterhou/stinger/core/common"
	"github.com/ritterhou/stinger/local/socks"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var download, upload uint64

// 计算带宽以及流量
func bandwidthTraffic() {
	logrus.Info("Monitoring bandwidth traffic.")

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

// 与网页端进行WebSocket连接
func ws(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrade.Upgrade(w, req, nil)
	if err != nil {
		logrus.Warn(err)
		return
	}

	messageType, p, err := conn.ReadMessage()
	if err != nil {
		logrus.Warn(err)
		return
	}
	separator := string(p)
	logrus.Info("The separator is ", separator)

	ticker := time.NewTicker(1 * time.Second)
	lastDownload := download
	lastUpload := upload
	for range ticker.C {
		if lastDownload != download || lastUpload != upload {
			lastDownload = download
			lastUpload = upload
			message := fmt.Sprintf("%s%s%s", common.ByteFormat(download), separator, common.ByteFormat(upload))
			if err := conn.WriteMessage(messageType, []byte(message)); err != nil {
				logrus.Warn(err)
				conn.Close()
				break
			}
		}
	}
	logrus.Info("Stop sending traffic to ", conn.RemoteAddr())
}
