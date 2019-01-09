package network

import "log"

// 当前所有的连接
var connections = make([]Connection, 50)

func remove(conn Connection) {
	for i, c := range connections {
		if c == conn {
			connections = append(connections[:i], connections[i+1:]...)
			return
		}
	}
	log.Printf("%s not in connections\n", conn)
}

func add(conn Connection) {
	connections = append(connections, conn)
}

func GetConnections() []Connection {
	return connections
}
