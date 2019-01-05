package main

import (
	"flag"
	"github.com/ritterhou/stinger/core/network"
	localConf "github.com/ritterhou/stinger/local/conf"
	"github.com/ritterhou/stinger/local/http"
	"github.com/ritterhou/stinger/local/socks"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	confFile     string
	remoteServer string
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	flag.StringVar(&confFile, "c", "stinger_local.yaml", "Local configuration file.")

	localConf.LoadConf(confFile)
	conf := localConf.GetConf()

	pac := conf["pac"].(map[interface{}]interface{})
	pacPort := pac["port"].(int)
	global := pac["global"].(bool)
	domains := pac["domains"].([]interface{})

	localPort := conf["local_port"].(int)
	remoteServer = conf["server_address"].(string)

	http.CreatePacFile(localPort, global, domains)
	go http.StartServer(pacPort)

	go socks.BandwidthTraffic()

	startProxyServer(localPort)
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

	log.Println("Local ProxyServer working on " + host)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting:", err)
			continue
		}
		go handlerSocks5(network.Connection{Conn: conn})
	}
}

func handlerSocks5(conn network.Connection) {
	err := socks.AuthSocks5(conn)
	if err != nil {
		log.Println(err)
		return
	}
	remoteConn, err := socks.ConnectRemote(conn, remoteServer)
	if err != nil {
		log.Println(err)
		return
	}

	//log.Printf("Connect success %s -> %s, %s => %s\n", conn.RemoteAddress(), conn.LocalAddress(), remoteConn.LocalAddress(), remoteConn.RemoteAddress())
	socks.HandlerSocks5Data(conn, remoteConn)
}
