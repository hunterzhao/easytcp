package easytcp

import (
    "fmt"
    "net"
    "sync"
)

type Conn struct {
    server      *Server
    clientConn  *net.TCPConn
    closeChan   chan struct{}
    sendChan    chan Packet
    receiveChan chan Packet
    closeOnce   sync.Once
}

func newConn (conn *net.TCPConn, s *Server) *Conn {
    return &Conn{
        server:       s,
        clientConn:   conn,
        closeChan:    make(chan struct{}),
        sendChan:     make(chan Packet, s.config.MaxLengthofPackage),
        receiveChan:  make(chan Packet, s.config.MaxLengthofPackage),
    }
}

func (c *Conn) Close() {
    c.closeOnce.Do(func(){
        close(c.closeChan)
        close(c.receiveChan)
        close(c.sendChan)
        c.clientConn.Close()
        c.server.callback.OnClose(c)
    })
}

func (c *Conn) Run() {
    c.server.callback.OnConnect(c)

    go func() {
        defer func() {
            c.Close()
        }()

        for {
           select {
              case <-c.closeChan:
                  fmt.Printf("conn is been closed\n")
                  return
              default:
           }
           if packet, err := c.server.protocol.ReadPacket(c.clientConn); err != nil {
               return
           } else {
               c.receiveChan <- packet
           }
        }
    }()

    go func() {
        defer func() {
            c.Close()
        }()
        for {
            select {
            case <-c.closeChan:
                fmt.Printf("conn is been closed\n")
                return
            case p :=<- c.receiveChan:
                if p != nil {
                   c.server.callback.OnMessage(c, p)
               }
            }
        }
    }()

    go func() {
        defer func() {
            c.Close()
        }()
        for {
            select {
            case <-c.closeChan:
                fmt.Printf("conn is been closed\n")
                return
            case p :=<- c.sendChan:
                if p != nil{
                   c.clientConn.Write(p.Serialize())
                }
            }
        }
    }()
}


