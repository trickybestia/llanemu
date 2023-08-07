package main

import (
	"net"

	"github.com/alexflint/go-arg"
	"github.com/songgao/water"
	"github.com/trickybestia/llanemu/internal/llanemu"

	"log"
)

type Args struct {
	RemoteEndpoint string `arg:"required,-r,--remote-endpoint"`
	Name           string `arg:"required,-n,--name"`
	Address        string `arg:"required,-a,--address"`
}

var args Args
var address net.IP
var network net.IPNet

func parseArgs() {
	p := arg.MustParse(&args)

	var err error
	var address_network *net.IPNet

	address, address_network, err = net.ParseCIDR(args.Address)

	if err != nil {
		p.Fail(err.Error())
	}

	address = address.To4()

	if address == nil {
		p.Fail("Valid IPv4 address expected")
	}

	network = *address_network
}

func pipeFromTAPToConn(tap *water.Interface, conn net.Conn) {
	buf := make([]byte, 2048)

	for {
		read, err := tap.Read(buf)

		if err != nil {
			log.Fatal(err)
		}

		if err = llanemu.WritePacket(conn, buf[:read]); err != nil {
			log.Fatal(err)
		}
	}
}

func pipeFromConnToTAP(tap *water.Interface, conn net.Conn) {
	for {
		packet, err := llanemu.ReadPacket(conn)

		if err != nil {
			log.Fatal(err)
		}

		if _, err = tap.Write(packet); err != nil {
			log.Fatal(err, len(packet), packet)
		}
	}
}

func main() {
	parseArgs()

	conn, err := net.Dial("tcp4", args.RemoteEndpoint)

	if err != nil {
		log.Fatal(err)
	}

	tap, err := createTAP()

	if err != nil {
		log.Fatal(err)
	}

	go pipeFromConnToTAP(tap, conn)
	pipeFromTAPToConn(tap, conn)
}
