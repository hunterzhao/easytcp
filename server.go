package easytcp

import (
    "net"
    "fmt"
)

type Config struct {
    MaxLengthofPackage uint32

}

type Server struct {
    config   *Config
    callback Callback
    protocol Protocol
}

type Callback interface {
    OnConnect (*Conn)
    OnMessage (*Conn, Packet)
    OnClose   (*Conn)
}

func CreateServer (call Callback, proto Protocol, conf *Config) *Server {
    return &Server {
        callback: call,
        protocol: proto,
        config: conf,
    }
}

func (s *Server) Listen(listener *net.TCPListener) {
    defer func() {
        listener.Close()
    }()

    for {
        fmt.Printf("wait for another conn\n")
        conn, err := listener.AcceptTCP()
        if err != nil {
           return
       }
        newConn(conn, s).Run()
    }
}

func (s *Server) Stop() {

}
