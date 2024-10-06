package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message广播消息Channel的goroutine，一旦有消息就发送全部的在线User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		// 将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// BroadCast 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// 当前连接的业务
	//fmt.Println("链接建立成功")
	user := NewUser(conn, this)

	// 用户上线了，将用户加入OnlineMap中
	user.Online()

	//监听用户是否活跃的Channel
	isLive := make(chan bool)

	//广播当前用户上线信息
	//this.BroadCast(user, "已上线")

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			msg := string(buf[:n-1])

			//将得到的消息进行广播
			user.DoMessage(msg)

			//用户的任意消息，代表当前用户是一个活跃的
			isLive <- true
		}
	}()
	//当前handler阻塞
	for {
		select {
		case <-isLive:

		case <-time.After(time.Second * 300):
			//已经超时
			//将当前的User强制关闭
			user.SendMsg("你被踢了")

			//销毁用的资源
			close(user.C)

			//关闭连接
			conn.Close()

			//退出当前的handler
			return //runtime.Goexit()
		}
	}
}

// Start 启动服务器
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen.err:", err)
		return
	}
	// close socket
	defer listener.Close()

	// 启动监听Message的goroutine
	go this.ListenMessage()
	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
		}

		// do handler
		go this.Handler(conn)
	}
}
