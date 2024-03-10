package types

import (
	"bytes"
	"testing"
)

// go test -v ./types
func TestTypes(t *testing.T) {
	testBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x07, 0x98, 0x00}
	testString1 := "hello MQTT"
	testString2 := "i like you so"

	// 01 MQTT_BYTE from streams
	reader_1 := bytes.NewReader(testBytes)
	// expect 0x01 1 nil
	b_1_1, n, err := NewByte().FromStream(reader_1)
	if b_1_1.data != 0x01 || n != 1 || err != nil {
		t.Error()
	}
	// expect 0x02 1 nil
	b_1_2, n, err := NewByte().FromStream(reader_1)
	if b_1_2.data != 0x02 || n != 1 || err != nil {
		t.Error()
	}

	// 02 MQTT_BYTE to stream
	buffer_2 := bytes.NewBuffer(nil)
	// expect 1, nil, [0x01]
	n, err = b_1_1.ToStream(buffer_2)
	if n != 1 || err != nil || len(buffer_2.Bytes()) != 1 || buffer_2.Bytes()[0] != 0x01 {
		t.Error()
	}
	// expect 1, nil, [0x01, 0x02]
	n, err = b_1_2.ToStream(buffer_2)
	if n != 1 || err != nil || len(buffer_2.Bytes()) != 2 || buffer_2.Bytes()[0] != 0x01 || buffer_2.Bytes()[1] != 0x02 {
		t.Error()
	}

	// 03 MQTT_U16 from stream
	reader_3 := bytes.NewReader(testBytes)
	// expect 258 2 nil
	u16_3_1, n, err := NewU16().FromStream(reader_3)
	if u16_3_1.data != 258 || n != 2 || err != nil {
		t.Error()
	}
	// expect 772 2 nil
	u16_3_2, n, err := NewU16().FromStream(reader_3)
	if u16_3_2.data != 772 || n != 2 || err != nil {
		t.Error()
	}

	// 04 MQTT_U16 to stream
	buffer_4 := bytes.NewBuffer(nil)
	// expect 2, nil, [0x01, 0x02]
	n, err = u16_3_1.ToStream(buffer_4)
	if n != 2 || err != nil || len(buffer_4.Bytes()) != 2 ||
		buffer_4.Bytes()[0] != 0x01 || buffer_4.Bytes()[1] != 0x02 {
		t.Error()
	}
	// expect 2, nil, [0x01, 0x02, 0x03, 0x04]
	n, err = u16_3_2.ToStream(buffer_4)
	if n != 2 || err != nil || len(buffer_4.Bytes()) != 4 ||
		buffer_4.Bytes()[2] != 0x03 || buffer_4.Bytes()[3] != 0x04 {
		t.Error()
	}

	// 05 MQTT_U32 from stream
	reader_5 := bytes.NewReader(testBytes)
	// expect 16909060 4 nil
	u32_5_1, n, err := NewU32().FromStream(reader_5)
	if u32_5_1.data != 16909060 || n != 4 || err != nil {
		t.Error()
	}
	// expect 2442302356 4 nil
	u32_5_2, n, err := NewU32().FromStream(reader_5)
	if u32_5_2.data != 2442302356 || n != 4 || err != nil {
		t.Error()
	}

	// 06 MQTT_U32 to stream
	buffer_6 := bytes.NewBuffer(nil)
	// expect 4, nil, [0x01, 0x02, 0x03, 0x04]
	n, err = u32_5_1.ToStream(buffer_6)
	if n != 4 || err != nil || len(buffer_6.Bytes()) != 4 ||
		buffer_6.Bytes()[0] != 0x01 || buffer_6.Bytes()[1] != 0x02 ||
		buffer_6.Bytes()[2] != 0x03 || buffer_6.Bytes()[3] != 0x04 {
		t.Error()
	}
	// expect 4, nil, [0x01, 0x02, 0x03, 0x04, 0x91, 0x92, 0x93, 0x94]
	n, err = u32_5_2.ToStream(buffer_6)
	if n != 4 || err != nil || len(buffer_6.Bytes()) != 8 ||
		buffer_6.Bytes()[4] != 0x91 || buffer_6.Bytes()[5] != 0x92 ||
		buffer_6.Bytes()[6] != 0x93 || buffer_6.Bytes()[7] != 0x94 {
		t.Error()
	}

	// 07 MQTT_UTF8 from value
	// expect nil 10 [104 101 108 108 111 32 77 81 84 84]
	utf8_7, err := NewUtf8().FromValue(testString1)
	if err != nil || utf8_7.len.ToValue() != 10 || string(utf8_7.data) != testString1 {
		t.Error()
	}

	// 08 MQTT_UTF8 to stream
	buffer_8 := bytes.NewBuffer(nil)
	// expect nil 12 [0 10 104 101 108 108 111 32 77 81 84 84]
	n, err = utf8_7.ToStream(buffer_8)
	if err != nil || n != 12 ||
		buffer_8.String() != string([]byte{0, 10, 104, 101, 108, 108, 111, 32, 77, 81, 84, 84}) {
		t.Error()
	}

	// 09 MQTT_UTF8 from stream
	utf8_9, n, err := NewUtf8().FromStream(buffer_8)
	// expect nil 12 hello MQTT
	if err != nil || n != 12 || utf8_9.ToValue() != testString1 {
		t.Error()
	}

	// 10 MQTT_UTF8_PAIR from value
	up_10, err := NewUtf8Pair().FromValue(testString1, testString2)
	// expect nil, testString1, testString2
	if err != nil || up_10.key.ToValue() != testString1 || up_10.value.ToValue() != testString2 {
		t.Error()
	}

	// 11 MQTT_UTF8_PAIR to stream
	buffer_11 := bytes.NewBuffer(nil)
	expect_11 := []byte{0, 10, 104, 101,
		108, 108, 111, 32,
		77, 81, 84, 84,
		0, 13, 105, 32,
		108, 105, 107, 101,
		32, 121, 111, 117,
		32, 115, 111}
	n, err = up_10.ToStream(buffer_11)
	if err != nil || n != len(expect_11) || buffer_11.Bytes()[0] != expect_11[0] ||
		buffer_11.Bytes()[1] != expect_11[1] || buffer_11.Bytes()[3] != expect_11[3] ||
		buffer_11.Bytes()[5] != expect_11[5] {
		t.Error()
	}

	// 12 MQTT_UTF8_PARI from stream
	reader_12 := bytes.NewReader(expect_11)
	up_12, n, err := NewUtf8Pair().FromStream(reader_12)
	// expect nil, testString1, testString2
	if err != nil || n != len(expect_11) || up_12.key.ToValue() != testString1 || up_12.value.ToValue() != testString2 {
		t.Error()
	}
}

