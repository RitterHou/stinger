package socks

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ritterhou/stinger/core/codec"
	"github.com/ritterhou/stinger/core/common"
	"github.com/ritterhou/stinger/core/network"
	"log"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

var totalDownload uint64
var totalUpload uint64

// 显示带宽以及流量
func BandwidthTraffic() {
	log.Printf("Moniting bandwidth traffic.")

	ticker := time.NewTicker(1 * time.Second)
	lastDownload := totalDownload
	lastUpload := totalUpload
	for range ticker.C {
		t := time.Now()
		now := t.Format("2006-01-02 15:04:05")

		download := totalDownload - lastDownload
		upload := totalUpload - lastUpload
		if upload != 0 && download != 0 {
			fmt.Printf("%s %s ↓ %s ↑", now, common.ByteFormat(download), common.ByteFormat(upload))
			fmt.Printf("    (%s ↓ %s ↑)\n", common.ByteFormat(totalDownload), common.ByteFormat(totalUpload))
		}
		lastDownload = totalDownload
		lastUpload = totalUpload
	}
}

func AuthSocks5(conn network.Connection) {
	socksVersion := conn.Read(1)[0]
	if socksVersion != 5 {
		log.Fatal("Socks version should be 5, now is", socksVersion)
	}

	authWaysNum := conn.Read(1)[0]
	authWays := conn.Read(uint32(authWaysNum))
	if !common.In(byte(0), authWays) {
		log.Fatal("Only support [NO AUTHENTICATION REQUIRED] auth way.")
	}

	conn.Write([]byte{5, 0})
}

func ConnectRemote(conn network.Connection, remoteServer string) (network.Connection, error) {
	socksVersion := conn.Read(1)[0]
	if socksVersion != 5 {
		log.Fatal("Socks version should be 5, now is", socksVersion)
	}

	command := conn.Read(1)[0]
	// 仅支持通过TCP协议进行流量转发
	if command != 1 {
		log.Fatal("Only support [CONNECT] command")
	}

	conn.Read(1) // 保留字

	addrType := conn.Read(1)[0]

	var host string
	switch addrType {
	case 1: // ipv4
		data := conn.Read(4)
		host = fmt.Sprintf("%d.%d.%d.%d", data[0], data[1], data[2], data[3])
	case 3: // 域名
		hostLength := conn.Read(1)[0]
		host = string(conn.Read(uint32(hostLength)))
	default:
		log.Fatal("Not support address type", addrType)
	}

	port := binary.BigEndian.Uint16(conn.Read(2))
	// 构建最终目标的地址
	targetAddr := host + ":" + strconv.Itoa(int(port))
	// 尝试连接到远程主机
	c, err := net.Dial("tcp", remoteServer)
	if err != nil {
		conn.Write([]byte{5, 3, 0, 1, 0, 0, 0, 0, 0, 0})
		return network.Connection{}, errors.New("can't connect to remote server " + remoteServer)
	}
	serverConn := network.Connection{Conn: c}
	// 把最终目标的地址（域名或IP）发送到远程主机，由远程主机负责实现连接
	serverConn.WriteWithLength([]byte(targetAddr))
	// 获取远程主机的连接状态
	connectStatus := serverConn.Read(1)[0]
	if connectStatus != 0 {
		serverConn.Close()
		conn.Write([]byte{5, 4, 0, 1, 0, 0, 0, 0, 0, 0})
		return network.Connection{}, errors.New("remote server connect target address failed " + targetAddr)
	}
	// 向客户端发送连接成功的消息
	conn.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0})
	return serverConn, nil
}

func HandlerSocks5Data(localConn network.Connection, remoteConn network.Connection) {
	go func() {
		for {
			// 浏览器 -> local
			buf := localConn.Read(1024)
			if buf == nil {
				remoteConn.Close()
				break
			}

			buf = codec.Encrypt(buf)
			// 记载本地上传的流量
			atomic.AddUint64(&totalUpload, uint64(len(buf)))
			// local -> server
			remoteConn.WriteWithLength(buf)
		}
	}()

	go func() {
		for {
			// server -> local
			buf := remoteConn.ReadWithLength()
			if buf == nil {
				localConn.Close()
				break
			}
			// 记载本地下载的流量
			atomic.AddUint64(&totalDownload, uint64(len(buf)))

			buf = codec.Decrypt(buf)
			// local -> 浏览器
			localConn.Write(buf)
		}
	}()
}
