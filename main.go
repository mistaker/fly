package main

import (
	"fmt"

	"fly/server"
	"fly/socket"
)

func main() {
	ser, _ := server.NewServer("tcp", ":8000", func(conn *socket.Connection) {
		fmt.Println("socket open ", conn.Fd)
	}, func(conn *socket.Connection, data []byte) {

		fmt.Println("socket message ", string(data))
		conn.SendData(data)
	}, func(conn *socket.Connection) {
		fmt.Println("socket close ", conn.Fd)
	})

	ser.Run()
}
