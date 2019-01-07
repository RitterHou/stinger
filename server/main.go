package main

import (
	"flag"
	"fmt"
	"github.com/ritterhou/stinger/core/codec"
	"github.com/ritterhou/stinger/core/common"
	"github.com/ritterhou/stinger/core/network"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func init() {
	t := time.Now()
	now := t.Format("20060102150405")
	filename := fmt.Sprintf("stinger_server.%s.log", now)

	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("open log file failed", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var (
	confFile string
	password string
)

func main() {
	flag.StringVar(&confFile, "c", "stinger_server.yaml", "Server configuration file.")

	path := common.GetAbsPath(confFile)
	content := common.ReadFile(path)
	conf := common.MarshalYaml(content)

	serverPort := conf["server_port"].(int)

	pwd := conf["password"]
	switch v := pwd.(type) {
	case int:
		password = strconv.Itoa(v)
	case string:
		password = v
	default:
		log.Println("Unknown type ", v)
	}

	codec.SetKey(password)

	startProxyServer(serverPort)
}

func startProxyServer(proxyPort int) {
	var l net.Listener
	var err error
	var host = "0.0.0.0:" + strconv.Itoa(proxyPort)

	l, err = net.Listen("tcp", host)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer l.Close()

	log.Println("Server listening on " + host)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting:", err)
			continue
		}

		//log.Printf("Connection established %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		c := network.Connection{Conn: conn}
		go handlerClient(c)
	}
}

func handlerClient(localConn network.Connection) {
	clientPwdBytes, err := localConn.ReadWithLength()
	if err != nil {
		log.Println(err)
		return
	}
	clientPwd := string(clientPwdBytes)
	if clientPwd != password {
		log.Printf("client password %s not equals %s\n", clientPwd, password)
		localConn.Write([]byte{1})
		return
	}
	localConn.Write([]byte{0}) // 验证成功

	targetAddrBytes, err := localConn.ReadWithLength()
	if err != nil {
		log.Println(err)
		return
	}
	targetAddr := string(targetAddrBytes)
	//log.Println(targetAddr)

	c, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Println("can't connect to target address", targetAddr)
		localConn.Write([]byte{1}) // 远程主机连接失败
		return
	}

	localConn.Write([]byte{0}) // 连接成功
	remoteConn := network.Connection{Conn: c}

	go func() {
		for {
			// local -> server
			buf, err := localConn.ReadWithLength()
			if err != nil {
				log.Println(localConn.RemoteAddress() + " -> " + err.Error())
				remoteConn.Close()
				break
			}
			buf = codec.Decrypt(buf)
			// server -> remote
			remoteConn.Write(buf)
		}
	}()

	go func() {
		for {
			// remote -> server
			buf, err := remoteConn.Read(1024)
			if err != nil {
				log.Println(remoteConn.RemoteAddress() + " -> " + err.Error())
				localConn.Close()
				break
			}
			buf = codec.Encrypt(buf)
			// server -> local
			localConn.WriteWithLength(buf)
		}
	}()
}
