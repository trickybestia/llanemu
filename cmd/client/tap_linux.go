//go:build linux

package main

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

func createTAP() (ifce *water.Interface, err error) {
	config := water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: args.Name,
		},
	}

	tap, err := water.New(config)

	if err != nil {
		return nil, err
	}

	if err = configureTAP(); err != nil {
		tap.Close()

		return nil, err
	}

	return tap, nil
}

func configureTAP() error {
	mask_size, _ := network.Mask.Size()

	broadcastAddress := net.IP(make([]byte, len(address)))

	copy(broadcastAddress, address)

	for i, mask_octet := range network.Mask {
		broadcastAddress[i] |= ^mask_octet
	}

	if err := exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%d", address, mask_size), "brd", broadcastAddress.String(), "dev", args.Name).Run(); err != nil {
		return err
	}

	if err := exec.Command("ip", "link", "set", "up", "dev", args.Name).Run(); err != nil {
		return err
	}

	if err := exec.Command("ip", "route", "add", "255.255.255.255", "dev", args.Name).Run(); err != nil {
		return err
	}

	if err := exec.Command("ip", "route", "add", "224.0.0.0/24", "dev", args.Name).Run(); err != nil {
		return err
	}

	return nil
}
