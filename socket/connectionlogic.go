package socket

import (
	"fly/events"
	"sync"
)

var (
	ConnectionsMap sync.Map
	eventPoll      *events.EventLoop
	onConnectFunc  func(conn *Connection)
	onMessageFunc  func(conn *Connection, data []byte)
	onCLose        func(conn *Connection)
)

func RegisterLogic(poll *events.EventLoop, onConnectF func(conn *Connection), onMessageF func(conn *Connection, data []byte), onCloseF func(conn *Connection)) {
	onConnectFunc = onConnectF
	onMessageFunc = onMessageF
	onCLose = onCloseF
	eventPoll = poll
}
