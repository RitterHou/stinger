package network

import (
	"encoding/binary"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
)

// 对连接进行包装
type Connection struct {
	conn   net.Conn
	closed bool
}

// 使用已有连接创建自定义连接
func New(conn net.Conn) Connection {
	logrus.Debugf("New connection %s -> %s", conn.LocalAddr(), conn.RemoteAddr())
	return Connection{conn: conn, closed: false}
}

// 创建自定义连接并根据address进行远程连接
func Connect(address string) (Connection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		logrus.Warn(err)
		return Connection{}, err
	}
	logrus.Debugf("New connection %s -> %s", conn.LocalAddr(), conn.RemoteAddr())
	return Connection{conn: conn, closed: false}, nil
}

func (c Connection) Read(length uint32) ([]byte, error) {
	if c.closed {
		return nil, errors.New("Connection has been closed\n")
	}

	var buf = make([]byte, length)
	var bufSize, err = c.conn.Read(buf)
	if err != nil {
		c.Close()
		return nil, err
	}
	return buf[:bufSize], nil
}

func (c Connection) Write(data []byte) error {
	if c.closed {
		return errors.New("Connection has been closed\n")
	}

	// 读操作默认不会超时
	_, err := c.conn.Write(data)
	if err != nil {
		logrus.Warn(err)
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
	logrus.Debugf("Close connection %s -> %s", c.LocalAddress(), c.RemoteAddress())
	c.closed = true
	c.conn.Close()
}

func (c Connection) RemoteAddress() string {
	return c.conn.RemoteAddr().String()
}

func (c Connection) LocalAddress() string {
	return c.conn.LocalAddr().String()
}
