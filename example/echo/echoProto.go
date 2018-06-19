package echo

import (
    "net"
    "fmt"
    "bufio"
    "github.com/easytcp"
)

type EchoPacket struct {
    content string
}

func (e *EchoPacket) Serialize () []byte {
    return []byte(e.content)
}

func (e *EchoPacket) GetContent() string {
    return e.content
}

type EchoProto struct {

}

func (e * EchoProto) ReadPacket (conn *net.TCPConn) (easytcp.Packet, error) {
    reader := bufio.NewReader(conn)

    message,err := reader.ReadString('\n')
    fmt.Printf("read a packet : %s", message)
    if err != nil {
            conn.Close()
            //return easytcp.Packet{}, err
        }
        return &EchoPacket {
            content: message,
        }, nil
    }
