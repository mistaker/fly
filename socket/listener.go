package socket

import (
	"errors"
	"log"
	"net"
	"os"

	"fly/events"

	"golang.org/x/sys/unix"
)

type (
	Listener struct {
		file     *os.File
		listener net.Listener
		Fd       int
	}
)

func NewListener(network, addr string) (*Listener, error) {
	listener, err := net.Listen(network, addr)
	if err != nil {
		log.Printf("net.Listen err %v", err)
		return nil, err
	}

	tl, ok := listener.(*net.TCPListener)
	if !ok {
		return nil, errors.New("could not get file descriptor")
	}

	file, err := tl.File()
	if err != nil {
		return nil, err
	}

	fd := int(file.Fd())

	if err := poller.AddSocketEnableRead(fd); err != nil {
		log.Printf("listen poll filed err %v", err)
		return nil, err
	}

	return &Listener{
		file:     file,
		listener: listener,
		Fd:       fd,
	}, nil

}

func (l *Listener) HandleEvent(fd int, event events.Event) error {

	if event&events.EventRead == 1 {

		nfd, _, err := unix.Accept(l.Fd)

		if err != nil {
			if err != unix.EAGAIN {
				log.Printf("accept: %v", err)
			}
			return err
		}

		if err := unix.SetNonblock(nfd, true); err != nil {
			_ = unix.Close(nfd)
			log.Printf("set nonblock: %v", err)
			return err
		}

		if err := NewConnection(nfd); err != nil {
			log.Printf("new connection  err %v", err)
			return err
		}

		return nil
	}

	return nil
}
