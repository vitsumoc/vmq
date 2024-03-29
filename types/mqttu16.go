package types

import (
	"encoding/binary"
	"io"
)

func NewU16() *MQTT_U16 {
	return &MQTT_U16{}
}

func (mu16 *MQTT_U16) FromStream(input io.Reader) (int, error) {
	err := binary.Read(input, binary.BigEndian, &mu16.data)
	if err != nil {
		return 0, err
	}
	return mu16.Length(), nil
}

func (mu16 *MQTT_U16) ToStream(output io.Writer) (int, error) {
	err := binary.Write(output, binary.BigEndian, mu16.data)
	if err != nil {
		return 0, err
	}
	return mu16.Length(), nil
}

func (mu16 *MQTT_U16) FromValue(u16 uint16) *MQTT_U16 {
	mu16.data = u16
	return mu16
}

func (mu16 *MQTT_U16) ToValue() uint16 {
	return mu16.data
}

func (mu16 *MQTT_U16) Length() int {
	return 2
}
