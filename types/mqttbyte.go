package types

import (
	"encoding/binary"
	"io"
)

func NewByte() *MQTT_BYTE {
	return &MQTT_BYTE{}
}

func (mb *MQTT_BYTE) FromStream(input io.Reader) (*MQTT_BYTE, int, error) {
	err := binary.Read(input, binary.BigEndian, &mb.data)
	if err != nil {
		return nil, 0, err
	}
	return mb, 1, nil
}

func (mb *MQTT_BYTE) ToStream(output io.Writer) (int, error) {
	err := binary.Write(output, binary.BigEndian, mb.data)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (mb *MQTT_BYTE) FromValue(b byte) *MQTT_BYTE {
	mb.data = b
	return mb
}

func (mb *MQTT_BYTE) ToValue() byte {
	return mb.data
}
