package types

import (
	"encoding/binary"
	"io"
)

func NewU32() *MQTT_U32 {
	return &MQTT_U32{}
}

func (mu32 *MQTT_U32) FromStream(input io.Reader) (int, error) {
	err := binary.Read(input, binary.BigEndian, &mu32.data)
	if err != nil {
		return 0, err
	}
	return mu32.Length(), nil
}

func (mu32 *MQTT_U32) ToStream(output io.Writer) (int, error) {
	err := binary.Write(output, binary.BigEndian, mu32.data)
	if err != nil {
		return 0, err
	}
	return mu32.Length(), nil
}

func (mu32 *MQTT_U32) FromValue(u32 uint32) *MQTT_U32 {
	mu32.data = u32
	return mu32
}

func (mu32 *MQTT_U32) Length() int {
	return 4
}
