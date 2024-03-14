package packets

import (
	"bytes"
	"testing"
)

// go test -v ./packets
func TestProperties(t *testing.T) {
	p := NewProperties()

	// 1 test empty properties
	buffer := bytes.NewBuffer(nil)
	n, err := p.ToStream(buffer)
	// expect nil, 1, 1, [0x00]
	if err != nil || p.Length() != 1 || n != 1 || len(buffer.Bytes()) != 1 || buffer.Bytes()[0] != 0x00 {
		t.Error()
	}

	// 2 test set property
	p.SetProperty(PROPERTY_CONTENT_TYPE, "MQTT", nil)
	buffer.Reset()
	n, err = p.ToStream(buffer)
	// expect nil, 8, 8, [7 3 0 4 M Q T T]
	if err != nil || p.Length() != 8 || n != 8 || len(buffer.Bytes()) != 8 ||
		buffer.Bytes()[4] != 'M' || buffer.Bytes()[5] != 'Q' {
		t.Error()
	}

	// 3 test from stream
	p = NewProperties()
	n, err = p.FromStream(buffer)
	// expect nil, 8
	if err != nil || n != 8 || p.Length() != 8 {
		t.Error()
	}
	buffer.Reset()
	n, err = p.ToStream(buffer)
	// expect nil, 8, 8, [7 3 0 4 M Q T T]
	if err != nil || p.Length() != 8 || n != 8 || len(buffer.Bytes()) != 8 ||
		buffer.Bytes()[4] != 'M' || buffer.Bytes()[5] != 'Q' {
		t.Error()
	}

	// 4 test set user property
	err = p.SetProperty(PROPERTY_USER_PROPERTY, "M", "Q")
	if err != nil {
		t.Error()
	}
	// expect nil, 15, 15, [14 3 0 4 M Q T T 38 0 1 M 0 1 Q]
	buffer.Reset()
	n, err = p.ToStream(buffer)
	if err != nil || n != 15 || p.Length() != 15 ||
		buffer.Bytes()[0] != 14 || buffer.Bytes()[14] != 'Q' {
		t.Error()
	}
}
