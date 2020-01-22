package events

import (
	"log"
)

type (
	EventData struct {
		Fd     int
		Events Event
	}

	EventLoop struct {
		pool *pool
	}
)

var (
	EventChan = make(chan EventData, 10)
)

func NewEventLoop() (*EventLoop, error) {
	p, err := NewPoll()
	if err != nil {
		return nil, err
	}

	return &EventLoop{
		pool: p,
	}, nil
}

//删除相关socket
func (el *EventLoop) DeleteFdInLoop(fd int) {
	if err := el.pool.Del(fd); err != nil {
		log.Printf("el.pool.Del err %v", err)
	}
}

//注册相关 socket 并且监听可读事件
func (el *EventLoop) AddSocketEnableRead(fd int) error {
	err := el.pool.AddRead(fd)
	if err != nil {
		return err
	}

	return nil
}

//监听相关 socket 可读可写事件
func (el *EventLoop) EnableReadWrite(fd int) error {
	return el.pool.EnableReadWrite(fd)
}

//只监听相关 socket 可读事件
func (el *EventLoop) EnableRead(fd int) error {
	return el.pool.EnableRead(fd)
}

//注册给pool 用, 把事件发送到 worker 协程进行操作
func (el *EventLoop) handleEvent(fd int, events Event) {
	EventChan <- EventData{
		Fd:     fd,
		Events: events,
	}
}

func (el *EventLoop) RunPoll() {
	el.pool.Run(el.handleEvent)
}
