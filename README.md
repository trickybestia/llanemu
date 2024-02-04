# llanemu

Linux LAN emulator.

## How does it work?

Client-side part of llanemu creates TAP device and sends every packet to server. On server packets are sent to every other connected client.

Client also sets up routing of all outgoing broadcast and multicast packets (255.255.255.255 and 224.0.0.0/24 subnetwork respectively) via newly created TAP device.

## Tested software

* Minecraft (finds games on "local" network)
* Civilization VI

## Installation

### Manually compile from source

1. [Install Go](https://go.dev/).
2. Run the following shell script:

```sh
git clone https://github.com/trickybestia/llanemu.git

cd llanemu

export CGO_ENABLED=0

go build -C cmd/client -o ../../build/llanemu -ldflags "-s -w"
go build -C cmd/server -o ../../build/llanemu-server -ldflags "-s -w"
```

Built binaries will be at `llanemu/build` directory.

## Usage

### On server

```sh
llanemu-server --endpoint 10.10.0.1:12345
```

### On client A

```sh
llanemu --remote-endpoint 10.10.0.1:12345 --name llanemu-tap -a 192.168.123.1/24
```

### On client B

```sh
llanemu --remote-endpoint 10.10.0.1:12345 --name llanemu-tap -a 192.168.123.2/24
```

## License

Licensed under [GNU GPLv3](COPYING) only.
