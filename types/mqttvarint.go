package types

import (
	"errors"
	"io"
)

const (
	MQTT_VAR_INT_MIN = 0         // 0x00
	MQTT_VAR_INT_MAX = 268435455 // 0xFF 0xFF 0xFF 0x7F
)

func NewVarInt() *MQTT_VAR_INT {
	return &MQTT_VAR_INT{}
}

func (mvi *MQTT_VAR_INT) FromStream(input io.Reader) (*MQTT_VAR_INT, int, error) {
	mvi.data = make([]byte, 0)
	buf := make([]byte, 1)
	x := 0
	// maximum 4bytes
	for ; x < 4; x++ {
		_, err := input.Read(buf)
		if err != nil {
			return nil, x, err
		}
		b := buf[0]
		// height byte can't be 0
		if x > 0 && b == 0 {
			return nil, x, errors.New("MQTT_VAR_INT read error")
		}
		mvi.data = append(mvi.data, b)
		// is there after
		if b&0x80 > 0 {
			if x == 3 {
				return nil, x, errors.New("MQTT_VAR_INT read error, length > 4")
			}
			continue
		}
		// over
		break
	}
	return mvi, x + 1, nil
}

func (mvi *MQTT_VAR_INT) ToStream(output io.Writer) (int, error) {
	return output.Write(mvi.data)
}

func (mvi *MQTT_VAR_INT) FromValue(v int) (*MQTT_VAR_INT, error) {
	if v < MQTT_VAR_INT_MIN || v > MQTT_VAR_INT_MAX {
		return nil, errors.New("MQTT_VAR_INT parse error, value out of range")
	}
	mvi.data = make([]byte, 0)
	for {
		var b byte = (byte)(v % 0x80)
		v = v / 0x80
		// set top bit
		if v > 0 {
			b = b | 0x80
		}
		// put
		mvi.data = append(mvi.data, b)
		// over
		if v == 0 {
			break
		}
	}
	return mvi, nil
}

func (mvi *MQTT_VAR_INT) ToValue() int {
	multiplier := 1
	value := 0
	for x := 0; x < len(mvi.data); x++ {
		b := mvi.data[x]
		value += int(b&0x7F) * multiplier
		multiplier *= 0x80
	}
	return value
}
