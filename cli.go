package main

import (
	"fmt"
	"net"
)

func main() {

	conn, _ := net.Dial("tcp", "127.0.0.1:8000")

	fmt.Println("start write")

	n, err := conn.Write([]byte("hello fay"))
	fmt.Println(n, err)

	var data = make([]byte, 4)
	_, err = conn.Read(data)
	fmt.Println(err)
	fmt.Println(data)

}
