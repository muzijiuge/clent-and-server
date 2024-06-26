package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var clients = make(map[string]net.Conn)

func main() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("监听失败：", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("服务器正在运行并监听端口 :12345")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接失败：", err.Error())
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	username, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("读取用户名失败：", err.Error())
		return
	}
	username = strings.TrimSpace(username)
	fmt.Printf("用户 %s 已连接\n", username)
	clients[username] = conn

	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		strings.TrimSpace(msg)
		if msg == "exit" {
			break
		}

		var destUsername string
		var message string
		if idx := strings.Index(msg, ":"); idx != -1 {
			destUsername = msg[:idx]
			message = msg[idx+1:]
		} else {
			fmt.Fprintf(conn, "消息格式错误，请使用 '目标用户:消息内容'\n")
			continue
		}

		destConn, ok := clients[destUsername]
		if !ok {
			fmt.Fprintf(conn, "用户 '%s' 不存在或不在线\n", destUsername)
			continue
		}
		fmt.Fprintf(destConn, "[%s]: %s\n", username, message)
	}
	delete(clients, username)
	fmt.Printf("用户 '%s' 断开连接\n", username)
}