func TestVarInt(t *testing.T) {
	testBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x07, 0x98}
	testBytes_2 := []byte{0x80, 0x00}

	// varInt from stream
	// expect [0x01] 1 nil
	reader_1 := bytes.NewReader(testBytes)
	mvi_1, n, err := NewVarInt().FromStream(reader_1)
	if len(mvi_1.data) != 1 || mvi_1.data[0] != 0x01 || n != 1 || err != nil {
		t.Error()
	}

	// expect length error
	reader_2 := bytes.NewReader(testBytes[4:])
	_, _, err = NewVarInt().FromStream(reader_2)
	if err == nil {
		t.Error()
	}

	// expect [0x94, 0x95, 0x96, 0x07] 4 nil
	reader_3 := bytes.NewReader(testBytes[7:])
	mvi_3, n, err := NewVarInt().FromStream(reader_3)
	if len(mvi_3.data) != 4 || mvi_3.data[0] != 0x94 || mvi_3.data[3] != 0x07 || n != 4 || err != nil {
		t.Error()
	}

	// expect EOF error
	reader_4 := bytes.NewReader(testBytes[11:])
	_, _, err = NewVarInt().FromStream(reader_4)
	if err == nil {
		t.Error()
	}

	// expect height byte can't be 0 err
	reader_5 := bytes.NewReader(testBytes_2)
	_, _, err = NewVarInt().FromStream(reader_5)
	if err == nil {
		t.Error()
	}

	// var int from value
	v1 := 1
	v2 := 399
	v3 := -15
	v4 := MQTT_VAR_INT_MAX - 2
	v5 := MQTT_VAR_INT_MAX
	v6 := MQTT_VAR_INT_MAX + 1

	// expect [0x01] nil
	mvi_6, err := NewVarInt().FromValue(v1)
	if len(mvi_6.data) != 1 || mvi_6.data[0] != 0x01 || err != nil {
		t.Error()
	}

	// expect [0x8F 0x03] nil
	mvi_7, err := NewVarInt().FromValue(v2)
	if len(mvi_7.data) != 2 || mvi_7.data[0] != 0x8F || mvi_7.data[1] != 0x03 || err != nil {
		t.Error()
	}

	// expect value ranges err
	_, err = NewVarInt().FromValue(v3)
	if err == nil {
		t.Error()
	}

	// expect [0xFD, 0xFF, 0xFF, 0x7F] nil
	mvi_9, err := NewVarInt().FromValue(v4)
	if len(mvi_9.data) != 4 || mvi_9.data[0] != 0xFD ||
		mvi_9.data[1] != 0xFF || mvi_9.data[2] != 0xFF ||
		mvi_9.data[3] != 0x7F || err != nil {
		t.Error()
	}

	// expect [0xFF, 0xFF, 0xFF, 0x7F] nil
	mvi_10, err := NewVarInt().FromValue(v5)
	if len(mvi_10.data) != 4 || mvi_10.data[0] != 0xFF ||
		mvi_10.data[1] != 0xFF || mvi_10.data[2] != 0xFF ||
		mvi_10.data[3] != 0x7F || err != nil {
		t.Error()
	}

	// expect value ranges err
	_, err = NewVarInt().FromValue(v6)
	if err == nil {
		t.Error()
	}

	// var int to stream
	testToStream_1 := []byte{0x01}
	testToStream_2 := []byte{0x94, 0x95, 0x96, 0x07}

	// expect [0x01] nil
	reader_12 := bytes.NewReader(testToStream_1)
	mvi_12, _, _ := NewVarInt().FromStream(reader_12)
	buffer := bytes.NewBuffer(nil)
	n, err = mvi_12.ToStream(buffer)
	if len(buffer.Bytes()) != 1 || buffer.Bytes()[0] != 0x01 || n != 1 || err != nil {
		t.Error()
	}

	// expect [0x94, 0x95, 0x96, 0x07] nil
	reader_13 := bytes.NewReader(testToStream_2)
	mvi_13, _, _ := NewVarInt().FromStream(reader_13)
	buffer_2 := bytes.NewBuffer(nil)
	n, err = mvi_13.ToStream(buffer_2)
	if len(buffer_2.Bytes()) != 4 || buffer_2.Bytes()[0] != 0x94 ||
		buffer_2.Bytes()[1] != 0x95 || buffer_2.Bytes()[2] != 0x96 ||
		buffer_2.Bytes()[3] != 0x07 || n != 4 || err != nil {
		t.Error()
	}

	// var int to value
	testToValue_1 := []byte{0x01}
	testToValue_2 := []byte{0x94, 0x95, 0x96, 0x07}

	// expect 1
	reader_14 := bytes.NewReader(testToValue_1)
	mvi_14, _, _ := NewVarInt().FromStream(reader_14)
	v := mvi_14.ToValue()
	if v != 1 {
		t.Error()
	}

	// expect 15043220
	reader_15 := bytes.NewReader(testToValue_2)
	mvi_15, _, _ := NewVarInt().FromStream(reader_15)
	v = mvi_15.ToValue()
	// 0x14 * 128^0 + 0x15 * 128^1 + 0x16 * 128^2 + 0x07 * 128^3
	if v != 15043220 {
		t.Error()
	}
}

