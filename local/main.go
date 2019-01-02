package main

import (
	"encoding/binary"
	"fmt"
	"github.com/ritterhou/stinger/common/mylog"
	"github.com/ritterhou/stinger/common/network"
	"net"
	"strconv"
)

const localPort = 2680

var log = mylog.Info

func main() {
	var l net.Listener
	var err error
	var host = "0.0.0.0:" + strconv.Itoa(localPort)

	l, err = net.Listen("tcp", host)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer l.Close()

	log.Println("Listening on " + host + " ...")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting:", err)
			continue
		}

		//log.Printf("Connection established %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		c := network.Connection{Conn: conn}
		go handlerSocks5(c)
	}
}

func handlerSocks5(conn network.Connection) {
	authSocks5(conn)
	remoteConn := connectSocks5(conn)

	log.Printf("Connect success %s -> %s, %s => %s\n",
		conn.RemoteAddress(), conn.LocalAddress(),
		remoteConn.LocalAddress(), remoteConn.RemoteAddress())
	handlerSocks5Data(conn, remoteConn)
}

func authSocks5(conn network.Connection) {
	socksVersion := conn.Read(1)[0]
	if socksVersion != 5 {
		log.Fatal("Socks version should be 5, now is", socksVersion)
	}

	authWaysNum := conn.Read(1)[0]
	authWays := conn.Read(uint32(authWaysNum))
	if !in(byte(0), authWays) {
		log.Fatal("Only support [NO AUTHENTICATION REQUIRED] auth way.")
	}

	conn.Write([]byte{5, 0})
}

func connectSocks5(conn network.Connection) network.Connection {
	socksVersion := conn.Read(1)[0]
	if socksVersion != 5 {
		log.Fatal("Socks version should be 5, now is", socksVersion)
	}

	command := conn.Read(1)[0]
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

	addr := host + ":" + strconv.Itoa(int(port))

	c, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Can't connect to", addr)
	}

	conn.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0})

	remoteConn := network.Connection{Conn: c}
	return remoteConn
}

func handlerSocks5Data(localConn network.Connection, remoteConn network.Connection) {
	go func() {
		for {
			buf := localConn.Read(1024)
			if buf == nil {
				remoteConn.Close()
				break
			}
			remoteConn.Write(buf)
		}
	}()

	go func() {
		for {
			buf := remoteConn.Read(1024)
			if buf == nil {
				localConn.Close()
				break
			}
			localConn.Write(buf)
		}
	}()
}

func in(num byte, list []byte) bool {
	for _, e := range list {
		if e == num {
			return true
		}
	}
	return false
}
