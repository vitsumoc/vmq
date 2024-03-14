package types

import (
	"encoding/binary"
	"errors"
	"io"
)

const MQTT_UTF_8_MAX = 65535 // 0xFFFF

func NewUtf8() *MQTT_UTF8 {
	return &MQTT_UTF8{}
}

func (mutf8 *MQTT_UTF8) FromStream(input io.Reader) (int, error) {
	err := binary.Read(input, binary.BigEndian, &mutf8.len.data)
	if err != nil {
		return 0, err
	}
	mutf8.data = make([]byte, mutf8.len.data)
	err = binary.Read(input, binary.BigEndian, mutf8.data)
	if err != nil {
		return 0, err
	}
	return mutf8.Length(), nil
}

func (mutf8 *MQTT_UTF8) ToStream(output io.Writer) (int, error) {
	n, err := mutf8.len.ToStream(output)
	if err != nil {
		return n, err
	}
	err = binary.Write(output, binary.BigEndian, mutf8.data)
	if err != nil {
		return n, err
	}
	return mutf8.Length(), nil
}

func (mutf8 *MQTT_UTF8) FromValue(s string) (*MQTT_UTF8, error) {
	if len(s) > MQTT_UTF_8_MAX {
		return nil, errors.New("MQTT_UTF8 length error")
	}
	mutf8.len.data = uint16(len(s))
	mutf8.data = []byte(s)
	return mutf8, nil
}

func (mutf8 *MQTT_UTF8) ToValue() string {
	return string(mutf8.data)
}

func (mutf8 *MQTT_UTF8) Length() int {
	return mutf8.len.Length() + len(mutf8.data)
}
