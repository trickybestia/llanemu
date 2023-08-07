package main

import (
	"log"
	"net"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/trickybestia/llanemu/internal/llanemu"
)

type Args struct {
	Endpoint string `arg:"-e,--endpoint"`
}

var args Args
var conns = make(map[net.Conn]struct{})
var connsLock = sync.RWMutex{}

func relayFrame(frame []byte, sender net.Conn) {
	disconnectedConns := make([]net.Conn, 0)

	connsLock.RLock()

	for conn := range conns {
		if conn != sender {
			if err := llanemu.WritePacket(conn, frame); err != nil {
				disconnectedConns = append(disconnectedConns, conn)
			}
		}
	}

	connsLock.RUnlock()

	if len(disconnectedConns) != 0 {
		connsLock.Lock()

		for _, connection := range disconnectedConns {
			delete(conns, connection)

			connection.Close()
		}

		connsLock.Unlock()
	}
}

func handleConn(conn net.Conn) {
	connsLock.Lock()
	conns[conn] = struct{}{}
	connsLock.Unlock()

	for {
		packet, err := llanemu.ReadPacket(conn)

		if err != nil {
			break
		}

		relayFrame(packet, conn)
	}

	connsLock.Lock()
	delete(conns, conn)
	connsLock.Unlock()

	conn.Close()
}

func main() {
	arg.MustParse(&args)

	listener, err := net.Listen("tcp4", args.Endpoint)

	if err != nil {
		log.Fatalln(err)
	}

	totalConns := 0
	totalConnsLock := sync.Mutex{}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			totalConnsLock.Lock()
			totalConns += 1

			log.Printf("Client connected: %s. Total connections: %d.", conn.RemoteAddr().String(), totalConns)

			totalConnsLock.Unlock()

			handleConn(conn)

			totalConnsLock.Lock()
			totalConns -= 1

			log.Printf("Client disconnected: %s. Total connections: %d.", conn.RemoteAddr().String(), totalConns)

			totalConnsLock.Unlock()
		}()
	}
}
