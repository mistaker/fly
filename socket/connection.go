package socket

import (
	"bytes"
	"log"
	"sync"
	"sync/atomic"

	"fly/events"

	"golang.org/x/sys/unix"
)

type (
	Connection struct {
		isClose     int32 // 1是close
		readData    []byte
		readBuffer  bytes.Buffer
		writeBuffer bytes.Buffer
		rLock       sync.Mutex //todo 换成自旋锁
		wLock       sync.Mutex //todo 换成自旋锁
		Fd          int
		Agreement   interface{} //存储协议数据
	}
)

func NewConnection(fd int) error {
	conn := &Connection{
		Fd:       fd,
		readData: make([]byte, 500),
	}

	ConnectionsMap.Store(fd, conn)
	onConnectFunc(conn)

	if err := eventPoll.AddSocketEnableRead(fd); err != nil {
		return err
	}

	return nil
}

func (conn *Connection) SendData(data []byte) {
	conn.wLock.Lock()
	defer conn.wLock.Unlock()

	//如果 当前 write 缓冲区 还有 数据未发送,那直接 append 到缓冲区中
	if conn.writeBuffer.Len() != 0 {
		conn.writeBuffer.Write(data)
		return
	}

	//尝试对 socket 进行写数据
	nums, err := unix.Write(conn.Fd, data)
	if err != nil && err != unix.EAGAIN {
		conn.handleClose()
		return
	}

	if nums >= len(data) {
		return
	}

	conn.writeBuffer.Write(data[nums:])
	if conn.writeBuffer.Len() > 0 {
		if err := eventPoll.EnableReadWrite(conn.Fd); err != nil {
			log.Printf("enableReadWrite err %v", err)
		}
	}
}

func (conn *Connection) HandleEvent(fd int, event events.Event) {
	if event&events.EventErr != 0 {
		conn.handleClose()
		return
	}

	if event&events.EventWrite != 0 {
		conn.handleWrite()
	}

	if event&events.EventRead != 0 {
		conn.handleRead()
	}
}

func (conn *Connection) handleWrite() {
	conn.wLock.Lock()
	defer conn.wLock.Unlock()

	nums, err := unix.Write(conn.Fd, conn.writeBuffer.Bytes())
	if err != nil {
		if err == unix.EAGAIN {
			return
		}

		conn.handleClose()
		return
	}

	conn.writeBuffer.Next(nums)
	if conn.writeBuffer.Len() == 0 {
		if err = eventPoll.EnableRead(conn.Fd); err != nil {
			log.Printf("EnableRead err %v", err)
		}
	}
}

func (conn *Connection) handleRead() {
	conn.rLock.Lock()
	defer conn.rLock.Unlock()

	for {
		nums, err := unix.Read(conn.Fd, conn.readData)

		if err == unix.EAGAIN {
			return
		}

		if err != nil || nums == 0 {
			conn.handleClose()
			return
		}

		if nums == len(conn.readData) {
			conn.readBuffer.Write(conn.readData)
			continue
		}

		conn.readBuffer.Write(conn.readData[:nums])
		break
	}

	onMessageFunc(conn, conn.readBuffer.Bytes())
	conn.readBuffer.Reset()
}

func (conn *Connection) handleClose() {
	if atomic.CompareAndSwapInt32(&conn.isClose, 0, 1) == false {
		return
	}

	onCLose(conn)

	eventPoll.DeleteFdInLoop(conn.Fd)

	ConnectionsMap.Delete(conn.Fd)

	if err := unix.Close(conn.Fd); err != nil {
		log.Printf("unix close err %v", err)
	}
}
