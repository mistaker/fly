package server

import (
	"fmt"
	"log"
	"sync"

	"fly/common"
	"fly/events"
	"fly/socket"
)

type (
	Server struct {
		listen *socket.Listener
		ep     *events.EventLoop
		wg     sync.WaitGroup
		lock   sync.Mutex
	}
)

func NewServer(network, addr string, onConnect func(conn *socket.Connection),
	onMessage func(conn *socket.Connection, data []byte), onClose func(conn *socket.Connection)) (*Server, error) {

	ep, err := events.NewEventLoop()
	if err != nil {
		log.Printf("new epollevent err %v ", err)
		return nil, err
	}

	//注册逻辑
	socket.RegisterLogic(ep, onConnect, onMessage, onClose)

	listener, err := socket.NewListener(network, addr)
	if err != nil {
		log.Printf("new Listener err %v ", err)
		return nil, err
	}

	return &Server{
		ep:     ep,
		listen: listener,
	}, nil
}

func (s *Server) Worker() {
	for {
		select {
		case eventDataItem := <-events.EventChan:
			//tcp listener
			if eventDataItem.Fd == s.listen.Fd {
				_ = s.listen.HandleEvent(eventDataItem.Fd, eventDataItem.Events)
			}

			//connection write read
			if client, ok := socket.ConnectionsMap.Load(eventDataItem.Fd); ok {
				client.(*socket.Connection).HandleEvent(eventDataItem.Fd, eventDataItem.Events)
			}
		}
	}
}

func (s *Server) Run() {

	common.GoSafe(func() {
		s.ep.RunPoll()
	})

	common.GoSafe(func() {
		s.Worker()
	})

	common.GoSafe(func() {
		s.Worker()
	})

	common.GoSafe(func() {
		s.Worker()
	})

	fmt.Println("server start....")
	fmt.Println()

	select {}
}
