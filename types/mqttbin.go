package types

import (
	"encoding/binary"
	"errors"
	"io"
)

const MQTT_BIN_MAX = 65535 // 0xFFFF

func NewBin() *MQTT_BIN {
	return &MQTT_BIN{}
}

func (mbin *MQTT_BIN) FromStream(input io.Reader) (*MQTT_BIN, int, error) {
	err := binary.Read(input, binary.BigEndian, &mbin.len.data)
	if err != nil {
		return nil, 0, err
	}
	mbin.data = make([]byte, mbin.len.data)
	err = binary.Read(input, binary.BigEndian, mbin.data)
	if err != nil {
		return nil, 0, err
	}
	return mbin, 2 + int(mbin.len.data), nil
}

func (mbin *MQTT_BIN) ToStream(output io.Writer) (int, error) {
	n, err := mbin.len.ToStream(output)
	if err != nil {
		return n, err
	}
	err = binary.Write(output, binary.BigEndian, mbin.data)
	if err != nil {
		return n, err
	}
	return 2 + int(mbin.len.data), nil
}

func (mbin *MQTT_BIN) FromValue(bs []byte) (*MQTT_BIN, error) {
	if len(bs) > MQTT_BIN_MAX {
		return nil, errors.New("MQTT_BIN length error")
	}
	mbin.len.data = uint16(len(bs))
	mbin.data = bs
	return mbin, nil
}

func (mbin *MQTT_BIN) ToValue() []byte {
	return mbin.data
}
