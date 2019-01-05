package main

import (
	"github.com/ritterhou/stinger/core/codec"
	"github.com/ritterhou/stinger/core/network"
	"log"
	"net"
	"os"
	"strconv"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const localPort = 26800

func main() {
	var l net.Listener
	var err error
	var host = "0.0.0.0:" + strconv.Itoa(localPort)

	l, err = net.Listen("tcp", host)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer l.Close()

	log.Println("Server listening on " + host + " ...")
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
	targetAddrBytes := localConn.ReadWithLength()
	targetAddr := string(targetAddrBytes)
	//log.Println(targetAddr)

	c, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Fatal("Can't connect to", targetAddr)
	}

	remoteConn := network.Connection{Conn: c}

	go func() {
		for {
			// local -> server
			buf := localConn.ReadWithLength()
			if buf == nil {
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
			buf := remoteConn.Read(1024)
			if buf == nil {
				localConn.Close()
				break
			}
			buf = codec.Encrypt(buf)
			// server -> local
			localConn.WriteWithLength(buf)
		}
	}()
}
