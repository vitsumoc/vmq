package packets

import (
	"errors"
	"io"
)

const (
	REMAININGLENTH_MIN  = 0         // 0x00
	REMAININGLENGTH_MAX = 268435455 // 0xFF 0xFF 0xFF 0x7F
)

type remainingLength struct {
	bytes []byte
}

func remainingLengthFromBytes(input io.Reader) (rl *remainingLength, n int, err error) {
	rl = &remainingLength{}
	rl.bytes = make([]byte, 0)
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
			return nil, x, errors.New("remainingLength parse err, height byte can't be 0")
		}
		rl.bytes = append(rl.bytes, b)
		// is there after
		if b&0x80 > 0 {
			if x == 3 {
				return nil, x, errors.New("remainingLength parse err, length > 4")
			}
			continue
		}
		// over
		break
	}
	return rl, x + 1, nil
}

func remainingLengthFromValue(v int) (rl *remainingLength, err error) {
	if v < REMAININGLENTH_MIN || v > REMAININGLENGTH_MAX {
		return nil, errors.New("remainingLength make err, value ranges err")
	}
	rl = &remainingLength{}
	rl.bytes = make([]byte, 0)
	for {
		var b byte = (byte)(v % 0x80)
		v = v / 0x80
		// set top bit
		if v > 0 {
			b = b | 0x80
		}
		// put
		rl.bytes = append(rl.bytes, b)
		// over
		if v == 0 {
			break
		}
	}
	return rl, nil
}

func (rl remainingLength) toBytes(output io.Writer) (n int, err error) {
	return output.Write(rl.bytes)
}

func (rl remainingLength) toValue() int {
	multiplier := 1
	value := 0
	for x := 0; x < len(rl.bytes); x++ {
		b := rl.bytes[x]
		value += int(b&0x7F) * multiplier
		multiplier *= 0x80
	}
	return value
}
