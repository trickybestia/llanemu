package llanemu

import (
	"encoding/binary"
	"io"
)

func ReadPacket(reader io.Reader) ([]byte, error) {
	packetLenBuf := make([]byte, 2)

	_, err := io.ReadFull(reader, packetLenBuf)

	if err != nil {
		return nil, err
	}

	packetLen := binary.LittleEndian.Uint16(packetLenBuf)

	packetBuf := make([]byte, packetLen)

	_, err = io.ReadFull(reader, packetBuf)

	if err != nil {
		return nil, err
	}

	return packetBuf, nil
}

func WritePacket(writer io.Writer, data []byte) error {
	packetLenBuf := make([]byte, 2)

	binary.LittleEndian.PutUint16(packetLenBuf, uint16(len(data)))

	if _, err := writer.Write(packetLenBuf); err != nil {
		return err
	}

	if _, err := writer.Write(data); err != nil {
		return err
	}

	return nil
}