func TestBin(t *testing.T) {
	testBytes_1 := []byte{0x00, 0x03, 0x01, 0x02, 0x03, 0x04}
	testBytes_2 := []byte{0x00, 0x03, 0x01, 0x02}

	// 01 from stream
	// expect nil, 5, [0x01, 0x02, 0x03]
	reader_1 := bytes.NewReader(testBytes_1)
	bin_1, n, err := NewBin().FromStream(reader_1)
	if err != nil || n != 5 || bin_1.data[0] != 0x01 || bin_1.data[2] != 0x03 {
		t.Error()
	}

	// 02 err
	// expect err
	reader_2 := bytes.NewReader(testBytes_2)
	_, _, err = NewBin().FromStream(reader_2)
	if err == nil {
		t.Error()
	}

	// 03 to stream
	// expect nil, 5, [0x00, 0x03, 0x01, 0x02, 0x03]
	buffer_3 := bytes.NewBuffer(nil)
	n, err = bin_1.ToStream(buffer_3)
	if err != nil || n != 5 || buffer_3.Bytes()[0] != 0x00 || buffer_3.Bytes()[1] != 0x03 ||
		buffer_3.Bytes()[2] != 0x01 || buffer_3.Bytes()[4] != 0x03 {
		t.Error()
	}
}
