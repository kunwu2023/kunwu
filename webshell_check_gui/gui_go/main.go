package main

import (
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"webshell_gui_go/lib/api"
	"webshell_gui_go/lib/config"
)

func main() {
	config.SystemType = runtime.GOOS
	config.Arch = runtime.GOARCH

	// 创建一个TCP服务器
	ln, err := net.Listen("tcp", ":0") // 使用":0"让系统自动选择一个可用端口
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	// 输出监听的端口号，Electron程序可以从这里获取端口号
	addr := ln.Addr().String()
	fmt.Println("Listening on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn) // 使用goroutine处理每个连接
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var request map[string]interface{}
	err := json.NewDecoder(conn).Decode(&request)
	if err != nil {
		panic(err)
	}

	response := api.MiddlewareProcessing(request)
	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	// 将数据分成多个部分并发送
	chunkSize := 4096
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		conn.Write(data[i:end])
	}
	conn.Write([]byte("\n")) // 添加换行符作为数据结束标志
}
