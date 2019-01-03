package network

import (
	"encoding/binary"
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

func (c Connection) ReadWithLength() []byte {
	lengthBuf := c.Read(4)
	if len(lengthBuf) == 0 {
		// 链接断开
		return nil
	}
	//log.Println(lengthBuf)
	length := binary.BigEndian.Uint32(lengthBuf)
	return c.Read(length)
}

func (c Connection) WriteWithLength(source []byte) {
	length := uint32(len(source))
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, length)
	c.Write(lengthBuf)
	c.Write(source)
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
