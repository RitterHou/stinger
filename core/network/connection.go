package network

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
)

type Connection struct {
	Conn net.Conn
}

func (c Connection) Read(length uint32) ([]byte, error) {
	conn := c.Conn

	var buf = make([]byte, length)
	var bufSize, err = conn.Read(buf)
	if err != nil {
		c.Close()
		return nil, err
	}
	return buf[:bufSize], nil
}

func (c Connection) Write(data []byte) error {
	// 读操作默认不会超时
	_, err := c.Conn.Write(data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c Connection) ReadWithLength() ([]byte, error) {
	lengthBuf, err := c.Read(4)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)
	buf, err := c.Read(length)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// 读取一个单独的字节
func (c Connection) ReadByte() (byte, error) {
	buf, err := c.Read(1)
	if err != nil {
		return 0, err
	}
	if len(buf) == 0 {
		return 0, errors.New("read buf size is 0")
	}
	return buf[0], nil
}

func (c Connection) WriteWithLength(source []byte) error {
	length := uint32(len(source))
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, length)
	err := c.Write(lengthBuf)
	if err != nil {
		return err
	}
	err = c.Write(source)
	if err != nil {
		return err
	}
	return nil
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
