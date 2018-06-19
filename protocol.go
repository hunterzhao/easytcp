package easytcp

import (
    "net"
)
//define package type and protocol
type Packet interface {
    Serialize() []byte
}

type Protocol interface {
    ReadPacket(conn *net.TCPConn) (Packet, error)
}
