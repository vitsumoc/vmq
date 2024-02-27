package vmq

import (
	"bytes"
	"testing"
)

func TestRemainingLengthFromBytes(t *testing.T) {
	testBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x07, 0x98}
	testBytes_2 := []byte{0x80, 0x00}

	// expect [0x01] 1 nil
	reader_1 := bytes.NewReader(testBytes)
	rl_1, n, err := remainingLengthFromBytes(reader_1)
	if len(rl_1.bytes) != 1 || rl_1.bytes[0] != 0x01 || n != 1 || err != nil {
		t.Error()
	}

	// expect length error
	reader_2 := bytes.NewReader(testBytes[4:])
	_, _, err = remainingLengthFromBytes(reader_2)
	if err == nil {
		t.Error()
	}

	// expect [0x94, 0x95, 0x96, 0x07] 4 nil
	reader_3 := bytes.NewReader(testBytes[7:])
	rl_3, n, err := remainingLengthFromBytes(reader_3)
	if len(rl_3.bytes) != 4 || rl_3.bytes[0] != 0x94 || rl_3.bytes[3] != 0x07 || n != 4 || err != nil {
		t.Error()
	}

	// expect EOF error
	reader_4 := bytes.NewReader(testBytes[11:])
	_, _, err = remainingLengthFromBytes(reader_4)
	if err == nil {
		t.Error()
	}

	// expect height byte can't be 0 err
	reader_5 := bytes.NewReader(testBytes_2)
	_, _, err = remainingLengthFromBytes(reader_5)
	if err == nil {
		t.Error()
	}
}

func TestRemainingLengthFromValue(t *testing.T) {
	v1 := 1
	v2 := 399
	v3 := -15
	v4 := REMAININGLENGTH_MAX - 2
	v5 := REMAININGLENGTH_MAX
	v6 := REMAININGLENGTH_MAX + 1

	// expect [0x01] nil
	rl_1, err := remainingLengthFromValue(v1)
	if len(rl_1.bytes) != 1 || rl_1.bytes[0] != 0x01 || err != nil {
		t.Error()
	}

	// expect [0x8F 0x03] nil
	rl_2, err := remainingLengthFromValue(v2)
	if len(rl_2.bytes) != 2 || rl_2.bytes[0] != 0x8F || rl_2.bytes[1] != 0x03 || err != nil {
		t.Error()
	}

	// expect value ranges err
	_, err = remainingLengthFromValue(v3)
	if err == nil {
		t.Error()
	}

	// expect [0xFD, 0xFF, 0xFF, 0x7F] nil
	rl_4, err := remainingLengthFromValue(v4)
	if len(rl_4.bytes) != 4 || rl_4.bytes[0] != 0xFD ||
		rl_4.bytes[1] != 0xFF || rl_4.bytes[2] != 0xFF ||
		rl_4.bytes[3] != 0x7F || err != nil {
		t.Error()
	}

	// expect [0xFF, 0xFF, 0xFF, 0x7F] nil
	rl_5, err := remainingLengthFromValue(v5)
	if len(rl_5.bytes) != 4 || rl_5.bytes[0] != 0xFF ||
		rl_5.bytes[1] != 0xFF || rl_5.bytes[2] != 0xFF ||
		rl_5.bytes[3] != 0x7F || err != nil {
		t.Error()
	}

	// expect value ranges err
	_, err = remainingLengthFromValue(v6)
	if err == nil {
		t.Error()
	}
}

func TestToBytes(t *testing.T) {
	testBytes_1 := []byte{0x01}
	testBytes_2 := []byte{0x94, 0x95, 0x96, 0x07}

	// expect [0x01] nil
	reader_1 := bytes.NewReader(testBytes_1)
	rl_1, _, _ := remainingLengthFromBytes(reader_1)
	buffer := bytes.NewBuffer(nil)
	n, err := rl_1.toBytes(buffer)
	if len(buffer.Bytes()) != 1 || buffer.Bytes()[0] != 0x01 || n != 1 || err != nil {
		t.Error()
	}

	// expect [0x94, 0x95, 0x96, 0x07] nil
	reader_2 := bytes.NewReader(testBytes_2)
	rl_2, _, _ := remainingLengthFromBytes(reader_2)
	buffer_2 := bytes.NewBuffer(nil)
	n, err = rl_2.toBytes(buffer_2)
	if len(buffer_2.Bytes()) != 4 || buffer_2.Bytes()[0] != 0x94 ||
		buffer_2.Bytes()[1] != 0x95 || buffer_2.Bytes()[2] != 0x96 ||
		buffer_2.Bytes()[3] != 0x07 || n != 4 || err != nil {
		t.Error()
	}
}

func TestToValue(t *testing.T) {
	testBytes_1 := []byte{0x01}
	testBytes_2 := []byte{0x94, 0x95, 0x96, 0x07}

	// expect 1
	reader_1 := bytes.NewReader(testBytes_1)
	rl_1, _, _ := remainingLengthFromBytes(reader_1)
	v := rl_1.toValue()
	if v != 1 {
		t.Error()
	}

	// expect 15043220
	reader_2 := bytes.NewReader(testBytes_2)
	rl_2, _, _ := remainingLengthFromBytes(reader_2)
	v = rl_2.toValue()
	// 0x14 * 128^0 + 0x15 * 128^1 + 0x16 * 128^2 + 0x07 * 128^3
	if v != 15043220 {
		t.Error()
	}
}
