//go:generate statik -src=./assets
//go:generate go fmt statik/statik.go

package main

import (
	"flag"
	"github.com/ritterhou/stinger/core/codec"
	"github.com/ritterhou/stinger/core/mylog"
	"github.com/ritterhou/stinger/core/network"
	localConf "github.com/ritterhou/stinger/local/conf"
	"github.com/ritterhou/stinger/local/http"
	"github.com/ritterhou/stinger/local/socks"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
)

var (
	confFile     string
	remoteServer string
	password     string
)

func main() {
	flag.StringVar(&confFile, "c", "stinger_local.yaml", "Local configuration file.")

	localConf.LoadConf(confFile)
	conf := localConf.GetConf()

	logFile := conf["log_file"].(string)
	logLevel := conf["log_level"].(string)
	mylog.InitLog(logFile, logLevel)

	pac := conf["pac"].(map[interface{}]interface{})
	pacPort := pac["port"].(int)
	global := pac["global"].(bool)
	domains := pac["domains"].([]interface{})

	localPort := conf["local_port"].(int)
	remoteServer = conf["server_address"].(string)

	pwd := conf["password"]
	switch v := pwd.(type) {
	case int:
		password = strconv.Itoa(v)
	case string:
		password = v
	default:
		logrus.Warn("Unknown type ", v)
	}

	codec.SetKey(password)

	http.CreatePacFile(localPort, global, domains)
	go http.StartServer(pacPort)

	startProxyServer(localPort)
}

func startProxyServer(proxyPort int) {
	var l net.Listener
	var err error
	var host = "0.0.0.0:" + strconv.Itoa(proxyPort)

	l, err = net.Listen("tcp", host)
	if err != nil {
		logrus.Fatal("Error listening: ", err)
	}
	defer l.Close()

	logrus.Info("Local ProxyServer working on " + host)
	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.Warn("Error accepting ", err)
			continue
		}
		go handlerSocks5(network.New(conn))
	}
}

func handlerSocks5(conn network.Connection) {
	err := socks.AuthSocks5(conn)
	if err != nil {
		logrus.Warn(err)
		return
	}
	remoteConn, err := socks.ConnectRemote(conn, remoteServer, password)
	if err != nil {
		logrus.Warn(err)
		return
	}

	socks.HandlerSocks5Data(conn, remoteConn)
}
