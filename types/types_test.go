package types

import (
	"bytes"
	"testing"
)

// go test -v ./types
func TestTypes(t *testing.T) {
	testBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x07, 0x98, 0x00}
	testString1 := "hello MQTT"
	// testString2 := "i like you so"

	// 01 MQTT_BYTE from streams
	reader_1 := bytes.NewReader(testBytes)
	// expect 0x01 1 nil
	b_1_1, n, err := (&MQTT_BYTE{}).FromStream(reader_1)
	if b_1_1.data != 0x01 || n != 1 || err != nil {
		t.Error()
	}
	// expect 0x02 1 nil
	b_1_2, n, err := (&MQTT_BYTE{}).FromStream(reader_1)
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
	u16_3_1, n, err := (&MQTT_U16{}).FromStream(reader_3)
	if u16_3_1.data != 258 || n != 2 || err != nil {
		t.Error()
	}
	// expect 772 2 nil
	u16_3_2, n, err := (&MQTT_U16{}).FromStream(reader_3)
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
	u32_5_1, n, err := (&MQTT_U32{}).FromStream(reader_5)
	if u32_5_1.data != 16909060 || n != 4 || err != nil {
		t.Error()
	}
	// expect 2442302356 4 nil
	u32_5_2, n, err := (&MQTT_U32{}).FromStream(reader_5)
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
	utf8_7, err := (&MQTT_UTF8{}).FromValue(testString1)
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
	utf8_9, n, err := (&MQTT_UTF8{}).FromStream(buffer_8)
	// expect nil 12 hello MQTT
	if err != nil || n != 12 || utf8_9.ToValue() != testString1 {
		t.Error()
	}
}
