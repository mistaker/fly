# fly
基于epoll 边缘模式实现的 golang reactor socket 框架

### http协议 
```go
package main

import (
	"fly/agreement/http"
	"fmt"
)

func main() {
	http.Run(":8000", func(req *http.Request, resp *http.Response) {
		fmt.Println(req.Url)
		resp.Write([]byte("hello world "))
	})
}

``` 

### 可以自己模仿 ``fly/agreement/http`` 很轻松的实现自己的应用层协议


参考自
https://github.com/Allenxuxu/gev
