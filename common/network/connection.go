package network

import (
	"log"
	"net"
)

type Connection struct {
	Conn net.Conn
}

func (c Connection) Read(length uint32) []byte {
	conn := c.Conn

	var buf = make([]byte, length)
	var bufSize, err = conn.Read(buf)
	if bufSize == 0 {
		//log.Printf("Connection closed by client %s", conn.RemoteAddr())
		return nil
	}
	if err != nil {
		log.Println(err)
		c.Close()
		return nil
	}
	return buf[:bufSize]
}

func (c Connection) Write(data []byte) {
	_, err := c.Conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func (c Connection) Close() {
	c.Conn.Close()
}

func (c Connection) RemoteAddress() string {
	return c.Conn.RemoteAddr().String()
}

func (c Connection) LocalAddress() string {
	return c.Conn.LocalAddr().String()
}
