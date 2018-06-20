package main

import (
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
    "github.com/easytcp"
    "github.com/easytcp/example/echo"
)

type EchoCallback struct {

}

func (c *EchoCallback) OnConnect(conn *easytcp.Conn) {
    fmt.Printf("new connection coming\n")
}

func (c *EchoCallback) OnClose(conn *easytcp.Conn) {
    fmt.Printf("connecion closed\n")
}

func (c *EchoCallback) OnMessage(conn *easytcp.Conn, p easytcp.Packet) {
    ep := p.(*echo.EchoPacket)

    fmt.Printf("got a new message : %s ", ep.GetContent())
}

func main() {
    config := &easytcp.Config {
        MaxLengthofPackage: 20,
    }

    tcpAddr, _ := net.ResolveTCPAddr("tcp4", ":8989")
    listener, _ := net.ListenTCP("tcp", tcpAddr)

    svr := easytcp.CreateServer(&EchoCallback{}, &echo.EchoProto{}, config)

    go svr.Listen(listener)
    fmt.Printf("listen in 8989\n")

    chSig := make(chan os.Signal)
    signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
    fmt.Printf("Signal: ", <-chSig)
}




