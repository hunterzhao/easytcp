package easytcp

import (
    "fmt"
    "net"
)

type Conn struct {
    server      *Server
    clientConn  *net.TCPConn
    closeChan   chan struct{}
    sendChan    chan Packet
    receiveChan chan Packet
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
func (c *Conn) Run() {
    c.server.callback.OnConnect(c)

    go func() {
        for {
           select {
              case <-c.closeChan:
                  fmt.Printf("conn is been closed\n")
                  return
              default:
           }
           packet, _ := c.server.protocol.ReadPacket(c.clientConn)
           c.receiveChan <- packet
        }
    }()

    go func() {
        for {
            select {
            case <-c.closeChan:
                fmt.Printf("conn is been closed\n")
                return
            case p :=<- c.receiveChan:
                c.server.callback.OnMessage(c, p)

            }

        }
    }()

    go func() {
        for {
            select {
            case <-c.closeChan:
                fmt.Printf("conn is been closed\n")
                return
            case p :=<- c.sendChan:
                c.clientConn.Write(p.Serialize())
            }
        }
    }()
}


