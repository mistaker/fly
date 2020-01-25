package main

import "fly/agreement/http"

func main() {
	http.Run(":8000", func(req *http.Request, resp *http.Response) {
		resp.Write([]byte("hello world "))
	})
}
