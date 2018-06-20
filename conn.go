package easytcp

import (
    "fmt"
    "net"
    "sync"
    "sync/atomic"
)

type Conn struct {
    server      *Server
    clientConn  *net.TCPConn
    closeChan   chan struct{}
    sendChan    chan Packet
    receiveChan chan Packet
    closeOnce   sync.Once
    closeFlag   int32
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

func (c *Conn) IsClosed() bool {
    return atomic.LoadInt32(&c.closeFlag) == 1
}

func (c *Conn) Close() {
    c.closeOnce.Do(func(){
        atomic.StoreInt32(&c.closeFlag, 1)
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
                if c.IsClosed() != true {
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
                if c.IsClosed() != true{
                   c.clientConn.Write(p.Serialize())
                }
            }
        }
    }()
}


