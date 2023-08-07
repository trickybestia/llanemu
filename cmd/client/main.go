package main

import (
	"net"

	"github.com/alexflint/go-arg"
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

	address = address.To4()

	if err != nil {
		p.Fail(err.Error())
	}

	network = *address_network
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

	go func() {
		for {
			packet, err := llanemu.ReadPacket(conn)

			if err != nil {
				log.Fatal(err)
			}

			if _, err = tap.Write(packet); err != nil {
				log.Fatal(err, len(packet), packet)
			}
		}
	}()

	buf := make([]byte, 1600)

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
