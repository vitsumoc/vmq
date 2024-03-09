package types

import (
	"encoding/binary"
	"io"
)

func (mu32 *MQTT_U32) FromStream(input io.Reader) (*MQTT_U32, int, error) {
	err := binary.Read(input, binary.BigEndian, &mu32.data)
	if err != nil {
		return nil, 0, err
	}
	return mu32, 4, nil
}

func (mu32 *MQTT_U32) ToStream(output io.Writer) (int, error) {
	err := binary.Write(output, binary.BigEndian, mu32.data)
	if err != nil {
		return 0, err
	}
	return 4, nil
}

func (mu32 *MQTT_U32) FromValue(u32 uint32) *MQTT_U32 {
	mu32.data = u32
	return mu32
}

func (mu32 *MQTT_U32) ToValue() uint32 {
	return mu32.data
}
