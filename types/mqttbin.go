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

func (mbin *MQTT_BIN) FromStream(input io.Reader) (int, error) {
	err := binary.Read(input, binary.BigEndian, &mbin.len.data)
	if err != nil {
		return 0, err
	}
	mbin.data = make([]byte, mbin.len.data)
	err = binary.Read(input, binary.BigEndian, mbin.data)
	if err != nil {
		return 0, err
	}
	return mbin.Length(), nil
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
	return mbin.Length(), nil
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

func (mbin *MQTT_BIN) Length() int {
	return mbin.len.Length() + len(mbin.data)
}
