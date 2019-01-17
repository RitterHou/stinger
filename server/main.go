package main

import (
	"flag"
	"github.com/ritterhou/stinger/core/codec"
	"github.com/ritterhou/stinger/core/mylog"
	"github.com/ritterhou/stinger/core/network"
	serverConf "github.com/ritterhou/stinger/server/conf"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
)

var (
	confFile string
	password string
)

func main() {
	flag.StringVar(&confFile, "c", "stinger_server.yaml", "Server configuration file.")
	serverConf.LoadConf(confFile)

	conf := serverConf.GetConf()
	mylog.InitLog(conf.LogFile, conf.LogLevel)

	codec.SetKey(conf.Password)
	startProxyServer(conf.ServerPort)
}

func startProxyServer(proxyPort int) {
	var l net.Listener
	var err error
	var host = "0.0.0.0:" + strconv.Itoa(proxyPort)

	l, err = net.Listen("tcp", host)
	if err != nil {
		logrus.Fatal("Error listening:", err)
	}
	defer l.Close()

	logrus.Info("Server listening on " + host)
	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.Warn("Error accepting:", err)
			continue
		}

		//log.Printf("Connection established %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go handlerClient(network.New(conn))
	}
}

func handlerClient(localConn network.Connection) {
	clientPwdBytes, err := localConn.ReadWithLength()
	if err != nil {
		logrus.Warn(err)
		return
	}
	clientPwd := string(clientPwdBytes)
	if clientPwd != password {
		logrus.Printf("client password %s not equals %s", clientPwd, password)
		err = localConn.Write([]byte{1})
		if err != nil {
			localConn.Close()
			logrus.Warn(err)
		}
		return
	}
	err = localConn.Write([]byte{0}) // 验证成功
	if err != nil {
		localConn.Close()
		logrus.Warn(err)
		return
	}

	targetAddrBytes, err := localConn.ReadWithLength()
	if err != nil {
		logrus.Warn(err)
		return
	}
	targetAddr := string(targetAddrBytes)
	//log.Println(targetAddr)

	remoteConn, err := network.Connect(targetAddr)
	if err != nil {
		logrus.Info("can't connect to target address", targetAddr)
		err = localConn.Write([]byte{1}) // 远程主机连接失败
		if err != nil {
			localConn.Close()
			logrus.Warn(err)
		}
		return
	}

	err = localConn.Write([]byte{0}) // 连接成功
	if err != nil {
		localConn.Close()
		logrus.Warn(err)
		return
	}

	go func() {
		for {
			// local -> server
			buf, err := localConn.ReadWithLength()
			if err != nil {
				logrus.Info("local -> server " + localConn.RemoteAddress() + " -> " + err.Error())
				remoteConn.Close()
				break
			}
			buf = codec.Decrypt(buf)
			// server -> remote
			err = remoteConn.Write(buf)
			if err != nil {
				logrus.Info("server -> remote " + remoteConn.RemoteAddress() + " -> " + err.Error())
				localConn.Close()
				break
			}
		}
	}()

	go func() {
		for {
			// remote -> server
			buf, err := remoteConn.Read(1024)
			if err != nil {
				logrus.Info("remote -> server " + remoteConn.RemoteAddress() + " -> " + err.Error())
				localConn.Close()
				break
			}
			buf = codec.Encrypt(buf)
			// server -> local
			err = localConn.WriteWithLength(buf)
			if err != nil {
				logrus.Info("server -> local " + localConn.RemoteAddress() + " -> " + err.Error())
				remoteConn.Close()
				break
			}
		}
	}()
}
